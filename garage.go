package main

import "github.com/francesconi/roborock-garage/drv8825"

type Garage struct {
	stepper *drv8825.Driver
}

func NewGarage() (*Garage, error) {
	stepper, err := drv8825.New(drv8825.Config{
		PinEnable:    24,
		PinStep:      23,
		PinDirection: 22,
		StepMode:     drv8825.StepModeFull,
	})
	if err != nil {
		return nil, err
	}

	stepper.SetSpeed(60)

	return &Garage{stepper}, nil
}

func (g Garage) OpenDoor() {
	g.stepper.Enable()
	g.stepper.Move(200)
	g.stepper.Disable()
}

func (g Garage) CloseDoor() {
	g.stepper.Enable()
	g.stepper.Move(-200)
	g.stepper.Disable()
}

func (g Garage) Cleanup() error {
	return g.stepper.Close()
}
