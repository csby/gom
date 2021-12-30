package model

type ServerArgument struct {
	Name string `json:"name" required:"true" not:"名称"`
}
