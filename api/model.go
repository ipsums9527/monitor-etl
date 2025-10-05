package api

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/ipsums9527/monitor-etl/model"
)

type Client struct {
	SendChan chan *model.Message

	cli    SystemDataClient
	tk     *time.Ticker
	isStop atomic.Bool
	isRun  atomic.Bool
	wg     sync.WaitGroup
}

type SystemDataClient interface {
	Close() error
	GetSystemDataMessage() (*model.Message, error)
}
