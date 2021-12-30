package main

import (
	"github.com/kardianos/service"
)

type program struct {
	shell   shell
	service service.Service
}

func (s *program) Start(svc service.Service) error {
	log.Info("service '", svc.String(), "' started")

	go s.run()

	return nil
}

func (s *program) Stop(svc service.Service) error {
	log.Info("service '", svc.String(), "' stopped")

	return nil
}

func (s *program) Run() {
	if service.Interactive() {
		s.run()
	} else {
		if s.service != nil {
			s.service.Run()
		}
	}
}

func (s *program) run() {
	s.shell.Run()
}
