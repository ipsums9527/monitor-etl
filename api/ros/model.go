package ros

import (
	"github.com/ipsums9527/monitor-etl/api/ros/common"
)

type rosClient interface {
	Close() error
	GetSystemInfo() (*common.SystemInfo, error)

	// GetTrafficInfo get traffic info, return tx bit/s, rx bit/s, err
	GetTrafficInfo() (*common.EtherInfo, error)
	GetHealthInfo() (*common.HealthInfo, error)
}

type Client struct {
	rosClient
}
