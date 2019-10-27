package config

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/vgarvardt/rklotz/pkg/plugin"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/server/web"
)

// Config represents the application configuration
type Config struct {
	LogLevel     string `envconfig:"LOG_LEVEL" default:"info"`
	PostsDSN     string `envconfig:"POSTS_DSN" default:"file:///etc/rklotz/posts"`
	PostsPerPage int    `envconfig:"POSTS_PERPAGE" default:"10"`
	StorageDSN   string `envconfig:"STORAGE_DSN" default:"boltdb:///tmp/rklotz.db"`

	web.HTTPConfig
	web.SSLConfig
	plugin.Config
	renderer.UIConfig
	renderer.RootURLConfig
}

// Load loads app settings from environment variables
func Load() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
