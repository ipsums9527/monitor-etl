package common

type SystemInfo struct {
	CpuLoad     string `json:"cpu-load"`     // cpu percentage used
	FreeMemory  string `json:"free-memory"`  // unit M
	TotalMemory string `json:"total-memory"` // unit M
	Uptime      string `json:"uptime"`       // form: 12w34d5h6m7s
}

type EtherInfo struct {
	Name            string
	RxBitsPerSecond float64
	TxBitsPerSecond float64
}

type Options struct {
	Host     string
	User     string
	Password string
	Ethers   []struct {
		Name     string
		IsInvert bool
	}
}

type HealthInfo struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}
