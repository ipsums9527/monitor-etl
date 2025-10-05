package ros_http

import "resty.dev/v3"

type Client struct {
	*resty.Client
	ethers map[string]bool
}

type ethernet struct {
	Name            string `json:"name"`
	RxBitsPerSecond string `json:"rx-bits-per-second"`
	TxBitsPerSecond string `json:"tx-bits-per-second"`
}
