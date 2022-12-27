package renderer

import (
	"net/http"
	"net/url"
)

// UIConfig is the configuration for user interface
type UIConfig struct {
	Theme       string `envconfig:"UI_THEME" default:"foundation6"`
	Author      string `envconfig:"UI_AUTHOR" default:"Vladimir Garvardt"`
	Email       string `envconfig:"UI_EMAIL" default:"vgarvardt@gmail.com"`
	Description string `envconfig:"UI_DESCRIPTION" default:"rKlotz - simple golang-driven blog engine"`
	Language    string `envconfig:"UI_LANGUAGE" default:"en"`
	Title       string `envconfig:"UI_TITLE" default:"rKlotz"`
	Heading     string `envconfig:"UI_HEADING" default:"rKlotz"`
	Intro       string `envconfig:"UI_INTRO" default:"simple golang-driven blog engine"`
	// DateFormat is format for posts, see http://golang.org/pkg/time/#Time.Format
	DateFormat string `envconfig:"UI_DATEFORMAT" default:"2 Jan 2006"`
	AboutPath  string `envconfig:"UI_ABOUT_PATH" default:"/etc/rklotz/about.tpl"`
}

// RootURLConfig is the configuration for app root url
type RootURLConfig struct {
	Scheme string `envconfig:"ROOT_URL_SCHEME" default:"http"`
	Host   string `envconfig:"ROOT_URL_HOST"`
	Path   string `envconfig:"ROOT_URL_PATH" default:"/"`
}

// URL returns the URL for currently configured root url
func (u RootURLConfig) URL(r *http.Request) *url.URL {
	host := u.Host
	if len(host) < 1 {
		host = r.Host
	}
	return &url.URL{Scheme: u.Scheme, Host: host, Path: u.Path}
}
