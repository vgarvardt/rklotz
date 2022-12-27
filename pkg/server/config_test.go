package server

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_DefaultValues(t *testing.T) {
	cfg, err := LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, "info", cfg.Level)
	assert.Equal(t, "file:///etc/rklotz/posts", cfg.PostsDSN)
	assert.Equal(t, 10, cfg.PostsPerPage)
	assert.Equal(t, "boltdb:///tmp/rklotz.db", cfg.StorageDSN)

	assert.Equal(t, 8080, cfg.HTTPConfig.Port)
	assert.Equal(t, "/etc/rklotz/static", cfg.HTTPConfig.StaticPath)
	assert.Equal(t, "/etc/rklotz/templates", cfg.HTTPConfig.TemplatesPath)

	assert.Equal(t, false, cfg.SSLConfig.Enabled)
	assert.Equal(t, 8443, cfg.SSLConfig.Port)
	assert.Equal(t, "/tmp", cfg.SSLConfig.CacheDir)

	assert.Equal(t, "foundation6", cfg.UIConfig.Theme)
	assert.Equal(t, "Vladimir Garvardt", cfg.UIConfig.Author)
	assert.Equal(t, "vgarvardt@gmail.com", cfg.UIConfig.Email)
	assert.Equal(t, "rKlotz - simple golang-driven blog engine", cfg.UIConfig.Description)
	assert.Equal(t, "en", cfg.UIConfig.Language)
	assert.Equal(t, "rKlotz", cfg.UIConfig.Title)
	assert.Equal(t, "rKlotz", cfg.UIConfig.Heading)
	assert.Equal(t, "simple golang-driven blog engine", cfg.UIConfig.Intro)
	assert.Equal(t, "2 Jan 2006", cfg.UIConfig.DateFormat)
	assert.Equal(t, "/etc/rklotz/about.tpl", cfg.UIConfig.AboutPath)

	assert.Equal(t, "http", cfg.RootURLConfig.Scheme)
	assert.Equal(t, "", cfg.RootURLConfig.Host)
	assert.Equal(t, "/", cfg.RootURLConfig.Path)

	assert.Len(t, cfg.Config.Enabled, 0)
}

func TestLoad(t *testing.T) {
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("POSTS_DSN", "file:///path/to/posts")
	os.Setenv("POSTS_PERPAGE", "42")
	os.Setenv("STORAGE_DSN", "mysql://root@localhost/rklotz")

	os.Setenv("WEB_PORT", "8081")
	os.Setenv("WEB_STATIC_PATH", "/path/to/static")
	os.Setenv("WEB_TEMPLATES_PATH", "/path/to/templates")

	os.Setenv("UI_THEME", "premium")
	os.Setenv("UI_AUTHOR", "Neal Stephenson")
	os.Setenv("UI_EMAIL", "neal@stephenson.com")
	os.Setenv("UI_DESCRIPTION", "Novel by Neal Stephenson")
	os.Setenv("UI_LANGUAGE", "qwghlm")
	os.Setenv("UI_TITLE", "Cryptonomicon")
	os.Setenv("UI_HEADING", "Anathem")
	os.Setenv("UI_INTRO", "Reamde")
	os.Setenv("UI_DATEFORMAT", "Mon Jan 2 15:04:05 -0700 MST 2006")
	os.Setenv("UI_ABOUT_PATH", "/path/to/about.tpl")

	os.Setenv("ROOT_URL_SCHEME", "gopher")
	os.Setenv("ROOT_URL_HOST", "example.com")
	os.Setenv("ROOT_URL_PATH", "/blog")

	os.Setenv("PLUGINS_ENABLED", "foo,bar,baz")

	os.Setenv("PLUGINS_DISQUS", "shortname:foo")
	os.Setenv("PLUGINS_GA", "tracking_id:bar")
	os.Setenv("PLUGINS_YAMKA", "id:baz")
	os.Setenv("PLUGINS_HIGHLIGHTJS", "theme:foo,version:9.9.9")
	os.Setenv("PLUGINS_YASHA", "services:facebook twitter,l10n:de")

	appConfig, err := LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, "debug", appConfig.Level)
	assert.Equal(t, "file:///path/to/posts", appConfig.PostsDSN)
	assert.Equal(t, 42, appConfig.PostsPerPage)
	assert.Equal(t, "mysql://root@localhost/rklotz", appConfig.StorageDSN)

	assert.Equal(t, 8081, appConfig.HTTPConfig.Port)
	assert.Equal(t, "/path/to/static", appConfig.HTTPConfig.StaticPath)
	assert.Equal(t, "/path/to/templates", appConfig.HTTPConfig.TemplatesPath)

	assert.Equal(t, "premium", appConfig.UIConfig.Theme)
	assert.Equal(t, "Neal Stephenson", appConfig.UIConfig.Author)
	assert.Equal(t, "neal@stephenson.com", appConfig.UIConfig.Email)
	assert.Equal(t, "Novel by Neal Stephenson", appConfig.UIConfig.Description)
	assert.Equal(t, "qwghlm", appConfig.UIConfig.Language)
	assert.Equal(t, "Cryptonomicon", appConfig.UIConfig.Title)
	assert.Equal(t, "Anathem", appConfig.UIConfig.Heading)
	assert.Equal(t, "Reamde", appConfig.UIConfig.Intro)
	assert.Equal(t, "Mon Jan 2 15:04:05 -0700 MST 2006", appConfig.UIConfig.DateFormat)
	assert.Equal(t, "/path/to/about.tpl", appConfig.UIConfig.AboutPath)

	assert.Equal(t, "gopher", appConfig.RootURLConfig.Scheme)
	assert.Equal(t, "example.com", appConfig.RootURLConfig.Host)
	assert.Equal(t, "/blog", appConfig.RootURLConfig.Path)

	assert.Equal(t, []string{"foo", "bar", "baz"}, appConfig.Config.Enabled)

	assert.Equal(t, map[string]string{"shortname": "foo"}, appConfig.Config.Settings.Disqus)
	assert.Equal(t, map[string]string{"tracking_id": "bar"}, appConfig.Config.Settings.Ga)
	assert.Equal(t, map[string]string{"id": "baz"}, appConfig.Config.Settings.Yamka)
	assert.Equal(t, map[string]string{"theme": "foo", "version": "9.9.9"}, appConfig.Config.Settings.Highlightjs)
	assert.Equal(t, map[string]string{"services": "facebook twitter", "l10n": "de"}, appConfig.Config.Settings.Yasha)
}

func TestRootURL_URL(t *testing.T) {
	os.Unsetenv("ROOT_URL_SCHEME")
	os.Unsetenv("ROOT_URL_HOST")
	os.Unsetenv("ROOT_URL_PATH")

	cfg, err := LoadConfig()
	require.NoError(t, err)

	r := &http.Request{Host: "example.com"}

	assert.Equal(t, &url.URL{Scheme: "http", Host: "example.com", Path: "/"}, cfg.RootURLConfig.URL(r))

	cfg.RootURLConfig.Scheme = "https"
	cfg.RootURLConfig.Host = "protected.com"
	cfg.RootURLConfig.Path = "/blog"

	assert.Equal(t, &url.URL{Scheme: "https", Host: "protected.com", Path: "/blog"}, cfg.RootURLConfig.URL(r))
}
