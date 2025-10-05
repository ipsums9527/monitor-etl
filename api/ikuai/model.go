package ikuai

import "resty.dev/v3"

type Client struct {
	cli      *resty.Client
	user     string
	password string
}

type Options struct {
	Host     string
	User     string
	Password string
}

type sysStat struct {
	Cpu        []string `json:"cpu"`
	Cputemp    []string `json:"cputemp"`
	Freq       []string `json:"freq"`
	Gwid       string   `json:"gwid"`
	Hostname   string   `json:"hostname"`
	LinkStatus int      `json:"link_status"`
	Memory     struct {
		Available int    `json:"available"`
		Buffers   int    `json:"buffers"`
		Cached    int    `json:"cached"`
		Free      int    `json:"free"`
		Total     int    `json:"total"`
		Used      string `json:"used"`
	} `json:"memory"`
	OnlineUser struct {
		Count         int `json:"count"`
		Count2G       int `json:"count_2g"`
		Count5G       int `json:"count_5g"`
		CountWired    int `json:"count_wired"`
		CountWireless int `json:"count_wireless"`
	} `json:"online_user"`
	Stream struct {
		ConnectNum int `json:"connect_num"`
		Download   int `json:"download"`
		TotalDown  int `json:"total_down"`
		TotalUp    int `json:"total_up"`
		Upload     int `json:"upload"`
	} `json:"stream"`
	Uptime  float64 `json:"uptime"`
	Verinfo struct {
		Arch            string `json:"arch"`
		Bootguide       string `json:"bootguide"`
		BuildDate       int64  `json:"build_date"`
		IsEnterprise    int    `json:"is_enterprise"`
		Modelname       string `json:"modelname"`
		SupportDingtalk int    `json:"support_dingtalk"`
		SupportI18N     int    `json:"support_i18n"`
		SupportLcd      int    `json:"support_lcd"`
		Sysbit          string `json:"sysbit"`
		Verflags        string `json:"verflags"`
		Version         string `json:"version"`
		Verstring       string `json:"verstring"`
	} `json:"verinfo"`
}
