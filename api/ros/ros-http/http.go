package ros_http

import (
	"errors"
	"strconv"
	"strings"

	"resty.dev/v3"

	"github.com/ipsums9527/monitor-etl/api/ros/common"
)

func New(opt *common.Options) *Client {
	etherMaps := make(map[string]bool)
	for _, e := range opt.Ethers {
		etherMaps[e.Name] = e.IsInvert
	}

	return &Client{
		Client: resty.New().
			SetBaseURL(opt.Host+"/rest").
			SetBasicAuth(opt.User, opt.Password).
			SetDisableWarn(true).
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json"),
		ethers: etherMaps,
	}
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) GetSystemInfo() (*common.SystemInfo, error) {
	sysInfo := new(common.SystemInfo)
	resp, err := c.R().
		SetResult(sysInfo).
		SetQueryParam(".proplist", "cpu-load,free-memory,total-memory,uptime").
		Get("/system/resource")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	return sysInfo, nil
}

func (c *Client) GetTrafficInfo() (*common.EtherInfo, error) {
	var ethersInfo []ethernet
	allInterfacesName := ""
	for etherName := range c.ethers {
		allInterfacesName += etherName + ","
	}
	resp, err := c.R().SetResult(&ethersInfo).SetBody(map[string]string{
		"interface": allInterfacesName,
		"once":      "",
		".proplist": "name,rx-bits-per-second,tx-bits-per-second",
	}).Post("/interface/monitor-traffic")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	var totalTx, totalRx float64
	for _, e := range ethersInfo {
		rxBitsPerSecond, err := strconv.ParseFloat(e.RxBitsPerSecond, 64)
		if err != nil {
			return nil, err
		}
		txBitsPerSecond, err := strconv.ParseFloat(e.TxBitsPerSecond, 64)
		if err != nil {
			return nil, err
		}

		if c.ethers[e.Name] {
			totalTx += rxBitsPerSecond
			totalRx += txBitsPerSecond
		} else {
			totalTx += txBitsPerSecond
			totalRx += rxBitsPerSecond
		}
	}

	return &common.EtherInfo{
		Name:            allInterfacesName,
		RxBitsPerSecond: totalRx,
		TxBitsPerSecond: totalTx,
	}, nil
}

func (c *Client) GetHealthInfo() (*common.HealthInfo, error) {
	var healths []*common.HealthInfo
	resp, err := c.R().
		SetResult(&healths).
		Get("/system/health")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	var health *common.HealthInfo
	for _, e := range healths {
		if e.Name == "cpu-temperature" {
			return e, nil
		}

		if strings.Contains(e.Name, "temperature") {
			health = e
		}
	}
	if health != nil {
		return health, nil
	}

	return nil, common.ErrNotFoundTemperature
}
