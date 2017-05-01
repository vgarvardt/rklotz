package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/controller"
	"github.com/vgarvardt/rklotz/pkg/model"
	"github.com/vgarvardt/rklotz/pkg/svc"
)

func RunServer(cmd *cobra.Command, args []string) {
	defer model.DB.Close()

	appConfig, err := config.Load()
	if nil != err {
		log.WithError(err).Panic("Failed to load config")
	}

	logLevel, err := log.ParseLevel(appConfig.LogLevel)
	if nil != err {
		log.WithError(err).Panic("Failed to parse log level")
	}
	log.SetLevel(logLevel)

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format("2006/01/02 - 15:04:05")))
	instanceId := hex.EncodeToString(hasher.Sum(nil))[:5]

	log.WithFields(log.Fields{"version": version, "instance": instanceId}).Info("Starting rKlotz...")

	e := echo.New()
	e.SetDebug(false)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.SetRenderer(svc.Renderer(appConfig.Web.TemplatesPath, instanceId, appConfig.UI))

	e.GET("/", controller.FrontController)
	e.GET("/tag/:tag", controller.TagController)
	e.GET("/*", controller.PostController)

	feedController := controller.NewFeedController(appConfig.UI, appConfig.RootURL)

	feed := e.Group("/feed")
	feed.GET("/atom", feedController.AtomHandler)
	feed.GET("/rss", feedController.RssHandler)

	e.Static("/static", appConfig.Web.StaticPath)
	e.File("/favicon.ico", filepath.Join(appConfig.Web.StaticPath, appConfig.UI.Theme, "favicon.ico"))

	address := fmt.Sprintf(":%d", appConfig.Web.Port)
	std := standard.New(address)
	std.SetHandler(e)
	log.WithField("address", address).Info("Running...")

	gracehttp.Serve(std.Server)
}
