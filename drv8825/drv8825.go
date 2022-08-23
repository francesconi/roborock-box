package drv8825

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type StepMode uint8

const (
	StepModeFull StepMode = iota + 1
	StepModeHalf
	StepModeQuarter
	StepModeEighth
	StepModeSixteenth
	StepModeThirtySecond
)

func (sm StepMode) degreePerStep() float64 {
	switch sm {
	default:
		fallthrough
	case StepModeFull:
		return 1.8
	case StepModeHalf:
		return 0.9
	case StepModeQuarter:
		return 0.45
	case StepModeEighth:
		return 0.225
	case StepModeSixteenth:
		return 0.1125
	case StepModeThirtySecond:
		return 0.05625
	}
}

func (sm StepMode) stepsPerRevolution() uint {
	return uint(360 / sm.degreePerStep())
}

func (sm StepMode) delayPerStep(rpm uint) time.Duration {
	return time.Minute / time.Duration(sm.stepsPerRevolution()*rpm)
}

type Config struct {
	PinEnable    uint8
	PinStep      uint8
	PinDirection uint8
	StepMode     StepMode
}

type Driver struct {
	en       rpio.Pin
	step     rpio.Pin
	dir      rpio.Pin
	stepMode StepMode
	rpm      uint
}

func New(cfg Config) (*Driver, error) {
	err := rpio.Open()
	if err != nil {
		return nil, err
	}

	d := &Driver{
		en:       rpio.Pin(cfg.PinEnable),
		step:     rpio.Pin(cfg.PinStep),
		dir:      rpio.Pin(cfg.PinDirection),
		stepMode: cfg.StepMode,
	}

	d.en.Output()
	d.step.Output()
	d.dir.Output()
	d.step.Low()
	d.dir.Low()

	return d, nil
}

func (d *Driver) Enable() {
	d.en.Low()
}

func (d *Driver) Disable() {
	d.en.High()
}

func (d *Driver) Move(steps int) {
	if steps == 0 {
		return
	}
	if steps > 0 {
		d.dir.High()
	} else {
		d.dir.Low()
		steps = -steps
	}
	delay := d.stepMode.delayPerStep(d.rpm)
	for i := 0; i < steps; i++ {
		d.step.High()
		time.Sleep(delay)
		d.step.Low()
		time.Sleep(delay)
	}
}

func (d *Driver) Close() error {
	return rpio.Close()
}

func (d *Driver) SetSpeed(rpm uint) {
	d.rpm = rpm
}
