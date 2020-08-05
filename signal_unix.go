// +build darwin freebsd linux

package main

import (
	"schoolserver/logger"
	"os"
	"os/signal"
	"syscall"
)

func WaitForSignals() {
	c := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

	go func() {
		for {
			sig := <-c
			if sig == syscall.SIGUSR1 {
				logger.Info("SIGUSR1")

			} else {
				done <- true
			}
		}
	}()

	<-done
	logger.Info("Stopping server...")
}
