package drv8825

import (
	"errors"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// defaultRPM is used when SetSpeed has not been called.
const defaultRPM uint = 60

// StepMode defines the microstepping resolution of the DRV8825.
// Each mode specifies the step angle and the M0/M1/M2 pin levels required
// to configure the driver hardware.
//
// DRV8825 truth table (M2, M1, M0):
//
//	0 0 0 → full step   (1.8°)
//	0 0 1 → half step   (0.9°)
//	0 1 0 → 1/4 step    (0.45°)
//	0 1 1 → 1/8 step    (0.225°)
//	1 0 0 → 1/16 step   (0.1125°)
//	1 0 1 → 1/32 step   (0.05625°)
//
// Values assume a standard 1.8°/step motor (e.g. NEMA 17, 200 steps/revolution).
type StepMode struct {
	degreePerStep float64
	m0, m1, m2   bool
}

var (
	StepModeFull         = StepMode{1.8, false, false, false}
	StepModeHalf         = StepMode{0.9, true, false, false}
	StepModeQuarter      = StepMode{0.45, false, true, false}
	StepModeEighth       = StepMode{0.225, true, true, false}
	StepModeSixteenth    = StepMode{0.1125, false, false, true}
	StepModeThirtySecond = StepMode{0.05625, true, false, true}
)

func (sm StepMode) stepsPerRev() uint {
	return uint(360 / sm.degreePerStep)
}

func (sm StepMode) stepDelay(rpm uint) time.Duration {
	// Each step consists of two phases (high + low), so divide by 2
	// to achieve the requested RPM accurately.
	return time.Minute / time.Duration(sm.stepsPerRev()*rpm*2)
}

// Config holds the GPIO pin numbers and operating mode for the DRV8825.
//
// PinMode0/1/2 correspond to the M0/M1/M2 inputs on the DRV8825 and control
// microstepping resolution. Set them to the GPIO pins wired to M0/M1/M2.
// If a pin is 0 the driver will not drive it — the step mode must then be
// set via hardware (jumpers / pull resistors) and must match StepMode.
type Config struct {
	PinEnable    uint8
	PinStep      uint8
	PinDirection uint8
	PinMode0     uint8 // M0 — optional, 0 = hardware-controlled
	PinMode1     uint8 // M1 — optional, 0 = hardware-controlled
	PinMode2     uint8 // M2 — optional, 0 = hardware-controlled
	StepMode     StepMode
	RPM          uint
}

type Driver struct {
	en       rpio.Pin
	step     rpio.Pin
	dir      rpio.Pin
	stepMode StepMode
	rpm      uint
}

func New(cfg Config) (*Driver, error) {
	if cfg.StepMode.degreePerStep == 0 {
		return nil, errors.New("drv8825: StepMode must be set")
	}

	rpm := cfg.RPM
	if rpm == 0 {
		rpm = defaultRPM
	}

	d := &Driver{
		en:       rpio.Pin(cfg.PinEnable),
		step:     rpio.Pin(cfg.PinStep),
		dir:      rpio.Pin(cfg.PinDirection),
		stepMode: cfg.StepMode,
		rpm:      rpm,
	}

	d.en.Output()
	d.step.Output()
	d.dir.Output()
	d.en.High() // ENABLE is active-LOW; start disabled to avoid energizing the coils
	d.step.Low()
	d.dir.Low()

	// Configure microstepping mode pins if provided.
	setPinLevel(cfg.PinMode0, cfg.StepMode.m0)
	setPinLevel(cfg.PinMode1, cfg.StepMode.m1)
	setPinLevel(cfg.PinMode2, cfg.StepMode.m2)

	return d, nil
}

// setPinLevel drives pin to the given logic level. A pin value of 0 is a no-op
// (the pin is not connected to the Pi and is hardware-controlled instead).
func setPinLevel(pin uint8, high bool) {
	if pin == 0 {
		return
	}
	p := rpio.Pin(pin)
	p.Output()
	if high {
		p.High()
	} else {
		p.Low()
	}
}

func (d *Driver) SetSpeed(rpm uint) {
	if rpm == 0 {
		panic("drv8825: rpm must be greater than zero")
	}
	d.rpm = rpm
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

	// tDIR: DRV8825 requires DIR to be stable ≥650 ns before the rising STEP edge.
	time.Sleep(time.Microsecond)

	delay := d.stepMode.stepDelay(d.rpm)
	for i := 0; i < steps; i++ {
		d.step.High()
		time.Sleep(delay)
		d.step.Low()
		time.Sleep(delay)
	}
}

