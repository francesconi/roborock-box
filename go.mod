module github.com/francesconi/roborock-garage

go 1.18

require (
	github.com/kardianos/service v1.2.2
	github.com/stianeikeland/go-rpio/v4 v4.6.0
	github.com/vkorn/go-miio v0.0.0-20180929223642-adf1adb6425f
)

replace github.com/vkorn/go-miio => github.com/francesconi/go-miio v0.0.0-20230314075917-2917d4107018

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40 // indirect
	github.com/nickw444/miio-go v0.0.0-20190825225226-379bc4c72748 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
)
