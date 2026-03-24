package main

import (
	"bufio"
	"log/slog"
	"os"
	"strings"

	"github.com/francesconi/roborock-box/drv8825"
	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
)

const envFile = "/etc/default/roborock-box"

func main() {
	// Suppress verbose logging from the miio library
	logrus.SetLevel(logrus.InfoLevel)

	// Load IP and TOKEN from the env file so the service works regardless of
	// whether the init system sources it automatically.
	loadEnvFile(envFile)

	svcConfig := &service.Config{
		Name:        "roborock-box",
		DisplayName: "Roborock Box",
		Description: "Automatically opens and closes the robot vacuum box.",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target",
		},
	}

	ip, token := os.Getenv("IP"), os.Getenv("TOKEN")

	p := newProgram(VacuumWatcher(ip, token), BoxConfig{
		PinEnable:    24,
		PinStep:      23,
		PinDirection: 22,
		StepMode:     drv8825.StepModeFull,
		RPM:          300,
		DoorSteps:    4000,
	})

	svc, err := service.New(p, svcConfig)
	if err != nil {
		slog.Error("Failed to create service", slog.Any("error", err))
		os.Exit(1)
	}

	// Support: install, uninstall, start, stop, restart, status
	// e.g. sudo service roborock-box start
	if len(os.Args) > 1 {
		if err := service.Control(svc, os.Args[1]); err != nil {
			slog.Error("Service control failed", slog.String("command", os.Args[1]), slog.Any("error", err))
			os.Exit(1)
		}
		return
	}

	if ip == "" || token == "" {
		slog.Error("IP and TOKEN must be set in " + envFile)
		os.Exit(1)
	}

	if err := svc.Run(); err != nil {
		slog.Error("Service exited with error", slog.Any("error", err))
		os.Exit(1)
	}
}

// loadEnvFile reads KEY=VALUE pairs from path and sets them as environment
// variables. Existing variables are not overwritten, so values set by the
// calling environment always take precedence. Lines starting with # are ignored.
func loadEnvFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // file not found is expected when running outside the Pi
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok || os.Getenv(key) != "" {
			continue
		}
		os.Setenv(key, value)
	}
}
