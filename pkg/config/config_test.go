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

	assert.Equal(t, appConfig.LogLevel, "info")
	assert.Equal(t, appConfig.PostsDSN, "file:///etc/rklotz/posts")
	assert.Equal(t, appConfig.PostsPerPage, 10)
	assert.Equal(t, appConfig.StorageDSN, "boltdb:///tmp/rklotz.db")

	assert.Equal(t, appConfig.Web.Port, 8080)
	assert.Equal(t, appConfig.Web.StaticPath, "/etc/rklotz/static")
	assert.Equal(t, appConfig.Web.TemplatesPath, "/etc/rklotz/templates")

	assert.Equal(t, appConfig.UI.Theme, "foundation")
	assert.Equal(t, appConfig.UI.Author, "Vladimir Garvardt")
	assert.Equal(t, appConfig.UI.Email, "vgarvardt@gmail.com")
	assert.Equal(t, appConfig.UI.Description, "rKlotz - simple golang-driven blog engine")
	assert.Equal(t, appConfig.UI.Language, "en")
	assert.Equal(t, appConfig.UI.Title, "rKlotz")
	assert.Equal(t, appConfig.UI.Heading, "rKlotz")
	assert.Equal(t, appConfig.UI.Intro, "simple golang-driven blog engine")
	assert.Equal(t, appConfig.UI.DateFormat, "2 Jan 2006")
	assert.Equal(t, appConfig.UI.AboutPath, "/etc/rklotz/about.html")

	assert.Equal(t, appConfig.RootURL.Scheme, "http")
	assert.Equal(t, appConfig.RootURL.Host, "")
	assert.Equal(t, appConfig.RootURL.Path, "/")
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

	appConfig, err := Load()
	assert.NoError(t, err)

	assert.Equal(t, appConfig.LogLevel, "debug")
	assert.Equal(t, appConfig.PostsDSN, "file:///path/to/posts")
	assert.Equal(t, appConfig.PostsPerPage, 42)
	assert.Equal(t, appConfig.StorageDSN, "mysql://root@localhost/rklotz")

	assert.Equal(t, appConfig.Web.Port, 8081)
	assert.Equal(t, appConfig.Web.StaticPath, "/path/to/static")
	assert.Equal(t, appConfig.Web.TemplatesPath, "/path/to/templates")

	assert.Equal(t, appConfig.UI.Theme, "premium")
	assert.Equal(t, appConfig.UI.Author, "Neal Stephenson")
	assert.Equal(t, appConfig.UI.Email, "neal@stephenson.com")
	assert.Equal(t, appConfig.UI.Description, "Novel by Neal Stephenson")
	assert.Equal(t, appConfig.UI.Language, "qwghlm")
	assert.Equal(t, appConfig.UI.Title, "Cryptonomicon")
	assert.Equal(t, appConfig.UI.Heading, "Anathem")
	assert.Equal(t, appConfig.UI.Intro, "Reamde")
	assert.Equal(t, appConfig.UI.DateFormat, "Mon Jan 2 15:04:05 -0700 MST 2006")
	assert.Equal(t, appConfig.UI.AboutPath, "/path/to/about.html")

	assert.Equal(t, appConfig.RootURL.Scheme, "gopher")
	assert.Equal(t, appConfig.RootURL.Host, "example.com")
	assert.Equal(t, appConfig.RootURL.Path, "/blog")
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
