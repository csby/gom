package model

type MonitorNetworkIOArgument struct {
	Name string `json:"name" note:"网卡名称"`
}

type MonitorNetworkIO struct {
	Name     string `json:"name" note:"网卡名称"`
	MaxCount int    `json:"maxCount" note:"流量最大数量"`
	CurCount int    `json:"CurCount" note:"流量当前数量"`

	Flows MonitorNetworkIOThroughputCollection `json:"flows" note:"流量"`
}
