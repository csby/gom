package config

import (
	"encoding/json"
	"fmt"
	"github.com/csby/gwsf/gcfg"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	sync.RWMutex
	gcfg.Config
}

func NewConfig() *Config {
	return &Config{
		Config: gcfg.Config{
			Log: gcfg.Log{
				Folder: "",
				Level:  "error|warning|info",
			},
			Http: gcfg.Http{
				Enabled:     true,
				Port:        9618,
				BehindProxy: false,
			},
			Https: gcfg.Https{
				Enabled:     false,
				Port:        9613,
				BehindProxy: false,
				Cert: gcfg.Crt{
					Ca: gcfg.CrtCa{
						File: "",
					},
					Server: gcfg.CrtPfx{
						File:     "",
						Password: "",
					},
				},
				RequestClientCert: false,
			},
			Site: gcfg.Site{
				Doc: gcfg.SiteDoc{
					Enabled: true,
				},
				Opt: gcfg.SiteOpt{
					Users: []*gcfg.SiteOptUser{
						{
							Account:  "admin",
							Password: "1",
							Name:     "管理员",
						},
					},
				},
			},
		},
	}
}

func (s *Config) FromFile() (*Config, error) {
	cfg := &Config{}
	err := cfg.LoadFromFile(s.Path)

	return cfg, err
}

func (s *Config) LoadFromFile(filePath string) error {
	s.Lock()
	defer s.Unlock()

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, s)
	if err == nil {
	}

	return err
}

func (s *Config) SaveToFile(filePath string) error {
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

func (s *Config) String() string {
	bytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(bytes[:])
}

func (s *Config) FormatString() string {
	bytes, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return ""
	}

	return string(bytes[:])
}
