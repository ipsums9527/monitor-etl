package ros_tcp

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-routeros/routeros/v3"
	"github.com/go-viper/mapstructure/v2"
	"github.com/labstack/gommon/log"

	"github.com/ipsums9527/monitor-etl/api/ros/common"
)

func New(opt *common.Options) (*Client, error) {
	cli, err := routeros.Dial(opt.Host, opt.User, opt.Password)
	if err != nil {
		return nil, err
	}

	etherMaps := make(map[string]bool)
	for _, e := range opt.Ethers {
		etherMaps[e.Name] = e.IsInvert
	}

	return &Client{
		address:  opt.Host,
		user:     opt.User,
		password: opt.Password,
		ethers:   etherMaps,
		cli:      cli,
	}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) GetSystemInfo() (*common.SystemInfo, error) {
	re, err := c.cli.RunArgs([]string{"/system/resource/print", ".proplist=cpu-load,free-memory,total-memory,uptime"})
	if err != nil {
		c.reconnect()
		return nil, err
	}

	return &common.SystemInfo{
		CpuLoad:     re.Re[0].Map["cpu-load"],
		FreeMemory:  re.Re[0].Map["free-memory"],
		TotalMemory: re.Re[0].Map["total-memory"],
		Uptime:      re.Re[0].Map["uptime"],
	}, nil
}

func (c *Client) GetTrafficInfo() (*common.EtherInfo, error) {
	allInterfacesName := ""
	for etherName := range c.ethers {
		allInterfacesName += etherName + ","
	}
	re, err := c.cli.RunArgs([]string{
		"/interface/monitor-traffic",
		"=.proplist=name,rx-bits-per-second,tx-bits-per-second",
		"=interface=" + allInterfacesName,
		"=once=",
	})
	if err != nil {
		c.reconnect()
		return nil, err
	}

	var totalRx, totalTx float64
	for _, e := range re.Re {
		rx, err := strconv.ParseFloat(e.Map["rx-bits-per-second"], 64)
		if err != nil {
			return nil, err
		}
		tx, err := strconv.ParseFloat(e.Map["tx-bits-per-second"], 64)
		if err != nil {
			return nil, err
		}

		if c.ethers[e.Map["name"]] {
			tx, rx = rx, tx
		}
		totalRx += rx
		totalTx += tx
	}

	return &common.EtherInfo{
		Name:            allInterfacesName,
		RxBitsPerSecond: totalRx,
		TxBitsPerSecond: totalTx,
	}, nil
}

func (c *Client) GetHealthInfo() (*common.HealthInfo, error) {
	re, err := c.cli.RunArgs([]string{
		"/system/health/print",
	})
	if err != nil {
		c.reconnect()
		return nil, err
	}

	health := new(common.HealthInfo)
	var isFound bool
	for _, e := range re.Re {
		if e.Map["name"] == "cpu-temperature" {
			if err := mapstructure.Decode(e.Map, health); err != nil {
				log.Error(err)
			}
			return health, nil
		}
		if strings.Contains(e.Map["name"], "temperature") {
			isFound = true
			if err := mapstructure.Decode(e.Map, health); err != nil {
				log.Error(err)
			}
		}
	}
	if isFound {
		return health, nil
	}

	return nil, common.ErrNotFoundTemperature
}

func (c *Client) reconnect() {
	defer time.Sleep(time.Second * 5)

	log.Info("reconnecting...")
	if err := c.cli.Close(); err != nil {
		log.Error(err)
	}
	cli, err := routeros.Dial(c.address, c.user, c.password)
	if err != nil {
		log.Error(err)
		return
	}
	if err := c.cli.Close(); err != nil {
		log.Error(err)
	}
	c.cli = cli

	return
}
