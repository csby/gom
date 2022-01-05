package main

import (
	"github.com/csby/gom/controller"
	"github.com/csby/gwsf/gtype"
)

type Controllers struct {
	svc     *controller.Service
	monitor *controller.Monitor
}

func (s *Controllers) initController(wsc gtype.SocketChannelCollection) {
	s.svc = controller.NewService(log, cfg, wsc)
	s.monitor = controller.NewMonitor(log, cfg, wsc)
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
	router.GET(path.Uri("/svc/custom/download/:name"), preHandle,
		s.svc.DownloadCustom, s.svc.DownloadCustomDoc)
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

	// 系统服务-tomcat
	router.POST(path.Uri("/svc/tomcat/svc/list"), preHandle,
		s.svc.GetTomcats, s.svc.GetTomcatsDoc)
	router.POST(path.Uri("/svc/tomcat/svc/start"), preHandle,
		s.svc.StartTomcat, s.svc.StartTomcatDoc)
	router.POST(path.Uri("/svc/tomcat/svc/stop"), preHandle,
		s.svc.StopTomcat, s.svc.StopTomcatDoc)
	router.POST(path.Uri("/svc/tomcat/svc/restart"), preHandle,
		s.svc.RestartTomcat, s.svc.RestartTomcatDoc)

	router.POST(path.Uri("/svc/tomcat/app/list"), preHandle,
		s.svc.GetTomcatApps, s.svc.GetTomcatAppsDoc)
	router.GET(path.Uri("/svc/tomcat/app/download/:name/:app"), preHandle,
		s.svc.DownloadTomcatApp, s.svc.DownloadTomcatAppDoc)
	router.POST(path.Uri("/svc/tomcat/app/mod"), preHandle,
		s.svc.ModTomcatApp, s.svc.ModTomcatAppDoc)
	router.POST(path.Uri("/svc/tomcat/app/del"), preHandle,
		s.svc.DelTomcatApp, s.svc.DelTomcatAppDoc)

	router.POST(path.Uri("/svc/tomcat/cfg/tree"), preHandle,
		s.svc.GetTomcatConfigs, s.svc.GetTomcatConfigsDoc)
	router.GET(path.Uri("/svc/tomcat/cfg/file/content/:name/:path"), preHandle,
		s.svc.ViewTomcatConfigFile, s.svc.ViewTomcatConfigFileDoc)
	router.GET(path.Uri("/svc/tomcat/cfg/file/download/:name/:path"), preHandle,
		s.svc.DownloadTomcatConfigFile, s.svc.DownloadTomcatConfigFileDoc)
	router.POST(path.Uri("/svc/tomcat/cfg/folder/add"), preHandle,
		s.svc.CreateTomcatConfigFolder, s.svc.CreateTomcatConfigFolderDoc)
	router.POST(path.Uri("/svc/tomcat/cfg/mod"), preHandle,
		s.svc.ModTomcatConfig, s.svc.ModTomcatConfigDoc)
	router.POST(path.Uri("/svc/tomcat/cfg/del"), preHandle,
		s.svc.DeleteTomcatConfig, s.svc.DeleteTomcatConfigDoc)

	router.POST(path.Uri("/svc/tomcat/log/tree"), preHandle,
		s.svc.GetTomcatLogs, s.svc.GetTomcatLogsDoc)
	router.GET(path.Uri("/svc/tomcat/log/file/content/:name/:path"), preHandle,
		s.svc.ViewTomcatLogFile, s.svc.ViewTomcatLogFileDoc)
	router.GET(path.Uri("/svc/tomcat/log/file/download/:name/:path"), preHandle,
		s.svc.DownloadTomcatLogFile, s.svc.DownloadTomcatLogFileDoc)
	router.POST(path.Uri("/svc/tomcat/log/del"), preHandle,
		s.svc.DeleteTomcatLog, s.svc.DeleteTomcatLogDoc)

	// 系统服务-其他
	router.POST(path.Uri("/svc/other/svc/list"), preHandle,
		s.svc.GetOthers, s.svc.GetOthersDoc)
	router.POST(path.Uri("/svc/other/svc/start"), preHandle,
		s.svc.StartOther, s.svc.StartOtherDoc)
	router.POST(path.Uri("/svc/other/svc/stop"), preHandle,
		s.svc.StopOther, s.svc.StopOtherDoc)
	router.POST(path.Uri("/svc/other/svc/restart"), preHandle,
		s.svc.RestartOther, s.svc.RestartOtherDoc)

	// 系统资源-磁盘
	router.POST(path.Uri("/monitor/disk/usage/list"), preHandle,
		s.monitor.GetDiskPartitionUsages, s.monitor.GetDiskPartitionUsagesDoc)

	// 系统资源-网络
	router.POST(path.Uri("/monitor/network/interface/list"), preHandle,
		s.monitor.GetNetworkInterfaces, s.monitor.GetNetworkInterfacesDoc)
	router.POST(path.Uri("/monitor/network/throughput/list"), preHandle,
		s.monitor.GetNetworkThroughput, s.monitor.GetNetworkThroughputDoc)

}
