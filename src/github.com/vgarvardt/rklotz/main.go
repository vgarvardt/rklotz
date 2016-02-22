package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/vgarvardt/rklotz/cfg"
	"github.com/vgarvardt/rklotz/controller"
	"github.com/vgarvardt/rklotz/model"
)

func main() {
	defer model.GetDB().Close()

	if len(cfg.GetOptions().Update) > 0 {
		cfg.Log(fmt.Sprintf("Trying to update post UUID %s", cfg.GetOptions().Update))
		if err := model.UpdatePostField(cfg.GetOptions().Update, cfg.GetOptions().Field, cfg.GetOptions().Value); err != nil {
			panic(err)
		}
	}

	if err := model.RebuildIndex(); err != nil {
		panic(err)
	}

	if cfg.GetRunWebServer() {
		if cfg.Bool("debug") {
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
			cfg.String("auth.name"): cfg.String("auth.password"),
		}))
		authorized.GET("/new", controller.FormController)
		authorized.POST("/new", controller.FormController)
		authorized.GET("/edit/:uuid", controller.FormController)
		authorized.POST("/edit/:uuid", controller.FormController)
		authorized.GET("/drafts", controller.DraftsController)
		authorized.GET("/published", controller.PublishedController)

		router.NoRoute(controller.PostController)

		router.Static("/static", fmt.Sprintf("%s/static", cfg.GetRootDir()))

		addr := cfg.String("addr")
		cfg.Log(fmt.Sprintf("Running @ %s", addr))
		router.Run(addr)
	}
}
