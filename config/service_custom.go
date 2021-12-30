package config

type ServiceCustom struct {
	App string `json:"app" note:"程序根目录"`
	Log string `json:"log" note:"日志根目录"`
}
