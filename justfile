module := "roborock-garage"
ip := "192.168.x.x"

build:
    GOOS=linux GOARCH=arm GOARM=6 go build -o {{module}}

deploy: build
    scp {{module}} pi@{{ip}}:/home/pi/

deploy-and-run: deploy
    ssh -t pi@{{ip}} "sudo ./{{module}}"