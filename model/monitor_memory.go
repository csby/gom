package model

type MonitorMemoryUsage struct {
	MaxCount int `json:"maxCount" note:"最大数量"`
	CurCount int `json:"CurCount" note:"当前数量"`

	Percents MonitorMemoryPercentCollection `json:"percents" note:"使用率"`
}
