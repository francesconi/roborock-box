package main

import (
	"log"

	"github.com/vkorn/go-miio"
)

func main() {
	garage, err := NewGarage()
	if err != nil {
		log.Fatal(err)
	}
	defer garage.Cleanup()

	vacuum, err := miio.NewVacuum("<ip>", "<token>")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for msg := range vacuum.UpdateChan {
			log.Printf("Got message %+v", msg.State)
			state := msg.State.(*miio.VacuumState)
			switch state.State {
			case miio.VacStateCleaning:
				garage.OpenDoor()
			case miio.VacStateCharging:
				garage.CloseDoor()
			}
		}
	}()

	vacuum.UpdateStatus()
}
