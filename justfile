module := "roborock-garage"
ip := "192.168.178.25"

build:
    GOOS=linux GOARCH=arm GOARM=6 go build -o {{module}}

deploy: build
    #!/usr/bin/env sh
    ssh pi@{{ip}} << EOF
        if service --status-all | grep -q "{{module}}"; then
            sudo service {{module}} stop
        fi
    EOF
    scp {{module}} pi@{{ip}}:/home/pi/
    ssh pi@{{ip}} << EOF
        if ! service --status-all | grep -q "{{module}}"; then
            sudo service {{module}} install
        fi
        sudo service {{module}} start
    EOF
