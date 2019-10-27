package server

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	wErrors "github.com/pkg/errors"
	"github.com/vgarvardt/rklotz/pkg/handler"
	"github.com/vgarvardt/rklotz/pkg/loader"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/server/web"
	"github.com/vgarvardt/rklotz/pkg/storage"
	"go.uber.org/zap"
)

// Run initializes and runs web-server instance
func Run(cfg *Config, version string) error {
	logger, err := cfg.LogConfig.BuildLogger()

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

	htmlRenderer, err := renderer.NewHTML(
		renderer.HTMLConfig{
			TemplatesPath: cfg.HTTPConfig.TemplatesPath,
			InstanceID:    instanceID,
			UICfg:         cfg.UIConfig,
			PluginsCfg:    cfg.Config,
			RootURLCfg:    cfg.RootURLConfig,
		},
		logger,
	)
	if err != nil {
		return wErrors.Wrap(err, "failed to initialise HTML renderer")
	}

	feedRenderer := renderer.NewFeed(cfg.UIConfig, cfg.RootURLConfig)

	postsHandler := handler.NewPosts(storageInstance, htmlRenderer)
	feedHandler := handler.NewFeed(storageInstance, feedRenderer)

	r := web.NewRouter(postsHandler, feedHandler, logger)

	web.ServeStatic(r, cfg.HTTPConfig, cfg.UIConfig.Theme)

	return web.ListenAndServe(r, cfg.SSLConfig, cfg.HTTPConfig, logger)
}
