package main

import (
	"errors"
	"log"

	"github.com/kardianos/service"
	"github.com/vkorn/go-miio"
)

func main() {
	p, err := NewProgram()
	if err != nil {
		log.Fatal(err)
	}

	svcConfig := &service.Config{
		Name:        "roborock-garage",
		DisplayName: "Roborock Garage",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target",
		},
	}
	svc, err := service.New(p, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	status, err := svc.Status()
	if errors.Is(err, service.ErrNotInstalled) || status == service.StatusUnknown {
		svc.Install()
	}

	miio.LOGGER.SetLevel(miio.LogLevelInfo)

	if err = svc.Run(); err != nil {
		log.Fatal(err)
	}
}
