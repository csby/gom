package controller

import (
	"bytes"
	"fmt"
	"github.com/csby/gom/model"
	"github.com/csby/gom/socket"
	"github.com/csby/gwsf/gfile"
	"github.com/csby/gwsf/gtype"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) StartNginx(ctx gtype.Context, ps gtype.Params) {
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
	info := s.cfg.Sys.Svc.GetNginxByServiceName(argument.Name)
	if info == nil {
		ctx.Error(gtype.ErrInternal, fmt.Sprintf("服务(%s)不存在", argument.Name))
		return
	}

	err = s.start(argument.Name)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	svcStatus := &model.ServiceStatus{Name: argument.Name}
	svcStatus.Status, err = s.getStatus(argument.Name)
	if err == nil {
		go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
	}

	ctx.Success(info)
}

func (s *Service) StartNginxDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogNginx)
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

func (s *Service) StopNginx(ctx gtype.Context, ps gtype.Params) {
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
	info := s.cfg.Sys.Svc.GetNginxByServiceName(argument.Name)
	if info == nil {
		ctx.Error(gtype.ErrInternal, fmt.Sprintf("服务(%s)不存在", argument.Name))
		return
	}

	err = s.stop(argument.Name)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	svcStatus := &model.ServiceStatus{Name: argument.Name}
	svcStatus.Status, err = s.getStatus(argument.Name)
	if err == nil {
		go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
	}

	ctx.Success(info)
}

func (s *Service) StopNginxDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogNginx)
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

func (s *Service) RestartNginx(ctx gtype.Context, ps gtype.Params) {
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
	info := s.cfg.Sys.Svc.GetNginxByServiceName(argument.Name)
	if info == nil {
		ctx.Error(gtype.ErrInternal, fmt.Sprintf("服务(%s)不存在", argument.Name))
		return
	}

	err = s.restart(argument.Name)
	if err != nil {
		ctx.Error(gtype.ErrInternal, err)
		return
	}

	svcStatus := &model.ServiceStatus{Name: argument.Name}
	svcStatus.Status, err = s.getStatus(argument.Name)
	if err == nil {
		go s.writeOptMessage(socket.WSSvcStatusChanged, svcStatus)
	}

	ctx.Success(info)
}

func (s *Service) RestartNginxDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogNginx)
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

func (s *Service) GetNginxes(ctx gtype.Context, ps gtype.Params) {
	results := make([]*model.ServiceNginxInfo, 0)

	if s.cfg != nil {
		items := s.cfg.Sys.Svc.Nginxes
		c := len(items)
		for i := 0; i < c; i++ {
			item := items[i]
			if item == nil {
				continue
			}
			if len(item.ServiceName) < 1 {
				continue
			}

			result := &model.ServiceNginxInfo{
				Name:        item.Name,
				ServiceName: item.ServiceName,
				Remark:      item.Remark,
				Locations:   make([]*model.ServiceNginxLocation, 0),
			}
			if len(result.Name) < 1 {
				result.Name = result.ServiceName
			}
			result.Status, _ = s.getStatus(result.ServiceName)

			lc := len(item.Locations)
			for li := 0; li < lc; li++ {
				l := item.Locations[li]
				if l == nil {
					continue
				}

				location := &model.ServiceNginxLocation{
					Name: l.Name,
					Root: l.Root,
					Urls: make([]string, 0),
				}
				location.Version, location.DeployTime, _ = s.getNginxAppInfo(l.Root)

				uc := len(l.Urls)
				for ui := 0; ui < uc; ui++ {
					location.Urls = append(location.Urls, l.Urls[ui])
				}

				result.Locations = append(result.Locations, location)
			}

			results = append(results, result)
		}
	}

	ctx.Success(results)
}

func (s *Service) GetNginxesDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogNginx)
	function := catalog.AddFunction(method, uri, "获取服务列表")
	function.SetOutputDataExample([]*model.ServiceNginxInfo{
		{
			Name:        "example",
			ServiceName: "nginx",
			Locations: []*model.ServiceNginxLocation{
				{
					Name: "api",
					Urls: []string{},
				},
			},
		},
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) ModNginxApp(ctx gtype.Context, ps gtype.Params) {
	svcName := strings.TrimSpace(ctx.Request().FormValue("svcName"))
	if len(svcName) < 1 {
		ctx.Error(gtype.ErrInput, "服务名称(svcName)为空")
		return
	}
	appName := strings.TrimSpace(ctx.Request().FormValue("appName"))
	if len(appName) < 1 {
		ctx.Error(gtype.ErrInput, "站点名称(appName)为空")
		return
	}

	info := s.cfg.Sys.Svc.GetNginxByServiceName(svcName)
	if info == nil {
		ctx.Error(gtype.ErrInternal, fmt.Sprintf("服务名称(%s)不存在", svcName))
		return
	}

	app := info.GetLocationByName(appName)
	if info == nil {
		ctx.Error(gtype.ErrInternal, fmt.Sprintf("站点名称(%s)不存在", appName))
		return
	}

	rootFolder := app.Root
	if len(rootFolder) < 1 {
		ctx.Error(gtype.ErrInternal, "站点服务物理根路径为空")
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

	appFolder := rootFolder
	err = os.RemoveAll(appFolder)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "删除原应用程序失败: ", err)
		return
	}

	err = os.Rename(tempFolder, appFolder)
	if err != nil {
		ctx.Error(gtype.ErrInternal, "重命名文件夹失败: ", err)
		return
	}

	argument := &model.ServiceNginxArgument{
		SvcName: svcName,
		AppName: appName,
	}
	argument.Version, argument.DeployTime, _ = s.getNginxAppInfo(rootFolder)

	go s.writeOptMessage(socket.WSNginxAppUpdated, argument)

	ctx.Success(argument)
}

func (s *Service) ModNginxAppDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogNginx)
	function := catalog.AddFunction(method, uri, "上传站点程序")
	function.SetNote("上传应用程序打包文件(.war, .zip或.tar.gz)，并解压缩到站点根目录下，成功时返回服务名称及站点名称信息")
	function.SetRemark("压缩包内的文件应为网站内容文件，不要嵌套在文件夹中")
	function.AddInputHeader(true, "content-type", "内容类型", gtype.ContentTypeFormData)
	function.AddInputForm(true, "svcName", "服务名称", gtype.FormValueKindText, "")
	function.AddInputForm(true, "appName", "站点名称", gtype.FormValueKindText, "")
	function.AddInputForm(true, "file", "应用程序打包文件(.war, .zip或.tar.gz)", gtype.FormValueKindFile, nil)
	function.SetOutputDataExample(&model.ServiceNginxArgument{
		SvcName: "nginx",
		AppName: "api",
	})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s *Service) GetNginxAppDetail(ctx gtype.Context, ps gtype.Params) {
	argument := &model.ServiceNginxDetailArgument{}
	err := ctx.GetJson(argument)
	if err != nil {
		ctx.Error(gtype.ErrInput, err)
		return
	}
	if len(argument.SvcName) < 1 {
		ctx.Error(gtype.ErrInternal, "服务名称(svcName)为空")
		return
	}
	if len(argument.AppName) < 1 {
		ctx.Error(gtype.ErrInternal, "站点名称(appName)为空")
		return
	}

	info := s.cfg.Sys.Svc.GetNginxByServiceName(argument.SvcName)
	if info == nil {
		ctx.Error(gtype.ErrInternal, fmt.Sprintf("服务名称(%s)不存在", argument.SvcName))
		return
	}

	app := info.GetLocationByName(argument.AppName)
	if app == nil {
		ctx.Error(gtype.ErrInternal, fmt.Sprintf("站点名称(%s)不存在", argument.AppName))
		return
	}

	cfg := &model.FileInfo{}
	s.getFileInfos(cfg, app.Root)

	cfg.Sort()
	ctx.Success(cfg)
}

func (s *Service) GetNginxAppDetailDoc(doc gtype.Doc, method string, uri gtype.Uri) {
	catalog := s.createCatalog(doc, svcCatalogRoot, svcCatalogNginx)
	function := catalog.AddFunction(method, uri, "获取站点程序详细信息")
	function.SetInputJsonExample(&model.ServiceNginxDetailArgument{
		SvcName: "nginx",
		AppName: "api",
	})
	function.SetOutputDataExample(&model.FileInfo{})
	function.AddOutputError(gtype.ErrTokenEmpty)
	function.AddOutputError(gtype.ErrTokenInvalid)
	function.AddOutputError(gtype.ErrInternal)
	function.AddOutputError(gtype.ErrInput)
}

func (s Service) getNginxAppInfo(folder string) (version, deployTime string, err error) {
	version = ""
	deployTime = ""
	err = nil
	if len(folder) < 1 {
		return
	}
	fi, fe := os.Stat(folder)
	if fe != nil {
		err = fe
		return
	}
	if !fi.IsDir() {
		return
	}
	deployTime = fi.ModTime().Format("2006-01-02 15:04:05")

	version, _ = s.getTextVersion(folder)
	if len(version) < 1 {
		version, _ = s.getJsonVersion(folder)
		if len(version) < 1 {
			version = s.getTomcatAppVersion(folder)
		}
	}

	return
}
