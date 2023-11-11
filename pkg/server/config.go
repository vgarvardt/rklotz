package server

import (
	"context"
	"log/slog"
	"os"

	"github.com/cappuccinotm/slogx"
	"github.com/cappuccinotm/slogx/slogm"
	"github.com/sethvargo/go-envconfig"

	"github.com/vgarvardt/rklotz/pkg/server/plugin"
	"github.com/vgarvardt/rklotz/pkg/server/renderer"
	"github.com/vgarvardt/rklotz/pkg/server/web"
)

// Config represents server configuration
type Config struct {
	PostsDSN     string `env:"POSTS_DSN,default=file:///etc/rklotz/posts"`
	PostsPerPage int    `env:"POSTS_PERPAGE,default=10"`
	StorageDSN   string `env:"STORAGE_DSN,default=boltdb:///tmp/rklotz.db"`

	LogConfig
	web.HTTPConfig
	web.SSLConfig
	plugin.Config
	renderer.UIConfig
	renderer.RootURLConfig
}

// LogConfig represents logger configuration
type LogConfig struct {
	Level slog.Level `env:"LOG_LEVEL,default=info"`
	Type  string     `env:"LOG_TYPE,default=rklotz"`
}

// BuildLogger builds and initialises logger with the values from the config
func (c *LogConfig) BuildLogger() (*slog.Logger, error) {
	so := &slog.HandlerOptions{
		AddSource: true,
		Level:     c.Level,
	}

	logger := slog.
		New(slogx.NewChain(slog.NewJSONHandler(os.Stderr, so), slogm.StacktraceOnError())).
		With(slog.String("type", c.Type))

	return logger, nil
}

// LoadConfig loads app settings from environment variables
func LoadConfig(ctx context.Context) (*Config, error) {
	var cfg Config

	err := envconfig.Process(ctx, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
