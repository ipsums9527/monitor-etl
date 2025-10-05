package model

type Message struct {
	RAM struct {
		Free float64 `json:"free,omitempty"` // unit Mb
		Used float64 `json:"used,omitempty"` // unit Mb
	} `json:"ram,omitempty"`
	Uptime float64 `json:"uptime,omitempty"` // unit second
	CPU    struct {
		User   float64 `json:"user,omitempty"`   // percentage
		System float64 `json:"system,omitempty"` // percentage
	} `json:"cpu,omitempty"`
	Net struct {
		Download float64 `json:"download,omitempty"` // unit kbit/s
		Upload   float64 `json:"upload,omitempty"`   // unit kbit/s
	} `json:"net,omitempty"`
}
