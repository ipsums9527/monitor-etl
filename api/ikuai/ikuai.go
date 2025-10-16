package ikuai

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/labstack/gommon/log"
	"resty.dev/v3"

	"github.com/ipsums9527/monitor-etl/model"
)

func New(opt *Options) (*Client, error) {
	u, err := url.Parse(opt.Host)
	if err != nil {
		return nil, err
	}
	u.Path = "/Action"
	c := &Client{
		user:     opt.User,
		password: opt.Password,
	}

	httpCli := resty.New().
		SetBaseURL(u.String()).
		SetDisableWarn(true).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetRetryCount(3).AddRetryConditions(func(response *resty.Response, err error) bool {
		js, err := simplejson.NewJson(response.Bytes())
		if err != nil {
			log.Error(err)
			return true
		}
		if js.Get("ErrMsg").MustString() == "no login authentication" {
			if err := c.login(); err != nil {
				log.Error(err)
				return true
			}
			response.Request.SetCookies(c.cli.Cookies())
			return true
		}
		return false
	}).SetAllowNonIdempotentRetry(true)
	c.cli = httpCli

	return c, nil
}

func (c *Client) login() error {
	log.Info("start login")
	var rtn = &struct {
		ErrMsg string `json:"ErrMsg"`
	}{}

	resp, err := resty.New().SetBaseURL(c.cli.BaseURL()).R().SetResult(rtn).SetBody(map[string]string{
		"username": c.user,
		"passwd":   fmt.Sprintf("%x", md5.Sum([]byte(c.password))),
	}).Post("/login")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return errors.New(resp.String())
	}

	if rtn.ErrMsg != "Success" {
		return errors.New(rtn.ErrMsg)
	}
	c.cli.SetCookies(resp.Cookies())

	return nil
}

func (c *Client) getSysStat() (*sysStat, error) {
	if len(c.cli.Cookies()) == 0 {
		if err := c.login(); err != nil {
			return nil, err
		}
	}

	var rtn = &struct {
		ErrMsg string `json:"ErrMsg"`
		Data   struct {
			SysStat *sysStat `json:"sysstat"`
		} `json:"Data"`
	}{}
	resp, err := c.cli.R().SetBody(map[string]any{
		"func_name": "homepage",
		"action":    "show",
		"param":     map[string]string{"TYPE": "sysstat"},
	}).Post("/call")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	if err := json.Unmarshal(resp.Bytes(), rtn); err != nil {
		return nil, err
	}
	if rtn.ErrMsg != "Success" {
		return nil, errors.New(rtn.ErrMsg)
	}

	return rtn.Data.SysStat, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) GetSystemDataMessage() (*model.Message, error) {
	stat, err := c.getSysStat()
	if err != nil {
		return nil, err
	}

	msg := new(model.Message)
	cpu, _ := strconv.ParseFloat(strings.ReplaceAll(stat.Cpu[0], "%", ""), 64)
	msg.CPU.System = cpu
	msg.Uptime = stat.Uptime
	msg.RAM.Free = float64(stat.Memory.Free) / 1024
	msg.RAM.Used = float64(stat.Memory.Total-stat.Memory.Free) / 1024
	msg.Net.Upload = float64(stat.Stream.Upload) / 1024 * 8
	msg.Net.Download = float64(stat.Stream.Download) / 1024 * 8

	if len(stat.Cputemp) > 0 {
		msg.Temp, _ = strconv.ParseFloat(stat.Cputemp[0], 64)
	} else {
		msg.Temp = msg.CPU.System/100*30 + 40
	}

	return msg, nil
}
