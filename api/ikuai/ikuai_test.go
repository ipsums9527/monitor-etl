package ikuai

import (
	"testing"

	"github.com/go-viper/mapstructure/v2"

	"github.com/ipsums9527/monitor-etl/config"
)

func TestLogin(t *testing.T) {
	opt := new(Options)
	if err := mapstructure.Decode(config.ReadConfig().Api, opt); err != nil {
		t.Fatal(err)
	}
	c, err := New(opt)
	if err != nil {
		t.Error(err)
	}

	if err := c.login(); err != nil {
		t.Error(err)
		return
	}

	t.Log(c.cli.Cookies())
}

func TestFetch(t *testing.T) {
	opt := new(Options)
	if err := mapstructure.Decode(config.ReadConfig().Api, opt); err != nil {
		t.Fatal(err)
	}
	c, err := New(opt)
	if err != nil {
		t.Error(err)
	}

	if err := c.login(); err != nil {
		t.Error(err)
		return
	}

	t.Log(c.getSysStat())
}
