package controller

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/csby/gom/model"
	"github.com/csby/gom/socket"
	"github.com/csby/gwsf/gfile"
	"github.com/csby/gwsf/gtype"
	"github.com/kardianos/service"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func (s *Service) GetCustoms(ctx gtype.Context, ps gtype.Params) {
	results := make([]*model.ServiceCustomInfo, 0)
	rootFolder := s.cfg.Sys.Svc.Custom.App
	if len(rootFolder) > 0 {
		fs, fe := ioutil.ReadDir(rootFolder)
		if fe == nil {
			for _, f := range fs {
				if f.IsDir() {
					infoPath := filepath.Join(rootFolder, f.Name(), "info.json")
					info := &model.ServiceCustomInfo{}
					err := info.LoadFromFile(infoPath)
					if err == nil {
						info.SystemName = info.ServiceName()
						info.Folder = filepath.Dir(infoPath)
						info.DeployTime = gtype.DateTime(f.ModTime())
						if len(info.DisplayName) < 1 {
							info.DisplayName = info.Name
						}
						info.Status, _ = s.getStatus(info.SystemName)
						results = append(results, info)
					}
				}
			}
		}
	}

	ctx.Success(results)
}

func (s *Service) GetCustomsDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "获取服务列表")
	function.SetOutputDataExample([]*model.ServiceCustomInfo{
		{
			Name:        "example",
			SystemName:  "svc-example",
			DisplayName: "自定义服务示例",
			DeployTime:  gtype.DateTime(time.Now()),
		},
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) AddCustom(ctx gtype.Context, ps gtype.Params) {
	rootFolder := s.cfg.Sys.Svc.Custom.App
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "服务物理根路径为空")
		return
	}

	uploadFile, _, err := ctx.Request().FormFile("file")
	if err != nil {
		ctx.Error(gtype.ErrInput, "上传文件无效: ", err)
		return
	}
	defer uploadFile.Close()

	buf := &bytes.Buffer{}
	fileSize, err := buf.ReadFrom(uploadFile)
	if err != nil {
		ctx.Error(gtype.ErrInput, "读取文件上传文件失败: ", err)
		return
	}
	if fileSize < 1 {
		ctx.Error(gtype.ErrInput, "上传的文件无效: 文件大小为0")
		return
	}

	tempFolder := filepath.Join(filepath.Dir(rootFolder), ctx.NewGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		ctx.Error(gtype.ErrInternal, fmt.Sprintf("创建临时文件夹'%s'失败: ", tempFolder), err)
		return
	}
	defer os.RemoveAll(tempFolder)

	fileData := buf.Bytes()
	zipFile := &gfile.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		tarFile := &gfile.Tar{}
		err = tarFile.DecompressMemory(fileData, tempFolder)
		if err != nil {
			ctx.Error(gtype.ErrInternal, "解压文件失败: ", err)
			return
		}
	}

	infoFilePath, err := s.getCustomInfoPath(tempFolder)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	info := &model.ServiceCustomInfo{}
	err = info.LoadFromFile(infoFilePath)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "读取信息文件(info.json)失败: ", err)
		return
	}
	if len(info.Name) < 1 {
		ctx.Error(gtype.ErrInput, "信息文件(info.json)中的名称(name)为空")
		return
	}

	srvFolder := filepath.Join(rootFolder, info.Name)
	_, err = os.Stat(srvFolder)
	if !os.IsNotExist(err) {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("服务(%s)已存在", info.Name))
		return
	}

	err = gfile.Copy(filepath.Dir(infoFilePath), srvFolder)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "拷贝文件夹失败: ", err)
		return
	}

	info.SystemName = info.ServiceName()
	info.DeployTime = gtype.DateTime(time.Now())
	info.Folder = srvFolder
	go s.writeOptMessage(socket.WSCustomSvcAdded, info)

	err = s.installCustom(info)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "上传成功, 但安装失败: ", err)
		return
	}

	err = s.start(info.ServiceName())
	if err != nil {
		ctx.Error(gtype.ErrInternal, "上传并安装成功, 但启动失败: ", err)
		return
	}

	svcStatus := &model.ServiceStatus{Name: info.ServiceName()}
	svcStatus.Status, err = s.getStatus(info.ServiceName())
	if err == nil {
		go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
	}

	ctx.Success(info)
}

func (s *Service) AddCustomDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "添加服务")
	function.SetNote("上传服务程序打包文件(.zip或.tar.gz)，并安装成系统服务，成功时返回服务信息")
	function.SetRemark("打包文件中的根目录必须包含服务信息文件(info.json)，且服务名称(name)和可执行程序(exec)不能能为空")
	function.AddInputHeader(true, "content-type", "内容类型", gtype.ContentTypeFormData)
	function.AddInputForm(true, "file", "服务程序打包文件(.zip或.tar.gz)", gtype.FormValueKindFile, nil)
	function.SetOutputDataExample(&model.ServiceCustomInfo{
		Name:        "example",
		SystemName:  "svc-example",
		DisplayName: "自定义服务示例",
		DeployTime:  gtype.DateTime(time.Now()),
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) ModCustom(ctx gtype.Context, ps gtype.Params) {
	name := strings.TrimSpace(strings.ToLower(ctx.Request().FormValue("name")))
	if len(name) < 1 {
		ctx.Error(gtype.ErrInput, "服务名称(name)为空")
		return
	}

	rootFolder := s.cfg.Sys.Svc.Custom.App
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "服务物理根路径为空")
		return
	}

	uploadFile, _, err := ctx.Request().FormFile("file")
	if err != nil {
		ctx.Error(gtype.ErrInput, "上传文件无效: ", err)
		return
	}
	defer uploadFile.Close()

	buf := &bytes.Buffer{}
	fileSize, err := buf.ReadFrom(uploadFile)
	if err != nil {
		ctx.Error(gtype.ErrInput, "读取文件上传文件失败: ", err)
		return
	}
	if fileSize < 1 {
		ctx.Error(gtype.ErrInput, "上传的文件无效: 文件大小为0")
		return
	}

	tempFolder := filepath.Join(filepath.Dir(rootFolder), ctx.NewGuid())
	err = os.MkdirAll(tempFolder, 0777)
	if err != nil {
		ctx.Error(gtype.ErrInternal, fmt.Sprintf("创建临时文件夹'%s'失败: ", tempFolder), err)
		return
	}
	defer os.RemoveAll(tempFolder)

	fileData := buf.Bytes()
	zipFile := &gfile.Zip{}
	err = zipFile.DecompressMemory(fileData, tempFolder)
	if err != nil {
		tarFile := &gfile.Tar{}
		err = tarFile.DecompressMemory(fileData, tempFolder)
		if err != nil {
			ctx.Error(gtype.ErrInternal, "解压文件失败: ", err)
			return
		}
	}

	infoFilePath, err := s.getCustomInfoPath(tempFolder)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	info := &model.ServiceCustomInfo{}
	err = info.LoadFromFile(infoFilePath)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "读取信息文件(info.json)失败: ", err)
		return
	}
	if len(info.Name) < 1 {
		ctx.Error(gtype.ErrInput, "信息文件(info.json)中的名称(name)为空")
		return
	}
	if info.Name != name {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("信息文件(info.json)中的名称(%s)和目标服务名称(%s)不一致", info.Name, name))
		return
	}

	svcStatus := &model.ServiceStatus{Name: info.ServiceName()}
	svcStatus.Status, err = s.getStatus(info.ServiceName())
	if err != nil {
		if svcStatus.Status == gtype.ServerStatusRunning {
			err = s.stop(info.ServiceName())
			if err != nil {
				ctx.Error(gtype.ErrInternal, "停止服务失败: ", err)
				return
			}

			svcStatus.Status, err = s.getStatus(info.ServiceName())
			if err == nil {
				go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
			}
		}
	}

	srvFolder := filepath.Join(rootFolder, info.Name)
	err = os.RemoveAll(srvFolder)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "删除原服务文件夹失败: ", err)
		return
	}

	err = gfile.Copy(filepath.Dir(infoFilePath), srvFolder)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "拷贝文件夹失败: ", err)
		return
	}

	info.SystemName = info.ServiceName()
	info.DeployTime = gtype.DateTime(time.Now())
	info.Folder = srvFolder
	go s.writeOptMessage(socket.WSCustomSvcUpdated, info)

	if svcStatus.Status == gtype.ServerStatusStopped {
		err = s.start(info.ServiceName())
		if err != nil {
			ctx.Error(gtype.ErrInternal, "更新成功, 但启动失败: ", err)
			return
		}
	}
	svcStatus.Status, err = s.getStatus(info.ServiceName())
	if err == nil {
		go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
	}

	ctx.Success(info)
}

func (s *Service) ModCustomDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "更新服务")
	function.SetNote("上传服务程序打包文件(.zip或.tar.gz)，并安装成系统服务，成功时返回服务信息")
	function.SetRemark("打包文件中的根目录必须包含服务信息文件(info.json)，且服务名称(name)和可执行程序(exec)不能能为空")
	function.AddInputHeader(true, "content-type", "内容类型", gtype.ContentTypeFormData)
	function.AddInputForm(true, "name", "目标服务名称(需和info.json中的name一致)", gtype.FormValueKindText, "")
	function.AddInputForm(true, "file", "服务程序打包文件(.zip或.tar.gz)", gtype.FormValueKindFile, nil)
	function.SetOutputDataExample(&model.ServiceCustomInfo{
		Name:        "example",
		SystemName:  "svc-example",
		DisplayName: "自定义服务示例",
		DeployTime:  gtype.DateTime(time.Now()),
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) DelCustom(ctx gtype.Context, ps gtype.Params) {
	argument := &model.ServerArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInternal, "服务名称(name)为空")
		return
	}

	rootFolder := s.cfg.Sys.Svc.Custom.App
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "服务物理根路径为空")
		return
	}

	infoFolder := filepath.Join(rootFolder, argument.Name)
	infoPath := filepath.Join(infoFolder, "info.json")
	info := &model.ServiceCustomInfo{}
	err = info.LoadFromFile(infoPath)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "读取信息文件(info.json)失败: ", err)
		return
	}
	if len(info.Name) < 1 {
		ctx.Error(gtype.ErrInput, "信息文件(info.json)中的名称(name)为空")
		return
	}
	info.Folder = infoFolder

	svcStatus := &model.ServiceStatus{Name: info.ServiceName()}
	svcStatus.Status, err = s.getStatus(info.ServiceName())
	if err != nil {
		ctx.Error(gtype.ErrInternal, "获取服务状态失败: ", err)
		return
	}
	if svcStatus.Status != gtype.ServerStatusUnknown {
		ctx.Error(gtype.ErrInternal, "服务未卸载")
		return
	}

	err = os.RemoveAll(infoFolder)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "删除原服务文件夹失败: ", err)
		return
	}

	go s.writeOptMessage(socket.WSCustomSvcDeleted, argument)

	ctx.Success(argument)
}

func (s *Service) DelCustomDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "删除服务")
	function.SetRemark("已卸载服务才能删除")
	function.SetInputJsonExample(&model.ServerArgument{
		Name: "example",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) InstallCustom(ctx gtype.Context, ps gtype.Params) {
	argument := &model.ServerArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInternal, "服务名称(name)为空")
		return
	}

	rootFolder := s.cfg.Sys.Svc.Custom.App
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "服务物理根路径为空")
		return
	}

	infoFolder := filepath.Join(rootFolder, argument.Name)
	info, ge := s.getCustomInfo(infoFolder)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	err = s.installCustom(info)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	svcStatus := &model.ServiceStatus{Name: info.ServiceName()}
	svcStatus.Status, err = s.getStatus(info.ServiceName())
	if err == nil {
		go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
	}

	ctx.Success(info)
}

func (s *Service) InstallCustomDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "安装服务")
	function.SetInputJsonExample(&model.ServerArgument{
		Name: "example",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) UninstallCustom(ctx gtype.Context, ps gtype.Params) {
	argument := &model.ServerArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInternal, "服务名称(name)为空")
		return
	}

	rootFolder := s.cfg.Sys.Svc.Custom.App
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "服务物理根路径为空")
		return
	}

	infoFolder := filepath.Join(rootFolder, argument.Name)
	info, ge := s.getCustomInfo(infoFolder)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	err = s.uninstall(info.ServiceName())
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}
	svcStatus := &model.ServiceStatus{Name: info.ServiceName()}
	svcStatus.Status, err = s.getStatus(info.ServiceName())
	go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)

	ctx.Success(info)
}

func (s *Service) UninstallCustomDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "卸载服务")
	function.SetInputJsonExample(&model.ServerArgument{
		Name: "example",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) StartCustom(ctx gtype.Context, ps gtype.Params) {
	argument := &model.ServerArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInternal, "服务名称(name)为空")
		return
	}

	rootFolder := s.cfg.Sys.Svc.Custom.App
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "服务物理根路径为空")
		return
	}

	infoFolder := filepath.Join(rootFolder, argument.Name)
	info, ge := s.getCustomInfo(infoFolder)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	err = s.start(info.ServiceName())
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	svcStatus := &model.ServiceStatus{Name: info.ServiceName()}
	svcStatus.Status, err = s.getStatus(info.ServiceName())
	if err == nil {
		go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
	}

	ctx.Success(info)
}

func (s *Service) StartCustomDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "启动服务")
	function.SetInputJsonExample(&model.ServerArgument{
		Name: "example",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) StopCustom(ctx gtype.Context, ps gtype.Params) {
	argument := &model.ServerArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInternal, "服务名称(name)为空")
		return
	}

	rootFolder := s.cfg.Sys.Svc.Custom.App
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "服务物理根路径为空")
		return
	}

	infoFolder := filepath.Join(rootFolder, argument.Name)
	info, ge := s.getCustomInfo(infoFolder)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	err = s.stop(info.ServiceName())
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}
	svcStatus := &model.ServiceStatus{Name: info.ServiceName()}
	svcStatus.Status, err = s.getStatus(info.ServiceName())
	if err == nil {
		go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
	}

	ctx.Success(info)
}

func (s *Service) StopCustomDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "停止服务")
	function.SetInputJsonExample(&model.ServerArgument{
		Name: "example",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) RestartCustom(ctx gtype.Context, ps gtype.Params) {
	argument := &model.ServerArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInternal, "服务名称(name)为空")
		return
	}

	rootFolder := s.cfg.Sys.Svc.Custom.App
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "服务物理根路径为空")
		return
	}

	infoFolder := filepath.Join(rootFolder, argument.Name)
	info, ge := s.getCustomInfo(infoFolder)
	if ge != nil {
		ctx.Error(ge)
		return
	}

	err = s.restart(info.ServiceName())
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}
	svcStatus := &model.ServiceStatus{Name: info.ServiceName()}
	svcStatus.Status, err = s.getStatus(info.ServiceName())
	if err == nil {
		go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
	}

	ctx.Success(info)
}

func (s *Service) RestartCustomDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "重启服务")
	function.SetInputJsonExample(&model.ServerArgument{
		Name: "example",
	})
	function.SetOutputDataExample(nil)
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) GetCustomLogFiles(ctx gtype.Context, ps gtype.Params) {
	argument := &model.ServerArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.Name) < 1 {
		ctx.Error(gtype.ErrInternal, "服务名称(name)为空")
		return
	}

	rootFolder := s.cfg.Sys.Svc.Custom.Log
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "服务物理根路径为空")
		return
	}
	logFolder := filepath.Join(rootFolder, argument.Name)

	results := s.getFiles(logFolder)

	ctx.Success(results)
}

func (s *Service) GetCustomLogFilesDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "获取服务日志文件列表")
	function.SetInputJsonExample(&model.ServerArgument{
		Name: "example",
	})
	function.SetOutputDataExample([]*model.ServiceLogFile{
		{
			Name:     "2021-12-29.log",
			Size:     12783,
			SizeText: s.sizeToText(float64(12783)),
			ModTime:  gtype.DateTime(time.Now()),
		},
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) DownloadCustomLogFile(ctx gtype.Context, ps gtype.Params) {
	path := ps.ByName("path")
	if len(path) < 1 {
		ctx.Error(gtype.ErrInput, "路径为空")
		return
	}

	pathData, err := base64.URLEncoding.DecodeString(path)
	if err != nil {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("路径(%s)无效", path))
		return
	}

	logPath := string(pathData)
	fi, fe := os.Stat(logPath)
	if os.IsNotExist(fe) {
		ctx.Error(gtype.ErrInternal, fe)
		return
	}
	logFile, le := os.OpenFile(logPath, os.O_RDONLY, 0666)
	if le != nil {
		ctx.Error(gtype.ErrInternal, le)
		return
	}
	defer logFile.Close()

	contentLength := fi.Size()
	fileName := fmt.Sprintf("%s", filepath.Base(logPath))
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprint("attachment; filename=", fileName))
	ctx.Response().Header().Set("Content-Length", fmt.Sprint(contentLength))

	io.Copy(ctx.Response(), logFile)
}

func (s *Service) DownloadCustomLogFileDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "下载服务日志文件")
	function.SetNote("服务日志文件(.log)")
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) ViewCustomLogFile(ctx gtype.Context, ps gtype.Params) {
	path := ps.ByName("path")
	if len(path) < 1 {
		ctx.Error(gtype.ErrInput, "路径为空")
		return
	}

	pathData, err := base64.URLEncoding.DecodeString(path)
	if err != nil {
		ctx.Error(gtype.ErrInput, fmt.Sprintf("路径(%s)无效", path))
		return
	}

	logPath := string(pathData)
	fi, fe := os.Stat(logPath)
	if os.IsNotExist(fe) {
		ctx.Error(gtype.ErrInternal, fe)
		return
	}
	logFile, le := os.OpenFile(logPath, os.O_RDONLY, 0666)
	if le != nil {
		ctx.Error(gtype.ErrInternal, le)
		return
	}
	defer logFile.Close()

	contentLength := fi.Size()
	ctx.Response().Header().Set("Content-Type", gtype.ContentTypeText)
	ctx.Response().Header().Set("Content-Length", fmt.Sprint(contentLength))

	io.Copy(ctx.Response(), logFile)
}

func (s *Service) ViewCustomLogFileDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogCustom)
	function := catalog.AddFunction(method, uri, "查看服务日志文件")
	function.SetNote("返回日志文件文本内容")
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) installCustom(info *model.ServiceCustomInfo) error {
	if info == nil {
		return fmt.Errorf("info is nil")
	}

	executable := filepath.Join(filepath.Dir(s.cfg.Module.Path), "gshell")
	if runtime.GOOS == "windows" {
		executable += ".exe"
	}
	logFolder := ""
	if len(s.cfg.Sys.Svc.Custom.Log) > 0 {
		logFolder = filepath.Join(s.cfg.Sys.Svc.Custom.Log, info.Name)
	}
	cfg := &service.Config{
		Name:        info.ServiceName(),
		Description: info.Description,
		Arguments:   []string{info.Folder, logFolder},
		Executable:  executable,
	}
	cfg.DisplayName = cfg.Name

	svc, err := service.New(nil, cfg)
	if err != nil {
		return err
	}

	return svc.Install()
}

func (s *Service) getCustomInfo(folder string) (*model.ServiceCustomInfo, gtype.Error) {
	path := filepath.Join(folder, "info.json")
	info := &model.ServiceCustomInfo{}
	err := info.LoadFromFile(path)
	if err != nil {
		return nil, gtype.ErrInternal.SetDetail("读取信息文件(info.json)失败: ", err)
	}
	if len(info.Name) < 1 {
		return nil, gtype.ErrInput.SetDetail("信息文件(info.json)中的名称(name)为空")
	}
	info.Folder = folder

	return info, nil
}

func (s *Service) getCustomInfoPath(folderPath string) (string, error) {
	filePath := filepath.Join(folderPath, "info.json")
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fs, e := ioutil.ReadDir(folderPath)
		if e != nil {
			return "", e
		}
		for _, f := range fs {
			if f.IsDir() {
				path, pe := s.getCustomInfoPath(filepath.Join(folderPath, f.Name()))
				if pe == nil {
					return path, nil
				}
			}
		}
		return "", fmt.Errorf("未包含服务信息文件(info.json)")
	} else {
		return filePath, nil
	}
}