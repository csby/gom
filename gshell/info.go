package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Info struct {
	sync.RWMutex

	Name string `json:"name" note:"服务名称"`
	Exec string `json:"exec" note:"可执行程序"`
	Args string `json:"args" note:"程序启动参数"`
}

func (s *Info) ServiceName() string {
	return fmt.Sprintf("svc-%s", s.Name)
}

func (s *Info) LoadFromFile(filePath string) error {
	s.Lock()
	defer s.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, s)
}

func (s *Info) SaveToFile(filePath string) error {
	s.Lock()
	defer s.Unlock()

	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return err
	}

	fileFolder := filepath.Dir(filePath)
	_, err = os.Stat(fileFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(fileFolder, 0777)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, string(bytes[:]))

	return err
}
