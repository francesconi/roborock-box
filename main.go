package main

import (
	"log"

	"github.com/kardianos/service"
)

func main() {
	prg, err := NewProgram()
	if err != nil {
		log.Fatal(err)
	}
	cfg := &service.Config{
		Name:        "roborock-garage",
		DisplayName: "Roborock Garage",
	}
	s, err := service.New(prg, cfg)
	if err != nil {
		log.Fatal(err)
	}

	// status, err := s.Status()
	// if errors.Is(err, service.ErrNotInstalled) || status == service.StatusUnknown {
	// 	s.Install()
	// 	if err = s.Start(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// } else {
	// 	if err = s.Run(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
