package main

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/vkorn/go-miio"
)

const (
	maxRetries    = 60
	retryInterval = time.Minute
)

var errTimeout = errors.New("connection timed out")

// VacuumWatcher returns a WatchFunc that monitors a Roborock vacuum at the
// given IP and token, opening the door when cleaning starts and closing it
// when the vacuum returns to its charger. It retries automatically on timeout.
func VacuumWatcher(ip, token string) WatchFunc {
	return func(exit <-chan struct{}, onOpen, onClose func()) error {
		for attempt := 0; attempt < maxRetries; attempt++ {
			if attempt > 0 {
				slog.Info("Retrying connection", slog.Int("attempt", attempt))
				select {
				case <-time.After(retryInterval):
				case <-exit:
					return nil
				}
			}

			err := watchVacuum(ip, token, exit, onOpen, onClose)
			if err == nil {
				return nil
			}
			if !errors.Is(err, errTimeout) {
				return err
			}
			slog.Warn("Connection timed out", slog.Int("attempt", attempt+1))
		}
		return fmt.Errorf("max retries (%d) reached", maxRetries)
	}
}

func watchVacuum(ip, token string, exit <-chan struct{}, onOpen, onClose func()) error {
	vacuum, err := miio.NewVacuum(ip, token)
	if err != nil {
		return fmt.Errorf("connect to vacuum: %w", err)
	}

	var (
		vacState miio.VacState
		wg       sync.WaitGroup
		done     = make(chan struct{})
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case msg, ok := <-vacuum.UpdateChan:
				if !ok {
					return
				}
				s, ok := msg.State.(*miio.VacuumState)
				if !ok || vacState == s.State {
					continue
				}
				vacState = s.State

				if name, err := vacStateString(s.State); err != nil {
					slog.Error("Unknown vacuum state", slog.Any("error", err))
				} else {
					slog.Info("Vacuum state changed", slog.String("state", name))
				}

				switch s.State {
				case miio.VacStateCleaning:
					onOpen()
				case miio.VacStateCharging:
					onClose()
				}
			case <-done:
				return
			}
		}
	}()

	ticker := time.NewTicker(pollInterval)
	var runErr error
loop:
	for {
		select {
		case <-ticker.C:
			if !vacuum.UpdateStatus() {
				runErr = errTimeout
				break loop
			}
		case <-exit:
			break loop
		}
	}
	ticker.Stop()
	close(done)
	vacuum.Stop()
	wg.Wait()
	return runErr
}

func vacStateString(s miio.VacState) (string, error) {
	switch s {
	case miio.VacStateUnknown:
		return "unknown", nil
	case miio.VacStateInitiating:
		return "initiating", nil
	case miio.VacStateSleeping:
		return "sleeping", nil
	case miio.VacStateWaiting:
		return "waiting", nil
	case miio.VacStateCleaning:
		return "cleaning", nil
	case miio.VacStateReturning:
		return "returning", nil
	case miio.VacStateCharging:
		return "charging", nil
	case miio.VacStatePaused:
		return "paused", nil
	case miio.VacStateSpot:
		return "spot", nil
	case miio.VacStateShuttingDown:
		return "shutting down", nil
	case miio.VacStateUpdating:
		return "updating", nil
	case miio.VacStateDocking:
		return "docking", nil
	case miio.VacStateZone:
		return "zone", nil
	case miio.VacStateFull:
		return "full", nil
	default:
		return "", fmt.Errorf("unknown vacuum state: %d", s)
	}
}
