package server

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/ipsums9527/monitor-etl/model"
)

func New(opt *Options) *Server {
	e := echo.New()
	//e.Use(middleware.Recover())

	s := &Server{
		Listen:      opt.Listen,
		Port:        opt.Port,
		ReceiveChan: make(chan *model.Message),
		e:           e,
		data:        new(sysData),
		tk:          time.NewTicker(time.Minute),
	}

	e.GET("/api/v1/data", s.Data)

	return s
}

func (s *Server) Start() error {
	// after get the first message, then start the web server
	s.storeMsg(<-s.ReceiveChan)

	go func() {
		for range s.tk.C {
			log.Infof("uptime: %s, cpu: %v%%, up: %.2fMbps, down: %.2fMbps",
				time.Duration(int64(s.data.uptime.Load().(float64)))*time.Second,
				s.data.cpuSys.Load(),
				s.data.netUp.Load().(float64)/1024, s.data.netDown.Load().(float64)/1024)
		}
	}()

	go s.getData()

	return s.e.Start(s.Listen + ":" + strconv.Itoa(s.Port))
}

func (s *Server) Stop() error {
	close(s.ReceiveChan)
	s.tk.Stop()
	return s.e.Close()
}

func (s *Server) getData() {
	for msg := range s.ReceiveChan {
		s.storeMsg(msg)
	}
}

func (s *Server) storeMsg(msg *model.Message) {
	s.data.memFree.Store(msg.RAM.Free)
	s.data.memUsed.Store(msg.RAM.Used)
	s.data.cpuSys.Store(msg.CPU.System)
	s.data.cpuUser.Store(msg.CPU.User)
	s.data.uptime.Store(msg.Uptime)
	s.data.netDown.Store(msg.Net.Download)
	s.data.netUp.Store(msg.Net.Upload)
	s.data.temp.Store(msg.Temp)
}
