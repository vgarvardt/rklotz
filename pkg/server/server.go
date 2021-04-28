package server

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vgarvardt/rklotz/pkg/loader"
	"github.com/vgarvardt/rklotz/pkg/server/handler"
	"github.com/vgarvardt/rklotz/pkg/server/renderer"
	"github.com/vgarvardt/rklotz/pkg/server/web"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// Run initializes and runs web-server instance
func Run(cfg *Config, version string) error {
	logger, err := cfg.LogConfig.BuildLogger()

	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			fmt.Printf("%s Could not sync logger: %v", time.Now().String(), err)
		}
	}()

	logger.Info("Starting rKlotz...", zap.String("version", version))

	storageInstance, err := storage.NewStorage(cfg.StorageDSN, cfg.PostsPerPage)
	if err != nil {
		return fmt.Errorf("failed to initialise storage: %w", err)
	}
	defer func() {
		if err := storageInstance.Close(); err != nil {
			logger.Error("Got an error while closing storage", zap.Error(err))
		}
	}()

	loaderInstance, err := loader.New(cfg.PostsDSN, logger)
	if err != nil {
		return fmt.Errorf("failed to initialise loader: %w", err)
	}

	err = loaderInstance.Load(storageInstance)
	if err != nil {
		return fmt.Errorf("failed to load posts: %w", err)
	}

	htmlRenderer, err := renderer.NewHTML(
		renderer.HTMLConfig{
			Debug:         cfg.LogConfig.Level == zapcore.DebugLevel.String(),
			TemplatesPath: cfg.HTTPConfig.TemplatesPath,
			UICfg:         cfg.UIConfig,
			PluginsCfg:    cfg.Config,
			RootURLCfg:    cfg.RootURLConfig,
		},
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to initialise HTML renderer: %w", err)
	}

	feedRenderer := renderer.NewFeed(cfg.UIConfig, cfg.RootURLConfig)

	postsHandler := handler.NewPosts(storageInstance, htmlRenderer)
	feedHandler := handler.NewFeed(storageInstance, feedRenderer)

	r := web.NewRouter(postsHandler, feedHandler, logger)

	web.ServeStatic(r, cfg.HTTPConfig, cfg.UIConfig.Theme)

	return web.ListenAndServe(r, cfg.SSLConfig, cfg.HTTPConfig, logger)
}
