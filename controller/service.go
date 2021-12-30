package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/csby/gom/config"
	"github.com/csby/gom/model"
	"github.com/csby/gwsf/gtype"
	"github.com/kardianos/service"
	"io/ioutil"
	"path/filepath"
	"runtime"
)

const (
	svcCatalogRoot   = "系统服务"
	svcCatalogCustom = "自定义"
)

func NewService(log gtype.Log, cfg *config.Config, wsc gtype.SocketChannelCollection) *Service {
	inst := &Service{}
	inst.SetLog(log)
	inst.cfg = cfg
	inst.wsc = wsc

	return inst
}

type Service struct {
	base
}

func (s *Service) getStatus(name string) (gtype.ServerStatus, error) {
	cfg := &service.Config{}
	if runtime.GOOS == "linux" {
		cfg.Name = fmt.Sprintf("%s.service", name)
	} else {
		cfg.Name = name
	}

	svc, err := service.New(nil, cfg)
	if err != nil {
		return gtype.ServerStatusUnknown, err
	}

	status, err := svc.Status()
	if err != nil {
		if err == service.ErrNotInstalled {
			return gtype.ServerStatusUnknown, nil
		}
		return gtype.ServerStatusUnknown, err
	}

	return gtype.ServerStatus(status), nil
}

func (s *Service) uninstall(name string) error {
	cfg := &service.Config{Name: name}
	svc, err := service.New(nil, cfg)
	if err != nil {
		return err
	}

	return svc.Uninstall()
}

func (s *Service) start(name string) error {
	cfg := &service.Config{Name: name}
	svc, err := service.New(nil, cfg)
	if err != nil {
		return err
	}

	return svc.Start()
}

func (s *Service) stop(name string) error {
	cfg := &service.Config{Name: name}
	svc, err := service.New(nil, cfg)
	if err != nil {
		return err
	}

	return svc.Stop()
}

func (s *Service) restart(name string) error {
	cfg := &service.Config{Name: name}
	svc, err := service.New(nil, cfg)
	if err != nil {
		return err
	}

	return svc.Restart()
}

func (s *Service) getFiles(folder string) []*model.ServiceLogFile {
	files := make([]*model.ServiceLogFile, 0)

	fs, fe := ioutil.ReadDir(folder)
	if fe != nil {
		return files
	}

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		file := &model.ServiceLogFile{}
		files = append(files, file)
		file.Name = f.Name()
		file.Size = f.Size()
		file.ModTime = gtype.DateTime(f.ModTime())
		file.SizeText = s.sizeToText(float64(file.Size))
		file.Path = base64.URLEncoding.EncodeToString([]byte(filepath.Join(folder, file.Name)))
	}

	return files
}