package controller

import (
	"fmt"
	"github.com/csby/gmonitor"
	"github.com/csby/gom/model"
	"github.com/csby/gom/socket"
	"github.com/csby/gwsf/gtype"
	"sort"
	"time"
)

func (s *Monitor) GetNetworkInterfaces(ctx gtype.Context, ps gtype.Params) {
	results, err := gmonitor.Interfaces()
	if err != nil {
		ctx.Error(gtype.ErrInternal.SetDetail(err))
		return
	}

	ctx.Success(results)
}

func (s *Monitor) GetNetworkInterfacesDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, monitorCatalogRoot, monitorCatalogNetwork)
	function := catalog.AddFunction(method, uri, "获取网卡信息")

	function.SetOutputDataExample([]*gmonitor.Interface{
		{
			Name:       "eth0",
			MTU:        1500,
			MacAddress: "00:15:5d:16:b9:00",
			IPAddress:  []string{"192.168.1.1/24", "172.16.1.1/16"},
			Flags:      []string{"up", "broadcast", "multicast"},
		},
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Monitor) GetNetworkThroughput(ctx gtype.Context, ps gtype.Params) {
	argument := &model.MonitorNetworkIOArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput.SetDetail(err))
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInput.SetDetail("网卡名称(name)为空"))
		return
	}

	face := s.faces.GetInterface(argument.Name)
	if face == nil {
		ctx.Error(gtype.ErrInput.SetDetail(fmt.Sprintf("网卡名称(%s)不存在", argument.Name)))
		return
	}

	results := make(model.MonitorNetworkIOThroughputCollection, 0)
	items := face.Counters()
	c := len(items)
	for i := 0; i < c; i++ {
		item := items[i]
		if item == nil {
			continue
		}

		interval := uint64(item.TimeInterval)
		if interval < 1 {
			continue
		}
		result := &model.MonitorNetworkIOThroughput{
			TimePoint:      item.TimePoint,
			BytesSpeedSent: item.BytesSent / interval,
			BytesSpeedRecv: item.BytesRecv / interval,
		}
		result.Time = gtype.DateTime(time.Unix(item.TimePoint, 0))
		result.BytesSpeedSentText = s.toSpeedText(float64(item.BytesSent) / float64(interval))
		result.BytesSpeedRecvText = s.toSpeedText(float64(item.BytesRecv) / float64(interval))

		results = append(results, result)
	}

	sort.Sort(results)
	ctx.Success(&model.MonitorNetworkIO{
		Name:     argument.Name,
		MaxCount: s.faces.MaxCounter,
		CurCount: len(results),
		Flows:    results,
	})
}

func (s *Monitor) GetNetworkThroughputDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, monitorCatalogRoot, monitorCatalogNetwork)
	function := catalog.AddFunction(method, uri, "获取网卡吞吐量")
	function.SetInputJsonExample(&model.MonitorNetworkIOArgument{
		Name: "eth0",
	})
	function.SetOutputDataExample(&model.MonitorNetworkIO{
		Name:     "eth0",
		MaxCount: 60,
		CurCount: 1,
		Flows: []*model.MonitorNetworkIOThroughput{
			{
				Time:               gtype.DateTime(time.Now()),
				BytesSpeedSent:     3 * 1024,
				BytesSpeedRecv:     5 * 1024,
				BytesSpeedSentText: s.toSpeedText(float64(3 * 1024)),
				BytesSpeedRecvText: s.toSpeedText(float64(5 * 1024)),
			},
		},
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Monitor) doStatNetworkIO() {
	interval := time.Second

	for {
		time.Sleep(interval)

		items, err := gmonitor.StatNetworkIOs()
		if err != nil {
			continue
		}

		t := time.Now().Unix()
		c := len(items)
		for i := 0; i < c; i++ {
			item := items[i]
			if item == nil {
				continue
			}

			counter := s.faces.AddIOCounter(t, item)
			if counter != nil {
				s.sentNetworkThroughput(item.Name, counter)
			}
		}
	}
}

func (s *Monitor) sentNetworkThroughput(name string, item *NetworkIOCounter) {
	if item == nil {
		return
	}
	interval := uint64(item.TimeInterval)
	if interval < 1 {
		return
	}

	argument := &model.MonitorNetworkIOThroughputArgument{
		Name: name,
		Flow: model.MonitorNetworkIOThroughput{
			TimePoint:      item.TimePoint,
			BytesSpeedSent: item.BytesSent / interval,
			BytesSpeedRecv: item.BytesRecv / interval,
		},
	}
	argument.Flow.Time = gtype.DateTime(time.Unix(item.TimePoint, 0))
	argument.Flow.BytesSpeedSentText = s.toSpeedText(float64(item.BytesSent) / float64(interval))
	argument.Flow.BytesSpeedRecvText = s.toSpeedText(float64(item.BytesRecv) / float64(interval))

	go s.writeOptMessage(socket.WSNetworkThroughput, argument)
}
