package main

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"

	"github.com/vgarvardt/rklotz/pkg/app"
	"github.com/vgarvardt/rklotz/pkg/controller"
	"github.com/vgarvardt/rklotz/pkg/model"
	"github.com/vgarvardt/rklotz/pkg/svc"
)

func main() {
	defer model.DB.Close()

	config := svc.Container.MustGet(svc.DI_CONFIG).(svc.Config)

	e := echo.New()
	e.SetDebug(config.Bool("debug"))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.SetRenderer(svc.Renderer(app.RootDir(), app.InstanceId()))

	e.GET("/", controller.FrontController)
	e.GET("/tag/:tag", controller.TagController)
	e.GET("/autocomplete", controller.AutoComplete)
	e.GET("/*", controller.PostController)

	feed := e.Group("/feed")
	feed.GET("/atom", controller.AtomController)
	feed.GET("/rss", controller.RssController)

	authorized := e.Group("/@", middleware.BasicAuth(func(username, password string) bool {
		return username == config.String("auth.name") && password == config.String("auth.password")
	}))
	authorized.GET("/", controller.AdmFrontController)
	authorized.GET("/new", controller.FormController)
	authorized.POST("/new", controller.FormController)
	authorized.GET("/edit/:uuid", controller.FormController)
	authorized.POST("/edit/:uuid", controller.FormController)
	authorized.GET("/drafts", controller.DraftsController)
	authorized.GET("/published", controller.PublishedController)

	e.Static("/static", fmt.Sprintf("%s/static", app.RootDir()))
	e.File("/favicon.ico", fmt.Sprintf("%s/static/images/favicon.ico", app.RootDir()))

	addr := config.String("addr")
	std := standard.New(addr)
	std.SetHandler(e)
	log.WithField("address", addr).Info("Running...")

	gracehttp.Serve(std.Server)
}
