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
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/repository"
)

func RunServer(cmd *cobra.Command, args []string) {
	appConfig, err := config.Load()
	failOnError(err, "Failed to load config")

	logLevel, err := log.ParseLevel(appConfig.LogLevel)
	failOnError(err, "Failed to parse log level")
	log.SetLevel(logLevel)

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	instanceId := hex.EncodeToString(hasher.Sum(nil))[:5]

	log.WithFields(log.Fields{"version": version, "instance": instanceId}).Info("Starting rKlotz...")

	storage, err := repository.NewStorage(appConfig.StorageDSN, appConfig.PostsPerPage)
	failOnError(err, "Failed to get storage instance")
	defer storage.Close()

	loader, err := repository.NewLoader(appConfig.PostsDSN)
	failOnError(err, "Failed to get loader instance")

	err = loader.Load(storage)
	failOnError(err, "Failed to load posts")

	htmlRenderer, err := renderer.NewHTMLRenderer(appConfig.Web.TemplatesPath, instanceId, appConfig.UI, appConfig.Plugins)
	failOnError(err, "Failed to init HTML Renderer")
	xmlRenderer := renderer.NewXmlRenderer()

	postsHandler := handler.NewPostsHandler(storage, htmlRenderer)
	feedHandler := handler.NewFeedHandler(storage, xmlRenderer, appConfig.UI, appConfig.RootURL)

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

func failOnError(err error, msg string) {
	if nil != err {
		log.WithError(err).Panic(msg)
	}
}
