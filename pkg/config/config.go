package config

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/vgarvardt/rklotz/pkg/plugin"
	"github.com/vgarvardt/rklotz/pkg/renderer"
)

// Config represents the application configuration
type Config struct {
	LogLevel     string `envconfig:"LOG_LEVEL" default:"info"`
	PostsDSN     string `envconfig:"POSTS_DSN" default:"file:///etc/rklotz/posts"`
	PostsPerPage int    `envconfig:"POSTS_PERPAGE" default:"10"`
	StorageDSN   string `envconfig:"STORAGE_DSN" default:"boltdb:///tmp/rklotz.db"`

	Web
	SSL
	plugin.Config
	renderer.UIConfig
	renderer.RootURLConfig
}

// Web is the configuration for web application
type Web struct {
	Port          int    `envconfig:"WEB_PORT" default:"8080"`
	StaticPath    string `envconfig:"WEB_STATIC_PATH" default:"/etc/rklotz/static"`
	TemplatesPath string `envconfig:"WEB_TEMPLATES_PATH" default:"/etc/rklotz/templates"`
}

// SSL is the configuration for TLS/SSL
type SSL struct {
	Enabled  bool   `envconfig:"SSL_ENABLED" default:"false"`
	Port     int    `envconfig:"SSL_PORT" default:"8443"`
	Host     string `envconfig:"SSL_HOST"`
	Email    string `envconfig:"SSL_EMAIL" default:"vgarvardt@gmail.com"`
	CacheDir string `envconfig:"SSL_CACHE_DIR" default:"/tmp"`
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
