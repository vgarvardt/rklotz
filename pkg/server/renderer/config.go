package renderer

import (
	"net/http"
	"net/url"
)

// UIConfig is the configuration for user interface
type UIConfig struct {
	Theme       string `env:"UI_THEME,default=foundation6"`
	Author      string `env:"UI_AUTHOR,default=Vladimir Garvardt"`
	Email       string `env:"UI_EMAIL,default=vgarvardt@gmail.com"`
	Description string `env:"UI_DESCRIPTION,default=rKlotz - simple golang-driven blog engine"`
	Language    string `env:"UI_LANGUAGE,default=en"`
	Title       string `env:"UI_TITLE,default=rKlotz"`
	Heading     string `env:"UI_HEADING,default=rKlotz"`
	Intro       string `env:"UI_INTRO,default=simple golang-driven blog engine"`
	// DateFormat is format for posts, see http://golang.org/pkg/time/#Time.Format
	DateFormat string `env:"UI_DATEFORMAT,default=2 Jan 2006"`
	AboutPath  string `env:"UI_ABOUT_PATH,default=/etc/rklotz/about.tpl"`
}

// RootURLConfig is the configuration for app root url
type RootURLConfig struct {
	Scheme string `env:"ROOT_URL_SCHEME,default=http"`
	Host   string `env:"ROOT_URL_HOST"`
	Path   string `env:"ROOT_URL_PATH,default=/"`
}

// URL returns the URL for currently configured root url
func (u RootURLConfig) URL(r *http.Request) *url.URL {
	host := u.Host
	if len(host) < 1 {
		host = r.Host
	}
	return &url.URL{Scheme: u.Scheme, Host: host, Path: u.Path}
}
