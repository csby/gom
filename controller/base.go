package controller

import (
	"fmt"
	"github.com/csby/gom/config"
	"github.com/csby/gwsf/gtype"
)

type base struct {
	gtype.Base

	cfg *config.Config
	wsc gtype.SocketChannelCollection
}

func (s *base) createCatalog(doc gtype.Doc, names ...string) gtype.Catalog {
	root := doc.AddCatalog("管理平台接口")

	count := len(names)
	if count < 1 {
		return root
	}

	child := root
	for i := 0; i < count; i++ {
		name := names[i]
		child = child.AddChild(name)
	}

	return child
}

func (s *base) writeOptMessage(id int, data interface{}) bool {
	if s.wsc == nil {
		return false
	}
	msg := &gtype.SocketMessage{
		ID:   id,
		Data: data,
	}

	s.wsc.Write(msg, nil)

	return true
}

func (s *base) sizeToText(v float64) string {
	kb := float64(1024)
	mb := 1024 * kb
	gb := 1024 * mb

	if v >= gb {
		return fmt.Sprintf("%.1fGB", v/gb)
	} else if v >= mb {
		return fmt.Sprintf("%.1fMB", v/mb)
	} else if v >= kb {
		return fmt.Sprintf("%.1fKB", v/kb)
	} else {
		return fmt.Sprintf("%.0fB", v)
	}
}
