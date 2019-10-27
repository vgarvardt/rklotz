package server

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/vgarvardt/rklotz/pkg/plugin"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/server/web"
)

// Config represents server configuration
type Config struct {
	PostsDSN     string `envconfig:"POSTS_DSN" default:"file:///etc/rklotz/posts"`
	PostsPerPage int    `envconfig:"POSTS_PERPAGE" default:"10"`
	StorageDSN   string `envconfig:"STORAGE_DSN" default:"boltdb:///tmp/rklotz.db"`

	LogConfig
	web.HTTPConfig
	web.SSLConfig
	plugin.Config
	renderer.UIConfig
	renderer.RootURLConfig
}

// LogConfig represents logger configuration
type LogConfig struct {
	Level string `envconfig:"LOG_LEVEL" default:"info"`
	Type  string `envconfig:"LOG_TYPE" default:"rklotz"`
}

// BuildLogger builds and initialises logger with the values from the config
func (c *LogConfig) BuildLogger() (*zap.Logger, error) {
	logConfig := zap.NewProductionConfig()

	logLevel := new(zap.AtomicLevel)
	if err := logLevel.UnmarshalText([]byte(c.Level)); err != nil {
		return nil, err
	}

	logConfig.Level = *logLevel
	logConfig.Development = logLevel.String() == zapcore.DebugLevel.String()
	logConfig.Sampling = nil
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	logConfig.InitialFields = map[string]interface{}{"type": c.Type}

	logger, err := logConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// LoadConfig loads app settings from environment variables
func LoadConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
