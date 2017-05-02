package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/spf13/cobra"
	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/handler"
	m "github.com/vgarvardt/rklotz/pkg/middleware"
	"github.com/vgarvardt/rklotz/pkg/model"
	"github.com/vgarvardt/rklotz/pkg/renderer"
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

	htmlRenderer := renderer.NewHTMLRenderer(appConfig.Web.TemplatesPath, instanceId, appConfig.UI)
	xmlRenderer := renderer.NewXmlRenderer()

	postsHandler := handler.NewPostsHandler(htmlRenderer)
	feedHandler := handler.NewFeedHandler(xmlRenderer, appConfig.UI, appConfig.RootURL)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&m.LoggerRequest{}))
	r.Use(middleware.Recoverer)

	r.Get("/", postsHandler.Front)
	r.Get("/tag/:tag", postsHandler.Tag)
	r.NotFound(postsHandler.Post)

	r.Route("/feed", func(r chi.Router) {
		r.Get("/atom", feedHandler.Atom)
		r.Get("/rss", feedHandler.Rss)
	})

	r.FileServer("/static", http.Dir(appConfig.Web.StaticPath))
	faviconPath := filepath.Join(appConfig.Web.StaticPath, appConfig.UI.Theme, "favicon.ico")
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})

	address := fmt.Sprintf(":%d", appConfig.Web.Port)
	log.WithField("address", address).Info("Running...")

	log.Fatal(http.ListenAndServe(address, r))
}
