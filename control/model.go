package control

import (
	"github.com/ipsums9527/monitor-etl/api"
	"github.com/ipsums9527/monitor-etl/app/server"
)

type Options struct {
	Server struct {
		Listen string
		Port   int
	}
	Client struct {
		Host      string
		EtherName string
		IsInvert  bool
		User      string
		Password  string
	}
}
type Control struct {
	srv *server.Server
	cli *api.Client
}
