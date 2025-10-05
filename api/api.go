package api

import (
	"errors"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/labstack/gommon/log"

	"github.com/ipsums9527/monitor-etl/api/ikuai"
	"github.com/ipsums9527/monitor-etl/api/ros"
	"github.com/ipsums9527/monitor-etl/api/ros/common"
	"github.com/ipsums9527/monitor-etl/model"
)

func New(opt map[string]any) (*Client, error) {
	var cli SystemDataClient
	switch opt["type"].(string) {
	case "routeros":
		rosConfig := new(common.Options)
		if err := mapstructure.Decode(opt, rosConfig); err != nil {
			return nil, err
		}

		rosCli, err := ros.New(rosConfig)
		if err != nil {
			return nil, err
		}
		cli = rosCli

	case "ikuai":
		ikuaiConfig := new(ikuai.Options)
		if err := mapstructure.Decode(opt, ikuaiConfig); err != nil {
			return nil, err
		}

		ikuaiCli, err := ikuai.New(ikuaiConfig)
		if err != nil {
			return nil, err
		}
		cli = ikuaiCli

	default:
		return nil, errors.New("unknown connection type: " + opt["type"].(string))
	}

	return &Client{
		SendChan: make(chan *model.Message),
		cli:      cli,
		tk:       time.NewTicker(time.Second * 2),
	}, nil
}

func (c *Client) Start() {
	// run once immediately at startup
	c.run()

	for range c.tk.C {
		if c.isStop.Load() {
			return
		}
		c.run()
	}
}

func (c *Client) run() {
	if c.isStop.Load() {
		return
	}

	c.wg.Add(1)
	c.isRun.Store(true)
	defer func() {
		c.isRun.Store(false)
		c.wg.Done()
	}()

	// check again if stop
	if c.isStop.Load() {
		return
	}

	msg, err := c.cli.GetSystemDataMessage()
	if err != nil {
		log.Error(err)
		return
	}

	select {
	case c.SendChan <- msg:
	case <-time.After(time.Second * 5):
		log.Error("send message timeout")
	default:
		log.Error("send message failed")
	}

}

func (c *Client) Stop() error {
	if !c.isStop.CompareAndSwap(false, true) {
		return nil
	}

	c.tk.Stop()
	c.wg.Wait()
	close(c.SendChan)

	return c.cli.Close()
}
