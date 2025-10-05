package ros_tcp

import (
	"github.com/go-routeros/routeros/v3"
)

type Client struct {
	address  string
	user     string
	password string
	ethers   map[string]bool
	cli      *routeros.Client
}
