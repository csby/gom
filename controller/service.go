package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/csby/gom/config"
	"github.com/csby/gom/model"
	"github.com/csby/gwsf/gtype"
	"github.com/kardianos/service"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	svcCatalogRoot   = "系统服务"
	svcCatalogTomcat = "tomcat"
	svcCatalogNginx  = "nginx"
	svcCatalogCustom = "自定义"
	svcCatalogOther  = "其他"
	svcCatalogFile   = "文件"
)

func NewService(log gtype.Log, cfg *config.Config, wsc gtype.SocketChannelCollection) *Service {
	inst := &Service{}
	inst.SetLog(log)
	inst.cfg = cfg
	inst.wsc = wsc

	inst.removeCustomLogs()

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
		} else if err.Error() == "service in failed state" {
			return gtype.ServerStatusStopped, nil
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

func (s *Service) getFiles(root, folder string) model.ServiceLogFileCollection {
	files := make(model.ServiceLogFileCollection, 0)

	fs, fe := ioutil.ReadDir(filepath.Join(root, folder))
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

func (s *Service) removeFiles(folder string, now time.Time, expiration time.Duration) {
	fs, fe := ioutil.ReadDir(folder)
	if fe != nil {
		return
	}

	for _, f := range fs {
		path := filepath.Join(folder, f.Name())
		if f.IsDir() {
			s.removeFiles(path, now, expiration)
		} else {
			if now.Sub(f.ModTime()) > expiration {
				os.Remove(path)
			}
		}
	}
}

func (s *Service) getJsonVersion(folderPath string) (string, error) {
	/*
		{
		  "version": "1.0.1.1"
		}
	*/
	filePath := filepath.Join(folderPath, "version.json")
	fi, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", err
	}
	if fi.IsDir() {
		return "", fmt.Errorf("%s is not file", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	app := &gtype.WebApp{}
	err = json.Unmarshal(data, app)
	if err != nil {
		return "", err
	}

	return app.Version, nil
}

func (s *Service) getTextVersion(folderPath string) (string, error) {
	/*
		1.0.1.1
	*/
	filePath := filepath.Join(folderPath, "version.txt")
	fi, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", err
	}
	if fi.IsDir() {
		return "", fmt.Errorf("%s is not file", filePath)
	}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (s *Service) getFileInfos(cfg *model.FileInfo, root string) {
	if cfg == nil {
		return
	}
	if cfg.Children == nil {
		cfg.Children = make(model.FileInfoCollection, 0)
	}

	if len(root) < 1 {
		return
	}
	fi, ie := os.Stat(root)
	if os.IsNotExist(ie) {
		return
	}
	if !fi.IsDir() {
		return
	}
	cfg.Name = fi.Name()
	cfg.Folder = true
	cfg.Size = 0
	cfg.Time = gtype.DateTime(fi.ModTime())
	cfg.Path = base64.URLEncoding.EncodeToString([]byte(root))

	s.getFileInfo(cfg, root, "")
}

func (s *Service) getFileInfo(cfg *model.FileInfo, root, folder string) {
	if cfg == nil {
		return
	}
	if cfg.Children == nil {
		cfg.Children = make(model.FileInfoCollection, 0)
	}

	if len(root) < 1 {
		return
	}

	fs, fe := ioutil.ReadDir(filepath.Join(root, folder))
	if fe != nil {
		return
	}

	for _, f := range fs {
		name := f.Name()
		path := filepath.Join(folder, name)
		item := &model.FileInfo{
			Name:     name,
			Path:     base64.URLEncoding.EncodeToString([]byte(filepath.Join(root, path))),
			Time:     gtype.DateTime(f.ModTime()),
			Children: make(model.FileInfoCollection, 0),
		}
		cfg.Children = append(cfg.Children, item)
		if f.IsDir() {
			item.Folder = true
			item.Size = 0
			s.getFileInfo(item, root, path)
		} else {
			item.Size = f.Size()
		}
		cfg.Size += item.Size
	}
}
