package api

import (
	"testing"

	"github.com/ipsums9527/monitor-etl/config"
)

func TestFetch(t *testing.T) {
	c, err := New(config.ReadConfig().Api)
	if err != nil {
		t.Error(err)
	}

	t.Log(c.cli.GetSystemDataMessage())
}
