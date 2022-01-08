package controller

import (
	"github.com/csby/gmonitor"
	"github.com/csby/gom/model"
	"github.com/csby/gom/socket"
	"github.com/csby/gwsf/gtype"
	"sort"
	"time"
)

func (s *Monitor) GetCpuUsage(ctx gtype.Context, ps gtype.Params) {
	results := make(model.MonitorCpuPercentCollection, 0)
	items := s.cpuUsage.Percents()
	c := len(items)
	for i := 0; i < c; i++ {
		item := items[i]
		if item == nil {
			continue
		}

		result := &model.MonitorCpuPercent{
			TimePoint: item.TimePoint,
			Usage:     item.UsagePercent,
		}
		result.Time = gtype.DateTime(time.Unix(item.TimePoint, 0))

		results = append(results, result)
	}

	sort.Sort(results)
	ctx.Success(&model.MonitorCpuUsage{
		CpuName:  s.getCpuName(),
		MaxCount: s.faces.MaxCounter,
		CurCount: len(results),
		Percents: results,
	})
}

func (s *Monitor) GetCpuUsageDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, monitorCatalogRoot, monitorCatalogCpu)
	function := catalog.AddFunction(method, uri, "获取CPU使用率")
	function.SetOutputDataExample(&model.MonitorCpuUsage{
		MaxCount: 60,
		CurCount: 1,
		Percents: []*model.MonitorCpuPercent{
			{
				Time:  gtype.DateTime(time.Now()),
				Usage: 35,
			},
		},
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Monitor) doStatCpuUsage(interval time.Duration) {
	for {
		time.Sleep(interval)

		all, busy, err := gmonitor.AllCpuTotalBusyTime()
		if err != nil {
			continue
		}

		t := time.Now().Unix()
		percent := s.cpuUsage.AddTime(t, all, busy)
		if percent != nil {
			s.sentCpuUsage(percent)
		}
	}
}

func (s *Monitor) sentCpuUsage(item *NetworkCpuPercent) {
	if item == nil {
		return
	}

	argument := &model.MonitorCpuPercent{
		Time:  gtype.DateTime(time.Unix(item.TimePoint, 0)),
		Usage: item.UsagePercent,
	}

	go s.writeOptMessage(socket.WSCpuUsage, argument)
}

func (s *Monitor) getCpuName() string {
	if len(s.cupName) > 0 {
		return s.cupName
	}

	name, err := gmonitor.CpuName()
	if err == nil {
		s.cupName = name
	}

	return s.cupName
}
