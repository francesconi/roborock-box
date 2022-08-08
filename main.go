package main

import (
	"log"
	"roborock-garage/drv8825"
	"time"
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

	stepper.Enable()
	stepper.Forward(200)
	time.Sleep(1 * time.Second)
	stepper.Backward(200)
	stepper.Disable()
}
