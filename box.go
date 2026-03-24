package main

import (
	"sync"

	"github.com/francesconi/roborock-box/drv8825"
)

type BoxConfig struct {
	PinEnable    uint8
	PinStep      uint8
	PinDirection uint8
	StepMode     drv8825.StepMode
	RPM          uint
	DoorSteps    int
}

type Box struct {
	mu        sync.Mutex
	stepper   *drv8825.Driver
	doorSteps int
	doorOpen  bool
}

func NewBox(cfg BoxConfig) (*Box, error) {
	stepper, err := drv8825.New(drv8825.Config{
		PinEnable:    cfg.PinEnable,
		PinStep:      cfg.PinStep,
		PinDirection: cfg.PinDirection,
		StepMode:     cfg.StepMode,
	})
	if err != nil {
		return nil, err
	}
	stepper.SetSpeed(cfg.RPM)
	return &Box{stepper: stepper, doorSteps: cfg.DoorSteps}, nil
}

func (b *Box) OpenDoor() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.doorOpen {
		b.stepper.Enable()
		b.stepper.Move(-b.doorSteps)
		b.stepper.Disable()
		b.doorOpen = true
	}
}

func (b *Box) CloseDoor() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.doorOpen {
		b.stepper.Enable()
		b.stepper.Move(b.doorSteps)
		b.stepper.Disable()
		b.doorOpen = false
	}
}

func (b *Box) Cleanup() error {
	return b.stepper.Close()
}
