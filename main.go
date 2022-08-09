package main

import (
	"log"
	"roborock-garage/drv8825"

	"github.com/vkorn/go-miio"
)

func main() {
	stepper, err := drv8825.New(drv8825.Config{
		EN:   24,
		STEP: 23,
		DIR:  22,
		MS1:  21,
		MS2:  21,
		MS3:  21,
		Mode: drv8825.StepModeFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer stepper.Stop()

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
				openGarage(stepper)
			case miio.VacStateCharging:
				closeGarage(stepper)
			}
		}
	}()

	vacuum.UpdateStatus()
}

func openGarage(stepper *drv8825.Driver) {
	stepper.Enable()
	stepper.Forward(200)
	stepper.Disable()
}

func closeGarage(stepper *drv8825.Driver) {
	stepper.Enable()
	stepper.Backward(200)
	stepper.Disable()
}
