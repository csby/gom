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

	WSNginxAppAdded   = 1051 // 添加nginx应用
	WSNginxAppUpdated = 1052 // 更新nginx应用
	WSNginxAppDeleted = 1053 // 删除nginx应用

	WSNetworkThroughput = 2011 // 网络吞吐量
	WSCpuUsage          = 2012 // CPU使用率
	WSMemUsage          = 2013 // 内存使用率

	WSReviseProxyServiceStatus  = 3001 // 反向代理服务状态信息
	WSReviseProxyConnectionOpen = 3002 // 反向代理连接已打开
	WSReviseProxyConnectionShut = 3003 // 反向代理连接已关闭

	WSReviseProxyServerAdd = 3011 // 反向代理添加服务器
	WSReviseProxyServerDel = 3012 // 反向代理删除服务器
	WSReviseProxyServerMod = 3013 // 反向代理修改服务器

	WSReviseProxyTargetAdd = 3021 // 反向代理添加目标地址
	WSReviseProxyTargetDel = 3022 // 反向代理删除目标地址
	WSReviseProxyTargetMod = 3023 // 反向代理修改目标地址

	WSReviseProxyTargetStatusChanged = 3031 // 反向代理目标地址活动状态改变
)
