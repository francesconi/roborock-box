package main

import (
	"log"
	"time"

	"github.com/kardianos/service"
	"github.com/vkorn/go-miio"
)

type program struct {
	garage *Garage
	vacuum *miio.Vacuum
	exit   chan struct{}
}

func NewProgram() (*program, error) {
	garage, err := NewGarage()
	if err != nil {
		return nil, err
	}

	// vacuum, err := miio.NewVacuum("<ip>", "<token>")
	// if err != nil {
	// 	defer garage.Cleanup()
	// 	return nil, err
	// }

	return &program{
		exit:   make(chan struct{}),
		garage: garage,
		vacuum: nil,
	}, nil
}

func (p program) Start(s service.Service) error {
	log.Print("Starting service...")
	go p.run()
	return nil
}

func (p program) run() error {
	log.Printf("Service running on %v.", service.Platform())

	// go func() {
	// 	for msg := range p.vacuum.UpdateChan {
	// 		log.Printf("Got message %+v", msg.State)
	// 		state := msg.State.(*miio.VacuumState)
	// 		switch state.State {
	// 		case miio.VacStateCleaning:
	// 			p.garage.OpenDoor()
	// 		case miio.VacStateCharging:
	// 			p.garage.CloseDoor()
	// 		}
	// 	}
	// }()

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Print("Requesting vacuum status update...")
			// p.vacuum.UpdateStatus()
		case <-p.exit:
			ticker.Stop()
			p.garage.Cleanup()
			return nil
		}
	}
}

func (p program) Stop(s service.Service) error {
	log.Print("Stopping service...")
	close(p.exit)
	return nil
}

func (p program) Shutdown(s service.Service) error {
	log.Print("Shutting down service...")
	close(p.exit)
	return nil
}
