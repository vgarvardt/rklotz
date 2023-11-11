package server

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_DefaultValues(t *testing.T) {
	ctx := context.Background()

	cfg, err := LoadConfig(ctx)
	assert.NoError(t, err)

	assert.Equal(t, slog.LevelInfo, cfg.Level)
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
	ctx := context.Background()

	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("POSTS_DSN", "file:///path/to/posts")
	t.Setenv("POSTS_PERPAGE", "42")
	t.Setenv("STORAGE_DSN", "mysql://root@localhost/rklotz")

	t.Setenv("WEB_PORT", "8081")
	t.Setenv("WEB_STATIC_PATH", "/path/to/static")
	t.Setenv("WEB_TEMPLATES_PATH", "/path/to/templates")

	t.Setenv("UI_THEME", "premium")
	t.Setenv("UI_AUTHOR", "Neal Stephenson")
	t.Setenv("UI_EMAIL", "neal@stephenson.com")
	t.Setenv("UI_DESCRIPTION", "Novel by Neal Stephenson")
	t.Setenv("UI_LANGUAGE", "qwghlm")
	t.Setenv("UI_TITLE", "Cryptonomicon")
	t.Setenv("UI_HEADING", "Anathem")
	t.Setenv("UI_INTRO", "Reamde")
	t.Setenv("UI_DATEFORMAT", "Mon Jan 2 15:04:05 -0700 MST 2006")
	t.Setenv("UI_ABOUT_PATH", "/path/to/about.tpl")

	t.Setenv("ROOT_URL_SCHEME", "gopher")
	t.Setenv("ROOT_URL_HOST", "example.com")
	t.Setenv("ROOT_URL_PATH", "/blog")

	t.Setenv("PLUGINS_ENABLED", "foo,bar,baz")

	t.Setenv("PLUGINS_DISQUS", "shortname:foo")
	t.Setenv("PLUGINS_GA", "tracking_id:bar")
	t.Setenv("PLUGINS_YAMKA", "id:baz")
	t.Setenv("PLUGINS_HIGHLIGHTJS", "theme:foo,version:9.9.9")
	t.Setenv("PLUGINS_YASHA", "services:facebook twitter,l10n:de")

	appConfig, err := LoadConfig(ctx)
	assert.NoError(t, err)

	assert.Equal(t, slog.LevelDebug, appConfig.Level)
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
	ctx := context.Background()

	err := os.Unsetenv("ROOT_URL_SCHEME")
	require.NoError(t, err)
	err = os.Unsetenv("ROOT_URL_HOST")
	require.NoError(t, err)
	err = os.Unsetenv("ROOT_URL_PATH")
	require.NoError(t, err)

	cfg, err := LoadConfig(ctx)
	require.NoError(t, err)

	r := &http.Request{Host: "example.com"}

	assert.Equal(t, &url.URL{Scheme: "http", Host: "example.com", Path: "/"}, cfg.RootURLConfig.URL(r))

	cfg.RootURLConfig.Scheme = "https"
	cfg.RootURLConfig.Host = "protected.com"
	cfg.RootURLConfig.Path = "/blog"

	assert.Equal(t, &url.URL{Scheme: "https", Host: "protected.com", Path: "/blog"}, cfg.RootURLConfig.URL(r))
}
