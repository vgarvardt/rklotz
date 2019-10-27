package server

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	wErrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/acme/autocert"

	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/handler"
	"github.com/vgarvardt/rklotz/pkg/loader"
	m "github.com/vgarvardt/rklotz/pkg/middleware"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// Run initializes and runs web-server instance
func Run(cfg *config.Config, version string) error {
	logger, err := initLogger(cfg.LogLevel)
	if err != nil {
		return wErrors.Wrap(err, "failed to initialize logger")
	}
	defer logger.Sync()

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	instanceID := hex.EncodeToString(hasher.Sum(nil))[:5]

	logger.Info("Starting rKlotz...", zap.String("version", version), zap.String("instance", instanceID))

	storageInstance, err := storage.NewStorage(cfg.StorageDSN, cfg.PostsPerPage)
	if err != nil {
		return wErrors.Wrap(err, "failed to initialise storage")
	}
	defer func() {
		if err := storageInstance.Close(); err != nil {
			logger.Error("Got an error while closing storage", zap.Error(err))
		}
	}()

	loaderInstance, err := loader.NewLoader(cfg.PostsDSN, logger)
	if err != nil {
		return wErrors.Wrap(err, "failed to initialise loader")
	}

	err = loaderInstance.Load(storageInstance)
	if err != nil {
		return wErrors.Wrap(err, "failed to load posts")
	}

	htmlRenderer, err := renderer.NewHTMLRenderer(
		renderer.HTMLRendererConfig{
			TemplatesPath: cfg.Web.TemplatesPath,
			InstanceID:    instanceID,
			UICfg:         cfg.UI,
			PluginsCfg:    cfg.Plugins,
			RootURLCfg:    cfg.RootURL,
		},
		logger,
	)
	if err != nil {
		return wErrors.Wrap(err, "failed to initialise HTML renderer")
	}

	xmlRenderer := renderer.NewXMLRenderer()

	postsHandler := handler.NewPostsHandler(storageInstance, htmlRenderer)
	feedHandler := handler.NewFeedHandler(storageInstance, xmlRenderer, cfg.UI, cfg.RootURL)

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

	serveStatic(r, cfg.Web, cfg.UI.Theme)

	return listenAndServe(r, cfg.SSL, cfg.Web, logger)
}

func serveStatic(r chi.Router, cfgWeb config.Web, theme string) {
	staticRoot := http.Dir(cfgWeb.StaticPath)
	staticPath := "/static"
	staticHandler := http.StripPrefix(staticPath, http.FileServer(staticRoot))

	faviconPath := filepath.Join(cfgWeb.StaticPath, theme, "favicon.ico")

	r.Get(staticPath+"/*", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})
}

func listenAndServe(handler chi.Router, cfgSSL config.SSL, cfgWeb config.Web, logger *zap.Logger) error {
	if !cfgSSL.Enabled {
		address := fmt.Sprintf(":%d", cfgWeb.Port)
		logger.Info("Running HTTP server...", zap.String("address", address))

		return http.ListenAndServe(address, handler)
	}

	logger.Info("SSL is enabled, starting HTTPS server")

	tlsCertManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfgSSL.Host),
		Cache:      autocert.DirCache(cfgSSL.CacheDir),
		Email:      cfgSSL.Email,
	}

	httpsAddress := fmt.Sprintf(":%d", cfgSSL.Port)
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

	httpAddress := fmt.Sprintf(":%d", cfgWeb.Port)

	logger.Info("Running HTTP to HTTPS redirect server...", zap.String("address", httpAddress))

	return http.ListenAndServe(httpAddress, tlsCertManager.HTTPHandler(nil))
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
