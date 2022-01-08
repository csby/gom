package model

import "github.com/csby/gwsf/gtype"

type MonitorCpuPercent struct {
	TimePoint int64 `json:"-" note:"时间点"`

	Time  gtype.DateTime `json:"time" note:"时间"`
	Usage float64        `json:"usage" note:"使用率"`
}

type MonitorCpuPercentCollection []*MonitorCpuPercent

func (s MonitorCpuPercentCollection) Len() int {
	return len(s)
}

func (s MonitorCpuPercentCollection) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s MonitorCpuPercentCollection) Less(i, j int) bool {
	return s[i].TimePoint < s[j].TimePoint
}
