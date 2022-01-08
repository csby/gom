package model

type MonitorCpuUsage struct {
	CpuName  string `json:"cpuName" note:"CPU名称"`
	MaxCount int    `json:"maxCount" note:"最大数量"`
	CurCount int    `json:"CurCount" note:"当前数量"`

	Percents MonitorCpuPercentCollection `json:"percents" note:"使用率"`
}
