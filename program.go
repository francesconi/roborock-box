package main

import (
	"fmt"
	"time"

	"github.com/kardianos/service"
	"github.com/vkorn/go-miio"
)

type program struct {
	garage   *Garage
	vacuum   *miio.Vacuum
	vacState miio.VacState
	logger   service.Logger
	exit     chan struct{}
}

func NewProgram() (*program, error) {
	garage, err := NewGarage()
	if err != nil {
		return nil, err
	}

	ip := ""
	token := ""
	vacuum, err := miio.NewVacuum(ip, token)
	if err != nil {
		defer garage.Cleanup()
		return nil, err
	}

	return &program{
		exit:   make(chan struct{}),
		garage: garage,
		vacuum: vacuum,
	}, nil
}

func (p program) Start(s service.Service) error {
	l, err := s.Logger(nil)
	if err != nil {
		return err
	}
	p.logger = l

	go p.run()
	return nil
}

func (p program) run() error {
	go func() {
		for msg := range p.vacuum.UpdateChan {
			s, ok := msg.State.(*miio.VacuumState)
			if !ok || p.vacState == s.State {
				continue
			}
			p.vacState = s.State

			vacStateStr, err := vacStateString(s.State)
			if err != nil {
				p.logger.Error(err)
			} else {
				p.logger.Infof("Got state %s", vacStateStr)
			}

			switch s.State {
			case miio.VacStateCleaning:
				p.garage.OpenDoor()
			case miio.VacStateCharging:
				p.garage.CloseDoor()
			}
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			p.vacuum.UpdateStatus()
		case <-p.exit:
			ticker.Stop()
			p.garage.Cleanup()
			return nil
		}
	}
}

func (p program) Stop(s service.Service) error {
	close(p.exit)
	return nil
}

func (p program) Shutdown(s service.Service) error {
	close(p.exit)
	return nil
}

func vacStateString(s miio.VacState) (string, error) {
	switch s {
	case miio.VacStateUnknown:
		return "VacStateUnknown", nil
	case miio.VacStateInitiating:
		return "VacStateInitiating", nil
	case miio.VacStateSleeping:
		return "VacStateSleeping", nil
	case miio.VacStateWaiting:
		return "VacStateWaiting", nil
	case miio.VacStateCleaning:
		return "VacStateCleaning", nil
	case miio.VacStateReturning:
		return "VacStateReturning", nil
	case miio.VacStateCharging:
		return "VacStateCharging", nil
	case miio.VacStatePaused:
		return "VacStatePaused", nil
	case miio.VacStateSpot:
		return "VacStateSpot", nil
	case miio.VacStateShuttingDown:
		return "VacStateShuttingDown", nil
	case miio.VacStateUpdating:
		return "VacStateUpdating", nil
	case miio.VacStateDocking:
		return "VacStateDocking", nil
	case miio.VacStateZone:
		return "VacStateZone", nil
	case miio.VacStateFull:
		return "VacStateFull", nil
	default:
		return "", fmt.Errorf("miio: unknown VacState: %d", s)
	}
}
