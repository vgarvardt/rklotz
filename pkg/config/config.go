package config

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/structs"
	"github.com/kelseyhightower/envconfig"

	"github.com/vgarvardt/rklotz/pkg/config/plugin"
)

// Config represents the application configuration
type Config struct {
	LogLevel     string `envconfig:"LOG_LEVEL" default:"info"`
	PostsDSN     string `envconfig:"POSTS_DSN" default:"file:///etc/rklotz/posts"`
	PostsPerPage int    `envconfig:"POSTS_PERPAGE" default:"10"`
	StorageDSN   string `envconfig:"STORAGE_DSN" default:"boltdb:///tmp/rklotz.db"`

	Web     WebSettings
	SSL     SSLSettings
	UI      UISetting
	RootURL RootURL
	Plugins Plugins
}

// WebSettings is the configuration for web application
type WebSettings struct {
	Port          int    `envconfig:"WEB_PORT" default:"8080"`
	StaticPath    string `envconfig:"WEB_STATIC_PATH" default:"/etc/rklotz/static"`
	TemplatesPath string `envconfig:"WEB_TEMPLATES_PATH" default:"/etc/rklotz/templates"`
}

// SSLSettings is the configuration for TLS/SSL
type SSLSettings struct {
	Enabled  bool   `envconfig:"SSL_ENABLED" default:"false"`
	Port     int    `envconfig:"SSL_PORT" default:"8443"`
	Host     string `envconfig:"SSL_HOST"`
	Email    string `envconfig:"SSL_EMAIL" default:"vgarvardt@gmail.com"`
	CacheDir string `envconfig:"SSL_CACHE_DIR" default:"/tmp"`
}

// UISetting is the configuration for user interface
type UISetting struct {
	Theme       string `envconfig:"UI_THEME" default:"foundation"`
	Author      string `envconfig:"UI_AUTHOR" default:"Vladimir Garvardt"`
	Email       string `envconfig:"UI_EMAIL" default:"vgarvardt@gmail.com"`
	Description string `envconfig:"UI_DESCRIPTION" default:"rKlotz - simple golang-driven blog engine"`
	Language    string `envconfig:"UI_LANGUAGE" default:"en"`
	Title       string `envconfig:"UI_TITLE" default:"rKlotz"`
	Heading     string `envconfig:"UI_HEADING" default:"rKlotz"`
	Intro       string `envconfig:"UI_INTRO" default:"simple golang-driven blog engine"`
	// DateFormat is format for posts, see http://golang.org/pkg/time/#Time.Format
	DateFormat string `envconfig:"UI_DATEFORMAT" default:"2 Jan 2006"`
	AboutPath  string `envconfig:"UI_ABOUT_PATH" default:"/etc/rklotz/about.html"`
}

// RootURL is the configuration for app root url
type RootURL struct {
	Scheme string `envconfig:"ROOT_URL_SCHEME" default:"http"`
	Host   string `envconfig:"ROOT_URL_HOST"`
	Path   string `envconfig:"ROOT_URL_PATH" default:"/"`
}

// URL returns the URL for currently configured root url
func (u RootURL) URL(r *http.Request) *url.URL {
	host := u.Host
	if len(host) < 1 {
		host = r.Host
	}
	return &url.URL{Scheme: u.Scheme, Host: host, Path: u.Path}
}

// Plugins is teh configuration for app plugins
type Plugins struct {
	Enabled  []string `envconfig:"PLUGINS_ENABLED"`
	Settings PluginsSettings
}

// Configure applies configuration for enabled plugins
func (p Plugins) Configure(instance plugin.Plugin) (map[string]string, error) {
	pluginName, err := plugin.GetName(instance)
	if err != nil {
		return nil, err
	}

	settingsMap := structs.Map(p.Settings)
	pluginSettings, ok := settingsMap[strings.Title(pluginName)]
	if !ok {
		return nil, errors.New("failed to get plugin settings")
	}

	pluginSettingsMap := pluginSettings.(map[string]string)
	if len(pluginSettingsMap) == 0 {
		return instance.Configure(instance.Defaults())
	}

	return instance.Configure(pluginSettingsMap)
}

// PluginsSettings is the configuration for available plugins
type PluginsSettings struct {
	Disqus      map[string]string `envconfig:"PLUGINS_DISQUS"`
	Ga          map[string]string `envconfig:"PLUGINS_GA"`
	Gtm         map[string]string `envconfig:"PLUGINS_GTM"`
	Yamka       map[string]string `envconfig:"PLUGINS_YAMKA"`
	Highlightjs map[string]string `envconfig:"PLUGINS_HIGHLIGHTJS"`
	Yasha       map[string]string `envconfig:"PLUGINS_YASHA"`
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
