package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cappuccinotm/slogx"

	"github.com/vgarvardt/rklotz/pkg/loader"
	"github.com/vgarvardt/rklotz/pkg/server/handler"
	"github.com/vgarvardt/rklotz/pkg/server/renderer"
	"github.com/vgarvardt/rklotz/pkg/server/web"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// Run initializes and runs web-server instance
func Run(ctx context.Context, cfg *Config, version string) error {
	logger, err := cfg.LogConfig.BuildLogger()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	logger.Info("Starting rKlotz...", slog.String("version", version))

	storageInstance, err := storage.NewStorage(cfg.StorageDSN, cfg.PostsPerPage)
	if err != nil {
		return fmt.Errorf("failed to initialise storage: %w", err)
	}
	defer func() {
		if err := storageInstance.Close(); err != nil {
			logger.Error("Got an error while closing storage", slogx.Error(err))
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
			Debug:         cfg.LogConfig.Level == slog.LevelDebug,
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

	return web.ListenAndServe(ctx, r, cfg.SSLConfig, cfg.HTTPConfig, logger)
}
