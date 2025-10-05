package server

import (
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/ipsums9527/monitor-etl/model"
)

type Server struct {
	Listen      string
	Port        int
	ReceiveChan chan *model.Message

	e    *echo.Echo
	data *sysData
	tk   *time.Ticker
}

type Options struct {
	Listen string
	Port   int
}

type sysData struct {
	memFree atomic.Value
	memUsed atomic.Value
	uptime  atomic.Value
	cpuUser atomic.Value
	cpuSys  atomic.Value
	netDown atomic.Value
	netUp   atomic.Value
}
