package main

import (
	"sync"

	"github.com/francesconi/roborock-box/drv8825"
	"github.com/stianeikeland/go-rpio/v4"
)

type BoxConfig struct {
	Stepper   drv8825.Config
	DoorSteps int
}

type Box struct {
	mu        sync.Mutex
	stepper   *drv8825.Driver
	doorSteps int
	doorOpen  bool
}

func NewBox(cfg BoxConfig) (*Box, error) {
	if err := rpio.Open(); err != nil {
		return nil, err
	}

	stepper, err := drv8825.New(cfg.Stepper)
	if err != nil {
		rpio.Close()
		return nil, err
	}

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
	return rpio.Close()
}
