package main

import (
	"fmt"
	"github.com/csby/gwsf/gopt"
	"github.com/csby/gwsf/gtype"
	"net/http"
)

func NewHandler(log gtype.Log) gtype.Handler {
	instance := &Handler{}
	instance.SetLog(log)

	return instance
}

type Handler struct {
	gtype.Base

	ctrl Controllers
}

func (s *Handler) InitRouting(router gtype.Router) {
}

func (s *Handler) BeforeRouting(ctx gtype.Context) {
	method := ctx.Method()
	// enable across access
	if method == http.MethodOptions {
		ctx.Response().Header().Add("Access-Control-Allow-Origin", "*")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "content-type,token")
		ctx.SetHandled(true)
		return
	}

	if method == http.MethodGet {
		schema := ctx.Schema()
		path := ctx.Path()

		// default to opt site
		if "/" == path || "" == path || gopt.WebPath == path || "/omw" == path || "/omw/" == path || "/omw/#/" == path {
			redirectUrl := fmt.Sprintf("%s://%s%s/", schema, ctx.Host(), gopt.WebPath)
			http.Redirect(ctx.Response(), ctx.Request(), redirectUrl, http.StatusMovedPermanently)
			ctx.SetHandled(true)
			return
		}

		// http to https
		if method == "http" {
			if cfg.Http.RedirectToHttps && cfg.Https.Enabled {
				redirectUrl := fmt.Sprintf("%s://%s%s", "https", ctx.Host(), path)
				http.Redirect(ctx.Response(), ctx.Request(), redirectUrl, http.StatusMovedPermanently)
				ctx.SetHandled(true)
				return
			}
		}
	}
}

func (s *Handler) AfterRouting(ctx gtype.Context) {

}

func (s *Handler) ExtendOptSetup(opt gtype.Option) {
	if opt != nil {
		opt.SetCloud(cfg.Cloud.Enabled)
		opt.SetNode(cfg.Node.Enabled)
	}
}

func (s *Handler) ExtendOptApi(router gtype.Router,
	path *gtype.Path,
	preHandle gtype.HttpHandle,
	opt gtype.Opt) {
	s.ctrl.initController(opt.Wsc())
	s.ctrl.initRouter(router, path, preHandle)
}
