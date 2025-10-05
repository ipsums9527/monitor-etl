package control

import (
	"errors"

	"github.com/labstack/gommon/log"

	"github.com/ipsums9527/monitor-etl/api"
	"github.com/ipsums9527/monitor-etl/app/server"
	"github.com/ipsums9527/monitor-etl/config"
)

func New(conf *config.Config) (*Control, error) {
	srv := server.New(&server.Options{
		Listen: conf.Listen,
		Port:   conf.Port,
	})

	cli, err := api.New(conf.Api)
	if err != nil {
		return nil, err
	}

	return &Control{
		srv: srv,
		cli: cli,
	}, nil
}

func (c *Control) Start() error {
	log.Info("start service")
	go c.cli.Start()
	go func() {
		for msg := range c.cli.SendChan {
			c.srv.ReceiveChan <- msg
		}
	}()
	return c.srv.Start()
}

func (c *Control) Stop() error {
	log.Info("stop service")
	return errors.Join(c.cli.Stop(), c.srv.Stop())
}
