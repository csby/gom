package main

import (
	"github.com/csby/gom/controller"
	"github.com/csby/gwsf/gtype"
)

type Controllers struct {
	svc *controller.Service
}

func (s *Controllers) initController(wsc gtype.SocketChannelCollection) {
	s.svc = controller.NewService(log, cfg, wsc)
}

func (s *Controllers) initRouter(router gtype.Router, path *gtype.Path, preHandle gtype.HttpHandle) {
	// 系统服务-自定义
	router.POST(path.Uri("/svc/custom/list"), preHandle,
		s.svc.GetCustoms, s.svc.GetCustomsDoc)
	router.POST(path.Uri("/svc/custom/add"), preHandle,
		s.svc.AddCustom, s.svc.AddCustomDoc)
	router.POST(path.Uri("/svc/custom/mod"), preHandle,
		s.svc.ModCustom, s.svc.ModCustomDoc)
	router.POST(path.Uri("/svc/custom/del"), preHandle,
		s.svc.DelCustom, s.svc.DelCustomDoc)
	router.POST(path.Uri("/svc/custom/install"), preHandle,
		s.svc.InstallCustom, s.svc.InstallCustomDoc)
	router.POST(path.Uri("/svc/custom/uninstall"), preHandle,
		s.svc.UninstallCustom, s.svc.UninstallCustomDoc)
	router.POST(path.Uri("/svc/custom/start"), preHandle,
		s.svc.StartCustom, s.svc.StartCustomDoc)
	router.POST(path.Uri("/svc/custom/stop"), preHandle,
		s.svc.StopCustom, s.svc.StopCustomDoc)
	router.POST(path.Uri("/svc/custom/restart"), preHandle,
		s.svc.RestartCustom, s.svc.RestartCustomDoc)
	router.POST(path.Uri("/svc/custom/log/file/list"), preHandle,
		s.svc.GetCustomLogFiles, s.svc.GetCustomLogFilesDoc)
	router.GET(path.Uri("/svc/custom/log/file/download/:path"), preHandle,
		s.svc.DownloadCustomLogFile, s.svc.DownloadCustomLogFileDoc)
	router.GET(path.Uri("/svc/custom/log/file/content/:path"), preHandle,
		s.svc.ViewCustomLogFile, s.svc.ViewCustomLogFileDoc)
}
