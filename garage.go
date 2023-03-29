package main

import (
	"sync"

	"github.com/francesconi/roborock-box/drv8825"
)

type Garage struct {
	mu      sync.Mutex
	stepper *drv8825.Driver

	DoorOpen bool
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

	stepper.SetSpeed(600)

	return &Garage{
		mu:      sync.Mutex{},
		stepper: stepper,
	}, nil
}

func (g *Garage) OpenDoor() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.DoorOpen {
		g.stepper.Enable()
		g.stepper.Move(-4000)
		g.stepper.Disable()
		g.DoorOpen = true
	}
}

func (g *Garage) CloseDoor() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.DoorOpen {
		g.stepper.Enable()
		g.stepper.Move(4000)
		g.stepper.Disable()
		g.DoorOpen = false
	}
}

func (g *Garage) Cleanup() error {
	return g.stepper.Close()
}
