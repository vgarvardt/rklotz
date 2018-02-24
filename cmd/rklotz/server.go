package main

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/spf13/cobra"
	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/handler"
	"github.com/vgarvardt/rklotz/pkg/loader"
	m "github.com/vgarvardt/rklotz/pkg/middleware"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/acme/autocert"
)

var logger *zap.Logger

// RunServer initializes and runs web-server instance
func RunServer(cmd *cobra.Command, args []string) {
	appConfig, err := config.Load()
	failOnError(err, "Failed to load config")

	logger, err = initLogger(appConfig.LogLevel)
	failOnError(err, "Failed to initialize logger")
	defer logger.Sync()

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	instanceID := hex.EncodeToString(hasher.Sum(nil))[:5]

	logger.Info("Starting rKlotz...", zap.String("version", version), zap.String("instance", instanceID))

	storageInstance, err := storage.NewStorage(appConfig.StorageDSN, appConfig.PostsPerPage)
	failOnError(err, "Failed to get storageInstance instance")
	defer storageInstance.Close()

	loaderInstance, err := loader.NewLoader(appConfig.PostsDSN, logger)
	failOnError(err, "Failed to get loaderInstance instance")

	err = loaderInstance.Load(storageInstance)
	failOnError(err, "Failed to load posts")

	htmlRenderer, err := renderer.NewHTMLRenderer(
		renderer.HTMLRendererConfig{
			TemplatesPath: appConfig.Web.TemplatesPath,
			InstanceID:    instanceID,
			UISettings:    appConfig.UI,
			Plugins:       appConfig.Plugins,
			RootURL:       appConfig.RootURL,
		},
		logger,
	)
	failOnError(err, "Failed to init HTML Renderer")
	xmlRenderer := renderer.NewXMLRenderer()

	postsHandler := handler.NewPostsHandler(storageInstance, htmlRenderer)
	feedHandler := handler.NewFeedHandler(storageInstance, xmlRenderer, appConfig.UI, appConfig.RootURL)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(m.NewLogger(logger).Handler)
	r.Use(middleware.RequestLogger(m.NewRequestLogger()))
	r.Use(middleware.Recoverer)

	r.Get("/", postsHandler.Front)
	r.Get("/tag/{tag}", postsHandler.Tag)
	r.NotFound(postsHandler.Post)

	r.Route("/feed", func(r chi.Router) {
		r.Get("/atom", feedHandler.Atom)
		r.Get("/rss", feedHandler.Rss)
	})

	serveStatic(r, appConfig)

	listenAndServe(r, appConfig)
}

func serveStatic(r chi.Router, appConfig *config.AppConfig) {
	staticRoot := http.Dir(appConfig.Web.StaticPath)
	staticPath := "/static"
	staticHandler := http.StripPrefix(staticPath, http.FileServer(staticRoot))

	faviconPath := filepath.Join(appConfig.Web.StaticPath, appConfig.UI.Theme, "favicon.ico")

	r.Get(staticPath+"/*", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})
}

func listenAndServe(handler chi.Router, appConfig *config.AppConfig) {
	if appConfig.SSL.Enabled {
		logger.Info("SSL is enabled, starting HTTPS server")

		tlsCertManager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(appConfig.SSL.Host),
			Cache:      autocert.DirCache(appConfig.SSL.CacheDir),
			Email:      appConfig.UI.Email,
		}

		httpsAddress := fmt.Sprintf(":%d", appConfig.SSL.Port)
		server := &http.Server{
			Addr:    httpsAddress,
			Handler: handler,
			TLSConfig: &tls.Config{
				GetCertificate: tlsCertManager.GetCertificate,
			},
		}

		go func() {
			logger.Info("Running HTTPS server...", zap.String("address", httpsAddress))
			logger.Fatal("Failed to run HTTPS server", zap.Error(server.ListenAndServeTLS("", "")))
		}()

		httpAddress := fmt.Sprintf(":%d", appConfig.Web.Port)

		logger.Info("Running HTTP to HTTPS redirect server...", zap.String("address", httpAddress))
		logger.Fatal(
			"Failed to run HTTP to HTTPS redirect server",
			zap.Error(http.ListenAndServe(httpAddress, tlsCertManager.HTTPHandler(nil))),
		)
	} else {
		address := fmt.Sprintf(":%d", appConfig.Web.Port)
		logger.Info("Running HTTP server...", zap.String("address", address))
		logger.Fatal("Failed to run HTTP server", zap.Error(http.ListenAndServe(address, handler)))
	}
}

func initLogger(level string) (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()

	logLevel := new(zap.AtomicLevel)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}

	logConfig.Development = logLevel.String() == zapcore.DebugLevel.String()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	logger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func failOnError(err error, msg string) {
	if nil != err {
		if logger != nil {
			logger.Panic(msg, zap.Error(err))
		} else {
			log.Panic()
		}
	}
}
