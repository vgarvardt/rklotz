package config

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/structs"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"github.com/vgarvardt/rklotz/pkg/config/plugin"
)

type AppConfig struct {
	LogLevel     string `envconfig:"LOG_LEVEL"`
	PostsDSN     string `envconfig:"POSTS_DSN"`
	PostsPerPage int    `envconfig:"POSTS_PERPAGE"`
	StorageDSN   string `envconfig:"STORAGE_DSN"`

	Web     WebSettings
	UI      UISetting
	RootURL RootURL
	Plugins Plugins
}

type WebSettings struct {
	Port          int    `envconfig:"WEB_PORT"`
	StaticPath    string `envconfig:"WEB_STATIC_PATH"`
	TemplatesPath string `envconfig:"WEB_TEMPLATES_PATH"`
}

type UISetting struct {
	Theme       string `envconfig:"UI_THEME"`
	Author      string `envconfig:"UI_AUTHOR"`
	Email       string `envconfig:"UI_EMAIL"`
	Description string `envconfig:"UI_DESCRIPTION"`
	Language    string `envconfig:"UI_LANGUAGE"`
	Title       string `envconfig:"UI_TITLE"`
	Heading     string `envconfig:"UI_HEADING"`
	Intro       string `envconfig:"UI_INTRO"`
	// DateFormat is format for posts, see http://golang.org/pkg/time/#Time.Format
	DateFormat string `envconfig:"UI_DATEFORMAT"`
	AboutPath  string `envconfig:"UI_ABOUT_PATH"`
}

type RootURL struct {
	Scheme string `envconfig:"ROOT_URL_SCHEME"`
	Host   string `envconfig:"ROOT_URL_HOST"`
	Path   string `envconfig:"ROOT_URL_PATH"`
}

func (u RootURL) URL(r *http.Request) *url.URL {
	host := u.Host
	if len(host) < 1 {
		host = r.Host
	}
	return &url.URL{Scheme: u.Scheme, Host: host, Path: u.Path}
}

type Plugins struct {
	Enabled  []string `envconfig:"PLUGINS_ENABLED"`
	Settings PluginsSettings
}

func (p Plugins) Configure(instance plugin.Plugin) (map[string]string, error) {
	pluginName, err := plugin.GetName(instance)
	if err != nil {
		return nil, err
	}

	settingsMap := structs.Map(p.Settings)
	pluginSettings, ok := settingsMap[strings.Title(pluginName)]
	if !ok {
		return nil, errors.New("Failed to get plugin settings")
	}
	return instance.Configure(pluginSettings.(map[string]string))
}

type PluginsSettings struct {
	Disqus      map[string]string `envconfig:"PLUGINS_DISQUS"`
	Ga          map[string]string `envconfig:"PLUGINS_GA"`
	Yamka       map[string]string `envconfig:"PLUGINS_YAMKA"`
	Highlightjs map[string]string `envconfig:"PLUGINS_HIGHLIGHTJS"`
	Yasha       map[string]string `envconfig:"PLUGINS_YASHA"`
}

func init() {
	viper.SetDefault("logLevel", "info")
	viper.SetDefault("postsDSN", "file:///etc/rklotz/posts")
	viper.SetDefault("postsPerPage", 10)
	viper.SetDefault("storageDSN", "boltdb:///tmp/rklotz.db")

	viper.SetDefault("web.port", 8080)
	viper.SetDefault("web.staticPath", "/etc/rklotz/static")
	viper.SetDefault("web.templatesPath", "/etc/rklotz/templates")

	viper.SetDefault("ui.theme", "foundation")
	viper.SetDefault("ui.author", "Vladimir Garvardt")
	viper.SetDefault("ui.email", "vgarvardt@gmail.com")
	viper.SetDefault("ui.description", "rKlotz - simple golang-driven blog engine")
	viper.SetDefault("ui.language", "en")
	viper.SetDefault("ui.title", "rKlotz")
	viper.SetDefault("ui.heading", "rKlotz")
	viper.SetDefault("ui.intro", "simple golang-driven blog engine")
	viper.SetDefault("ui.dateFormat", "2 Jan 2006")
	viper.SetDefault("ui.aboutPath", "/etc/rklotz/about.html")

	viper.SetDefault("rootURL.scheme", "http")
	viper.SetDefault("rootURL.host", "")
	viper.SetDefault("rootURL.path", "/")

	viper.SetDefault("plugins.enabled", []string{})
	for name, p := range plugin.GetAll() {
		viper.SetDefault(fmt.Sprintf("plugins.settings.%s", name), p.Defaults())
	}
}

func Load() (*AppConfig, error) {
	var appConfig AppConfig

	if err := viper.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	err := envconfig.Process("", &appConfig)
	if err != nil {
		return nil, err
	}

	return &appConfig, nil
}
