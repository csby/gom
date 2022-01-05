package controller

import (
	"fmt"
	"github.com/csby/gom/config"
	"github.com/csby/gwsf/gtype"
)

const (
	monitorCatalogRoot    = "系统资源"
	monitorCatalogDisk    = "磁盘"
	monitorCatalogNetwork = "网络"
)

func NewMonitor(log gtype.Log, cfg *config.Config, wsc gtype.SocketChannelCollection) *Monitor {
	inst := &Monitor{}
	inst.SetLog(log)
	inst.cfg = cfg
	inst.wsc = wsc

	inst.faces = &NetworkInterfaceCollection{
		MaxCounter: 30,
	}

	go inst.doStatNetworkIO()

	return inst
}

type Monitor struct {
	base

	faces *NetworkInterfaceCollection
}

func (s *Monitor) toSpeedText(v float64) string {
	kb := float64(1024)
	mb := 1024 * kb
	gb := 1024 * mb

	if v >= gb {
		return fmt.Sprintf("%.1fGbps", v/gb)
	} else if v >= mb {
		return fmt.Sprintf("%.1fMbps", v/mb)
	} else if v >= kb {
		return fmt.Sprintf("%.1fKbps", v/kb)
	} else {
		return fmt.Sprintf("%.0fBps", v)
	}
}
