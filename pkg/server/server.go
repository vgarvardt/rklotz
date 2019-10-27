package server

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	wErrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/handler"
	"github.com/vgarvardt/rklotz/pkg/loader"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/server/web"
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

	postsHandler := handler.NewPosts(storageInstance, htmlRenderer)
	feedHandler := handler.NewFeed(storageInstance, xmlRenderer, cfg.UI, cfg.RootURL)

	r := web.NewRouter(postsHandler, feedHandler, logger)

	web.ServeStatic(r, cfg.Web, cfg.UI.Theme)

	return web.ListenAndServe(r, cfg.SSL, cfg.Web, logger)
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
