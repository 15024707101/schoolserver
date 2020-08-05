// +build windows

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
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		for {
			sig := <-c
			logger.Info(sig)
			done <- true
		}
	}()

	<-done
	logger.Info("Stopping server...")
}
