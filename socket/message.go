package socket

const (
	WSSvcStatusChanged = 1011 // 服务状态改变

	WSCustomSvcAdded   = 1021 // 添加自定义服务
	WSCustomSvcUpdated = 1022 // 更新自定义服务
	WSCustomSvcDeleted = 1023 // 删除自定义服务

	WSTomcatAppAdded   = 1031 // 添加tomcat应用
	WSTomcatAppUpdated = 1032 // 更新tomcat应用
	WSTomcatAppDeleted = 1033 // 删除tomcat应用

	WSTomcatCfgAdded   = 1041 // 添加tomcat配置
	WSTomcatCfgUpdated = 1042 // 更新tomcat配置
	WSTomcatCfgDeleted = 1043 // 删除tomcat配置

	WSNetworkThroughput = 2011 // 网络吞吐量
)
