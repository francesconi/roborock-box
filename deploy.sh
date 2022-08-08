GOOS=linux GOARCH=arm GOARM=6 go build -o roborock-garage
scp roborock-garage pi@192.168.x.x:/home/pi/
ssh -t pi@192.168.x.x "sudo ./roborock-garage"