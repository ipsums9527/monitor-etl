package ros

import (
	"errors"
	"net/url"
	"strconv"
	"unicode"

	"github.com/labstack/gommon/log"

	"github.com/ipsums9527/monitor-etl/api/ros/common"
	"github.com/ipsums9527/monitor-etl/api/ros/ros-http"
	"github.com/ipsums9527/monitor-etl/api/ros/ros-tcp"
	"github.com/ipsums9527/monitor-etl/model"
)

func New(opt *common.Options) (*Client, error) {
	u, err := url.Parse(opt.Host)
	if err != nil {
		return nil, err
	}
	opt.Host = u.Host

	var client rosClient
	switch u.Scheme {
	case "http":
		client = ros_http.New(opt)

	case "tcp":
		tcpClient, err := ros_tcp.New(opt)
		if err != nil {
			return nil, err
		}
		client = tcpClient

	default:
		return nil, errors.New("unknown connection type: " + u.Scheme)
	}

	c := &Client{
		rosClient: client,
	}

	return c, nil
}

func (c *Client) Close() error {
	return c.rosClient.Close()
}

func (c *Client) GetSystemDataMessage() (*model.Message, error) {
	var errs error
	msg := new(model.Message)

	// get traffic info
	ei, err := c.GetTrafficInfo()
	if err != nil {
		errs = errors.Join(err)
	} else {
		msg.Net.Upload = ei.TxBitsPerSecond / 1024
		msg.Net.Download = ei.RxBitsPerSecond / 1024
	}

	// get system info
	sys, err := c.GetSystemInfo()
	if err != nil {
		errs = errors.Join(err)
	} else {
		load, _ := strconv.ParseFloat(sys.CpuLoad, 64)
		msg.CPU.System = load
		memFree, _ := strconv.ParseFloat(sys.FreeMemory, 64)
		totalMem, _ := strconv.ParseFloat(sys.TotalMemory, 64)
		msg.RAM.Free = memFree / 1024 / 1024
		msg.RAM.Used = (totalMem - memFree) / 1024 / 1024
		msg.Uptime = formatUptime(sys.Uptime)
	}

	//get temperature
	health, err := c.GetHealthInfo()
	if err != nil {
		if errors.Is(err, common.ErrNotFoundTemperature) {
			msg.Temp = msg.CPU.System/100*30 + 40
		} else {
			errs = errors.Join(err)
		}
	} else {
		msg.Temp, _ = strconv.ParseFloat(health.Value, 64)
	}

	if errs != nil {
		return nil, errs
	}

	return msg, nil
}

func formatUptime(timeStr string) float64 {
	n := ""
	uptime := 0
	for _, s := range timeStr {
		if unicode.IsDigit(s) {
			n += string(s)
		} else {
			switch s {
			case 'w':
				digit, _ := strconv.Atoi(n)
				uptime += digit * 60 * 60 * 24 * 7
			case 'd':
				digit, _ := strconv.Atoi(n)
				uptime += digit * 60 * 60 * 24
			case 'h':
				digit, _ := strconv.Atoi(n)
				uptime += digit * 60 * 60
			case 'm':
				digit, _ := strconv.Atoi(n)
				uptime += digit * 60
			case 's':
				digit, _ := strconv.Atoi(n)
				uptime += digit
			default:
				log.Errorf("unknown uptime unit: <%s>", timeStr)
				return 0
			}
			n = ""
		}
	}

	return float64(uptime)
}
