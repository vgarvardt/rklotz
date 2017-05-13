package config

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	appConfig, err := Load()
	assert.NoError(t, err)

	assert.Equal(t, "info", appConfig.LogLevel)
	assert.Equal(t, "file:///etc/rklotz/posts", appConfig.PostsDSN)
	assert.Equal(t, 10, appConfig.PostsPerPage)
	assert.Equal(t, "boltdb:///tmp/rklotz.db", appConfig.StorageDSN)

	assert.Equal(t, 8080, appConfig.Web.Port)
	assert.Equal(t, "/etc/rklotz/static", appConfig.Web.StaticPath)
	assert.Equal(t, "/etc/rklotz/templates", appConfig.Web.TemplatesPath)

	assert.Equal(t, "foundation", appConfig.UI.Theme)
	assert.Equal(t, "Vladimir Garvardt", appConfig.UI.Author)
	assert.Equal(t, "vgarvardt@gmail.com", appConfig.UI.Email)
	assert.Equal(t, "rKlotz - simple golang-driven blog engine", appConfig.UI.Description)
	assert.Equal(t, "en", appConfig.UI.Language)
	assert.Equal(t, "rKlotz", appConfig.UI.Title)
	assert.Equal(t, "rKlotz", appConfig.UI.Heading)
	assert.Equal(t, "simple golang-driven blog engine", appConfig.UI.Intro)
	assert.Equal(t, "2 Jan 2006", appConfig.UI.DateFormat)
	assert.Equal(t, "/etc/rklotz/about.html", appConfig.UI.AboutPath)

	assert.Equal(t, "http", appConfig.RootURL.Scheme)
	assert.Equal(t, "", appConfig.RootURL.Host)
	assert.Equal(t, "/", appConfig.RootURL.Path)

	assert.Equal(t, []string{}, appConfig.Plugins.Enabled)
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
	os.Setenv("UI_ABOUT_PATH", "/path/to/about.html")

	os.Setenv("ROOT_URL_SCHEME", "gopher")
	os.Setenv("ROOT_URL_HOST", "example.com")
	os.Setenv("ROOT_URL_PATH", "/blog")

	os.Setenv("PLUGINS_ENABLED", "foo,bar,baz")

	appConfig, err := Load()
	assert.NoError(t, err)

	assert.Equal(t, "debug", appConfig.LogLevel)
	assert.Equal(t, "file:///path/to/posts", appConfig.PostsDSN)
	assert.Equal(t, 42, appConfig.PostsPerPage)
	assert.Equal(t, "mysql://root@localhost/rklotz", appConfig.StorageDSN)

	assert.Equal(t, 8081, appConfig.Web.Port)
	assert.Equal(t, "/path/to/static", appConfig.Web.StaticPath)
	assert.Equal(t, "/path/to/templates", appConfig.Web.TemplatesPath)

	assert.Equal(t, "premium", appConfig.UI.Theme)
	assert.Equal(t, "Neal Stephenson", appConfig.UI.Author)
	assert.Equal(t, "neal@stephenson.com", appConfig.UI.Email)
	assert.Equal(t, "Novel by Neal Stephenson", appConfig.UI.Description)
	assert.Equal(t, "qwghlm", appConfig.UI.Language)
	assert.Equal(t, "Cryptonomicon", appConfig.UI.Title)
	assert.Equal(t, "Anathem", appConfig.UI.Heading)
	assert.Equal(t, "Reamde", appConfig.UI.Intro)
	assert.Equal(t, "Mon Jan 2 15:04:05 -0700 MST 2006", appConfig.UI.DateFormat)
	assert.Equal(t, "/path/to/about.html", appConfig.UI.AboutPath)

	assert.Equal(t, "gopher", appConfig.RootURL.Scheme)
	assert.Equal(t, "example.com", appConfig.RootURL.Host)
	assert.Equal(t, "/blog", appConfig.RootURL.Path)

	assert.Equal(t, []string{"foo", "bar", "baz"}, appConfig.Plugins.Enabled)
}

func TestRootURL_URL(t *testing.T) {
	os.Unsetenv("ROOT_URL_SCHEME")
	os.Unsetenv("ROOT_URL_HOST")
	os.Unsetenv("ROOT_URL_PATH")

	appConfig, err := Load()
	assert.NoError(t, err)

	r := &http.Request{Host: "example.com"}

	assert.Equal(t, &url.URL{Scheme: "http", Host: "example.com", Path: "/"}, appConfig.RootURL.URL(r))

	appConfig.RootURL.Scheme = "https"
	appConfig.RootURL.Host = "protected.com"
	appConfig.RootURL.Path = "/blog"

	assert.Equal(t, &url.URL{Scheme: "https", Host: "protected.com", Path: "/blog"}, appConfig.RootURL.URL(r))
}
