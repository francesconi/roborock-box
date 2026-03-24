package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// LightSensorWatcher returns a WatchFunc that monitors the light sensor on
// the given GPIO pin, closing the door when the beam is interrupted (robot
// present) and opening it when the beam is clear (robot absent).
func LightSensorWatcher(pin uint8) WatchFunc {
	return func(exit <-chan struct{}, onOpen, onClose func()) error {
		return watchLightSensor(pin, exit, onOpen, onClose)
	}
}

func watchLightSensor(pin uint8, exit <-chan struct{}, onOpen, onClose func()) error {
	if err := rpio.Open(); err != nil {
		return fmt.Errorf("initialize light sensor: %w", err)
	}
	defer rpio.Close()

	p := rpio.Pin(pin)
	p.Input()
	p.PullUp()
	p.Detect(rpio.AnyEdge)
	defer p.Detect(rpio.NoEdge)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if p.EdgeDetected() {
				if p.Read() == rpio.High {
					onClose() // HIGH = beam interrupted (robot present)
				} else {
					onOpen() // LOW = beam clear (robot absent)
				}
			}
		case <-exit:
			return nil
		}
	}
}
