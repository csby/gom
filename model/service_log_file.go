package model

import "github.com/csby/gwsf/gtype"

type ServiceLogFile struct {
	Name     string         `json:"name" note:"文件名称"`
	Size     int64          `json:"size" note:"大小, 单位字节"`
	SizeText string         `json:"sizeText" note:"大小文本信息"`
	ModTime  gtype.DateTime `json:"modTime" note:"修改时间"`
	Path     string         `json:"path" note:"路径, base64"`
}
