package main

import "github.com/takama/daemon"

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// StartDaemon start daemon
func (service *Service) StartDaemon() error {
	if _, err := service.Install(); err != nil {
		return err
	}
	if _, err := service.Start(); err != nil {
		return err
	}
	return nil
}

// StopDaemon stop daemon
func (service *Service) StopDaemon() error {
	if _, err := service.Stop(); err != nil {
		return err
	}
	if _, err := service.Remove(); err != nil {
		return err
	}
	return nil
}

// StatusDaemon daemon status
func (service *Service) StatusDaemon() error {
	if _, err := service.Status(); err != nil {
		return err
	}
	return nil
}
