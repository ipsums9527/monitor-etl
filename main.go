package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/gommon/log"

	"github.com/ipsums9527/monitor-etl/config"
	"github.com/ipsums9527/monitor-etl/control"
)

func main() {
	config.ShowVersion()

	conf := config.ReadConfig()
	c, err := control.New(conf)
	if err != nil {
		panic(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := c.Start(); err != nil {
			log.Error("start error: ", err)
			sigChan <- syscall.SIGTERM
		}
	}()

	log.Infof("received signal: %v, shutting down gracefully...", <-sigChan)
	if err := c.Stop(); err != nil {
		log.Error("stop error:", err)
		os.Exit(1)
	}

	log.Info("shutdown complete")
}
