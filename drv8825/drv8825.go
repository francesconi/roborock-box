package drv8825

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type Driver struct {
	en            rpio.Pin
	step          rpio.Pin
	dir           rpio.Pin
	ms1, ms2, ms3 rpio.Pin
	mode          StepMode
}

type StepMode int

const (
	_ StepMode = iota
	StepModeFull
	StepModeHalf
	StepModeQuarter
	StepModeEighth
	StepModeSixteenth
	StepModeThirtySecond
)

type Direction int

const (
	_ Direction = iota
	DirectionForward
	DirectionBackward
)

type Config struct {
	EN            uint8
	STEP          uint8
	DIR           uint8
	MS1, MS2, MS3 uint8
	Mode          StepMode
}

func New(cfg Config) (*Driver, error) {
	err := rpio.Open()
	if err != nil {
		return nil, err
	}

	d := &Driver{
		en:   rpio.Pin(cfg.EN),
		step: rpio.Pin(cfg.STEP),
		dir:  rpio.Pin(cfg.DIR),
		ms1:  rpio.Pin(cfg.MS1),
		ms2:  rpio.Pin(cfg.MS2),
		ms3:  rpio.Pin(cfg.MS3),
	}

	d.en.Output()
	d.step.Output()
	d.dir.Output()
	d.ms1.Output()
	d.ms2.Output()
	d.ms3.Output()

	if err = d.SetMode(cfg.Mode); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *Driver) Enable() {
	d.en.Low()
}

func (d *Driver) Disable() {
	d.en.High()
}

func (d *Driver) Forward(steps int) error {
	return d.Move(steps, DirectionForward)
}

func (d *Driver) Backward(steps int) error {
	return d.Move(steps, DirectionBackward)
}

func (d *Driver) Move(steps int, direction Direction) error {
	switch direction {
	case DirectionForward:
		d.dir.High()
	case DirectionBackward:
		d.dir.Low()
	default:
		return fmt.Errorf("drv8825: invalid direction %d", direction)
	}

	delay := 500 * time.Microsecond
	for i := 0; i < steps; i++ {
		d.step.High()
		time.Sleep(delay)
		d.step.Low()
		time.Sleep(delay)
	}

	return nil
}

func (d *Driver) Stop() error {
	if err := rpio.Close(); err != nil {
		return err
	}
	return nil
}

func (d *Driver) SetMode(mode StepMode) error {
	switch mode {
	case StepModeFull:
		d.ms1.Low()
		d.ms2.Low()
		d.ms3.Low()
	case StepModeHalf:
		d.ms1.High()
		d.ms2.Low()
		d.ms3.Low()
	case StepModeQuarter:
		d.ms1.Low()
		d.ms2.High()
		d.ms3.Low()
	case StepModeEighth:
		d.ms1.High()
		d.ms2.High()
		d.ms3.Low()
	case StepModeSixteenth:
		d.ms1.Low()
		d.ms2.Low()
		d.ms3.High()
	case StepModeThirtySecond:
		d.ms1.High()
		d.ms2.Low()
		d.ms3.High()
	default:
		return fmt.Errorf("drv8825: invalid format %d", mode)
	}
	d.mode = mode
	return nil
}
