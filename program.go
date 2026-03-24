package main

import (
	"log/slog"
	"sync"

	"github.com/kardianos/service"
)

type program struct {
	watch  WatchFunc
	boxCfg BoxConfig
	once   sync.Once
	exit   chan struct{}
}

func newProgram(watch WatchFunc, boxCfg BoxConfig) *program {
	return &program{
		watch:  watch,
		boxCfg: boxCfg,
		exit:   make(chan struct{}),
	}
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	p.once.Do(func() { close(p.exit) })
	return nil
}

func (p *program) Shutdown(s service.Service) error {
	return p.Stop(s)
}

func (p *program) run() {
	box, err := NewBox(p.boxCfg)
	if err != nil {
		slog.Error("Failed to initialize box", slog.Any("error", err))
		return
	}
	defer func() {
		if err := box.Cleanup(); err != nil {
			slog.Error("Box cleanup failed", slog.Any("error", err))
		}
	}()

	if err := p.watch(p.exit, box.OpenDoor, box.CloseDoor); err != nil {
		slog.Error("Watcher exited with error", slog.Any("error", err))
	}
}
