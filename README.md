# Roborock Box

A Raspberry Pi service that automatically opens and closes a motorised box door based on your Roborock vacuum's state — the door opens when the robot starts cleaning and closes when it returns to the charger.

Alternatively, a light sensor can be used instead of the Wi-Fi vacuum integration.

## Hardware

| Part | Notes |
|---|---|
| Raspberry Pi | Any model with GPIO |
| NEMA 17 Stepper Motor | 1.8°/step (200 steps/revolution) |
| DRV8825 Stepper Motor Driver | |
| Stepper Motor Controller Board | |

### Wiring

| DRV8825 | Raspberry Pi GPIO |
|---|---|
| ENABLE | 24 |
| STEP | 23 |
| DIR | 22 |

M0, M1, M2 (microstepping mode pins) are hardware-controlled via jumpers. Tie all three to GND for full-step mode (default).

For the optional light sensor:

| Sensor | Raspberry Pi GPIO |
|---|---|
| DO | 25 |

## Prerequisites

- Go 1.18+
- SSH access to the Raspberry Pi
- `make`

## Configuration

### 1. Obtain the vacuum token

See the [python-miio docs](https://python-miio.readthedocs.io/en/latest/discovery.html#obtaining-tokens) for instructions on retrieving your vacuum's local IP address and token.

### 2. Edit hardware config

Pin numbers, step mode, speed, and door travel are set in [`main.go`](main.go):

```go
p := newProgram(VacuumWatcher(os.Getenv("IP"), os.Getenv("TOKEN")), BoxConfig{
    PinEnable:    24,
    PinStep:      23,
    PinDirection: 22,
    StepMode:     drv8825.StepModeFull,
    RPM:          300,
    DoorSteps:    4000,
})
```

Adjust `RPM` and `DoorSteps` to match your physical setup.

### 3. Deploy

```sh
make deploy
```

This builds the binary for ARM, copies it to the Pi, registers it as a system service, and starts it.

### 4. Set credentials

On first deploy, set the vacuum's IP and token on the Pi:

```sh
ssh pi@<host> sudo nano /etc/default/roborock-box
```

```sh
IP=192.168.1.x
TOKEN=your_token_here
```

Then restart the service:

```sh
ssh pi@<host> sudo service roborock-box restart
```

## Service management

```sh
sudo service roborock-box status
sudo service roborock-box start
sudo service roborock-box stop
sudo service roborock-box restart
```

To uninstall:

```sh
make uninstall
```

## Using the light sensor

To use a light sensor instead of the vacuum Wi-Fi integration, swap the watcher in [`main.go`](main.go):

```go
// replace this:
p := newProgram(VacuumWatcher(os.Getenv("IP"), os.Getenv("TOKEN")), cfg)

// with this:
p := newProgram(LightSensorWatcher(25), cfg)
```

No credentials are required. The door closes when the sensor beam is interrupted (robot present) and opens when the beam is clear (robot absent).

## Credits

- [makerportal/nema17-python](https://github.com/makerportal/nema17-python)
- [gavinlyonsrepo/RpiMotorLib](https://github.com/gavinlyonsrepo/RpiMotorLib/blob/master/Documentation/Nema11DRV8825.md)
- [shanghuiyang/rpi-devices](https://github.com/shanghuiyang/rpi-devices)
- [dimschlukas/rpi_python_drv8825](https://github.com/dimschlukas/rpi_python_drv8825)
