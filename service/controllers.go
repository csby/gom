package main

import (
	"github.com/csby/gom/controller"
	"github.com/csby/gwsf/gtype"
)

type Controllers struct {
	svc     *controller.Service
	monitor *controller.Monitor
	proxy   *controller.Proxy
}

func (s *Controllers) initController(wsc gtype.SocketChannelCollection) {
	s.svc = controller.NewService(log, cfg, wsc)
	s.monitor = controller.NewMonitor(log, cfg, wsc)
	s.proxy = controller.NewProxy(log, cfg, wsc)
}

func (s *Controllers) initRouter(router gtype.Router, path *gtype.Path, preHandle gtype.HttpHandle) {
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
	router.POST(path.Uri("/svc/tomcat/app/detail"), preHandle,
		s.svc.GetTomcatAppDetail, s.svc.GetTomcatDetailDoc)

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

	// 系统服务-nginx
	router.POST(path.Uri("/svc/nginx/svc/list"), preHandle,
		s.svc.GetNginxes, s.svc.GetNginxesDoc)
	router.POST(path.Uri("/svc/nginx/svc/start"), preHandle,
		s.svc.StartNginx, s.svc.StartNginxDoc)
	router.POST(path.Uri("/svc/nginx/svc/stop"), preHandle,
		s.svc.StopNginx, s.svc.StopNginxDoc)
	router.POST(path.Uri("/svc/nginx/svc/restart"), preHandle,
		s.svc.RestartNginx, s.svc.RestartNginxDoc)
	router.POST(path.Uri("/svc/nginx/app/mod"), preHandle,
		s.svc.ModNginxApp, s.svc.ModNginxAppDoc)
	router.POST(path.Uri("/svc/nginx/app/detail"), preHandle,
		s.svc.GetNginxAppDetail, s.svc.GetNginxAppDetailDoc)

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
	router.POST(path.Uri("/svc/custom/app/detail"), preHandle,
		s.svc.GetCustomDetail, s.svc.GetCustomDetailDoc)

	router.POST(path.Uri("/svc/custom/log/file/list"), preHandle,
		s.svc.GetCustomLogFiles, s.svc.GetCustomLogFilesDoc)
	router.GET(path.Uri("/svc/custom/log/file/download/:path"), preHandle,
		s.svc.DownloadCustomLogFile, s.svc.DownloadCustomLogFileDoc)
	router.GET(path.Uri("/svc/custom/log/file/content/:path"), preHandle,
		s.svc.ViewCustomLogFile, s.svc.ViewCustomLogFileDoc)

	// 系统服务-其他
	router.POST(path.Uri("/svc/other/svc/list"), preHandle,
		s.svc.GetOthers, s.svc.GetOthersDoc)
	router.POST(path.Uri("/svc/other/svc/start"), preHandle,
		s.svc.StartOther, s.svc.StartOtherDoc)
	router.POST(path.Uri("/svc/other/svc/stop"), preHandle,
		s.svc.StopOther, s.svc.StopOtherDoc)
	router.POST(path.Uri("/svc/other/svc/restart"), preHandle,
		s.svc.RestartOther, s.svc.RestartOtherDoc)

	// 系统服务-文件
	router.GET(path.Uri("/svc/file/content/:path"), preHandle,
		s.svc.ViewFile, s.svc.ViewFileDoc)
	router.GET(path.Uri("/svc/file/download/:path"), preHandle,
		s.svc.DownloadFile, s.svc.DownloadFileDoc)
	router.POST(path.Uri("/svc/file/mod"), preHandle,
		s.svc.ModFile, s.svc.ModFileDoc)
	router.POST(path.Uri("/svc/file/del"), preHandle,
		s.svc.DeleteFile, s.svc.DeleteFileDoc)

	// 系统资源-磁盘
	router.POST(path.Uri("/monitor/disk/usage/list"), preHandle,
		s.monitor.GetDiskPartitionUsages, s.monitor.GetDiskPartitionUsagesDoc)

	// 系统资源-网络
	router.POST(path.Uri("/monitor/network/interface/list"), preHandle,
		s.monitor.GetNetworkInterfaces, s.monitor.GetNetworkInterfacesDoc)
	router.POST(path.Uri("/monitor/network/throughput/list"), preHandle,
		s.monitor.GetNetworkThroughput, s.monitor.GetNetworkThroughputDoc)

	// 系统资源-CPU
	router.POST(path.Uri("/monitor/cpu/usage/list"), preHandle,
		s.monitor.GetCpuUsage, s.monitor.GetCpuUsageDoc)

	// 系统资源-内存
	router.POST(path.Uri("/monitor/mem/usage/list"), preHandle,
		s.monitor.GetMemoryUsage, s.monitor.GetMemoryUsageDoc)

	// 反向代理-服务
	router.POST(path.Uri("/proxy/service/setting/get"), preHandle,
		s.proxy.GetProxyServiceSetting, s.proxy.GetProxyServiceSettingDoc)
	router.POST(path.Uri("/proxy/service/setting/set"), preHandle,
		s.proxy.SetProxyServiceSetting, s.proxy.SetProxyServiceSettingDoc)
	router.POST(path.Uri("/proxy/service/status"), preHandle,
		s.proxy.GetProxyServiceStatus, s.proxy.GetProxyServiceStatusDoc)
	router.POST(path.Uri("/proxy/service/start"), preHandle,
		s.proxy.StartProxyService, s.proxy.StartProxyServiceDoc)
	router.POST(path.Uri("/proxy/service/stop"), preHandle,
		s.proxy.StopProxyService, s.proxy.StopProxyServiceDoc)
	router.POST(path.Uri("/proxy/service/restart"), preHandle,
		s.proxy.RestartProxyService, s.proxy.RestartProxyServiceDoc)

	// 反向代理-连接
	router.POST(path.Uri("/proxy/conn/list"), preHandle,
		s.proxy.GetProxyLinks, s.proxy.GetProxyLinksDoc)

	// 反向代理-端口
	router.POST(path.Uri("/proxy/server/list"), preHandle,
		s.proxy.GetProxyServers, s.proxy.GetProxyServersDoc)
	router.POST(path.Uri("/proxy/server/add"), preHandle,
		s.proxy.AddProxyServer, s.proxy.AddProxyServerDoc)
	router.POST(path.Uri("/proxy/server/del"), preHandle,
		s.proxy.DelProxyServer, s.proxy.DelProxyServerDoc)
	router.POST(path.Uri("/proxy/server/mod"), preHandle,
		s.proxy.ModifyProxyServer, s.proxy.ModifyProxyServerDoc)

	// 反向代理-目标
	router.POST(path.Uri("/proxy/target/list"), preHandle,
		s.proxy.GetProxyTargets, s.proxy.GetProxyTargetsDoc)
	router.POST(path.Uri("/proxy/target/add"), preHandle,
		s.proxy.AddProxyTarget, s.proxy.AddProxyTargetDoc)
	router.POST(path.Uri("/proxy/target/del"), preHandle,
		s.proxy.DelProxyTarget, s.proxy.DelProxyTargetDoc)
	router.POST(path.Uri("/proxy/target/mod"), preHandle,
		s.proxy.ModifyProxyTarget, s.proxy.ModifyProxyTargetDoc)
}
