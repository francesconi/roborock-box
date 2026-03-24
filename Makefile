module  = roborock-box
host   ?=
user   ?=
remote  = $(user)@$(host)
bindir  = /usr/local/bin
envfile = /etc/default/$(module)
sshctl  = /tmp/.ssh-$(module)
ssh     = ssh -o ControlPath=$(sshctl)
scp     = scp -o ControlPath=$(sshctl)

ifndef user
$(error user is required: make deploy user=<username> host=<ip>)
endif

ifndef host
$(error host is required: make deploy user=<username> host=<ip>)
endif

.PHONY: build deploy uninstall

build:
	GOOS=linux GOARCH=arm GOARM=6 go build -o $(module)

deploy: build
	# Open a single SSH connection — password entered once
	ssh -MNf -o ControlPath=$(sshctl) -o ControlPersist=yes $(remote)
	# Stop and remove any existing installation
	-$(ssh) $(remote) "sudo service $(module) stop 2>/dev/null; sudo $(bindir)/$(module) uninstall 2>/dev/null"
	# Copy binary
	$(scp) $(module) $(remote):/tmp/$(module)
	$(ssh) $(remote) "sudo mv /tmp/$(module) $(bindir)/$(module) && sudo chmod +x $(bindir)/$(module)"
	# Create env file with placeholders if it doesn't exist yet
	$(ssh) $(remote) "[ -f $(envfile) ] || printf 'IP=\nTOKEN=\n' | sudo tee $(envfile) > /dev/null"
	# Install and start
	$(ssh) $(remote) "sudo $(bindir)/$(module) install && sudo service $(module) start"
	# Close the shared connection
	$(ssh) -O exit $(remote) 2>/dev/null; true
	@echo ""
	@echo "Deployed. Set IP and TOKEN in $(envfile) on the Pi if not done yet:"
	@echo "  ssh $(remote) sudo nano $(envfile)"
	@echo "Then restart:"
	@echo "  ssh $(remote) sudo service $(module) restart"

uninstall:
	ssh $(remote) "sudo service $(module) stop; sudo $(bindir)/$(module) uninstall; sudo rm -f $(bindir)/$(module) $(envfile)"
