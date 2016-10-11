package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/Sirupsen/logrus"

	"github.com/vgarvardt/rklotz/app"
	"github.com/vgarvardt/rklotz/controller"
	"github.com/vgarvardt/rklotz/model"
	"github.com/vgarvardt/rklotz/svc"
)

func main() {
	defer model.GetDB().Close()

	logger := svc.Container.MustGet(svc.DI_LOGGER).(*log.Logger)
	config := svc.Container.MustGet(svc.DI_CONFIG).(svc.Config)

	switch app.Command() {
	case app.COMMAND_UPDATE:
		updateParams := app.GetUpdateParams()
		logger.WithField("UUID", updateParams.UUID).Info("Trying to update post")
		if err := model.UpdatePostField(updateParams.UUID, updateParams.Field, updateParams.Value); err != nil {
			panic(err)
		}

	case app.COMMAND_REBUILD:
		if err := model.RebuildIndex(); err != nil {
			panic(err)
		}

	case app.COMMAND_RUN:
		if config.Bool("debug") {
			gin.SetMode(gin.DebugMode)
		} else {
			gin.SetMode(gin.ReleaseMode)
		}

		router := gin.Default()

		router.GET("/", controller.FrontController)
		router.GET("/tag/:tag", controller.TagController)
		router.GET("/autocomplete", controller.AutoComplete)

		feed := router.Group("/feed")
		feed.GET("/atom", controller.AtomController)
		feed.GET("/rss", controller.RssController)

		authorized := router.Group("/@", gin.BasicAuth(gin.Accounts{
			config.String("auth.name"): config.String("auth.password"),
		}))
		authorized.GET("/new", controller.FormController)
		authorized.POST("/new", controller.FormController)
		authorized.GET("/edit/:uuid", controller.FormController)
		authorized.POST("/edit/:uuid", controller.FormController)
		authorized.GET("/drafts", controller.DraftsController)
		authorized.GET("/published", controller.PublishedController)

		router.NoRoute(controller.PostController)

		router.Static("/static", fmt.Sprintf("%s/static", app.RootDir()))

		addr := config.String("addr")
		logger.WithField("address", addr).Info("Running...")
		router.Run(addr)
	}
}
