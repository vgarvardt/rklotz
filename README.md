# rKlotz

[![Coverage Status](https://codecov.io/gh/vgarvardt/rklotz/branch/master/graph/badge.svg)](https://codecov.io/gh/vgarvardt/rklotz)

> Yet another simple single-user file-based golang-driven blog engine

## Run locally

You need to have [Docker](https://www.docker.com/) installed and running.

```bash
docker run --rm -it -p 8080:8080 vgarvardt/rklotz
```

Then open `http://127.0.0.1:8080` in your browser.

## Build and run locally

You need to have Go 1.7+ and [Docker](https://www.docker.com/) installed and running.

```bash
$ git clone git@github.com:vgarvardt/rklotz.git
$ cd rklotz
$ make deps
$ make build
$ docker run -it -p 8080:8080 vgarvardt/rklotz
```

Then open `http://127.0.0.1:8080` in your browser.

## Build your own blog based on rKlotz

See [github.com/vgarvardt/itskrig.com](https://github.com/vgarvardt/itskrig.com) for example
on how to build your blog using `rKlotz` as base image.

## Posts

Posts in rKlotz are just files written in some markup language.
Currently only [Markdown](https://daringfireball.net/projects/markdown/syntax) (`md` extension) is supported.

Post file has the following structure:

* *Line 1*: Post title
* *Line 2*: Post publishing date - posts are ordered by publishing date in reverse chronological order.
  Date must be in [`RFC822Z` format](https://golang.org/pkg/time/#pkg-constants)
* *Line 3*: Post tags - comma-separated tags list
* *Line 4*: Reserved for further usage
* *Line 5*: Post delimiter - `+++` for Markdown, not necessary that line number, may be preceded by any number
  of lines before delimiter
* *Line 6*: Post body - may be preceded by any number of lines before post body, after delimiter

Post path is determined automatically from its path, relative to posts root path (see settings).

Posts examples are available in [asserts/posts](./assets/posts).

## Settings

Currently The following settings (environment variables) are available:

### Base application settings

* `LOG_LEVEL` (default `info`) - logging level
* `POSTS_DSN` (default `file:///etc/rklotz/posts`) - posts root path in the format `storage://<path>`.
  Currently the following storage types are supported:
  * `file` - local file system
* `POSTS_PERPAGE` (default `10`) - number of posts per page
* `STORAGE_DSN` (default `boltdb:///tmp/rklotz.db`) - posts storage in run-time in the format `storage://path`.
  Currently the following storage types are supported:
  * `boltdb` - storage on top of [BoltDB](https://github.com/boltdb/bolt) and [Storm](https://github.com/asdine/storm)
  * `memory` - store all posts in memory, perfect for hundreds or even thousands of posts

### Web application settings

* `WEB_PORT` (default `8080`) - port to run the `http` server
* `WEB_STATIC_PATH` (default `/etc/rklotz/static`) - static files root path
* `WEB_TEMPLATES_PATH` (default `/etc/rklotz/templates`) - templates root path

### SSL settings

`rKlotz` supports SSL with [Let's Encrypt](https://letsencrypt.org/).

* `SSL_ENABLED` (default `false`) - enables SSL/TLS
* `SSL_PORT` (default `8443`) - SSL port
* `SSL_HOST` - host to validate for SSL
* `SSL_REDIRECT_HTTP` (default `true`) - redirect `http` requests to `https` if SSL is enabled,
  otherwise both HTTP and HTTPS will be served
* `SSL_CACHE_DIR` (default `/tmp`) - directory to cache retrieved certificate

### HTML and UI settings

* `UI_THEME` (default `foundation`) - theme name. Themes list available in [templates](./templates)
  (except for `plugins`, that are plugins templates, see bellow)
* `UI_AUTHOR` (default `Vladimir Garvardt`) - blog author name (html head meta)
* `UI_EMAIL` (default `vgarvardt@gmail.com`) - blog author email
* `UI_DESCRIPTION` (default `rKlotz - simple golang-driven blog engine`) - blog description (html head meta)
* `UI_LANGUAGE` (default `en`) - blog language (html lang)
* `UI_TITLE` (default `rKlotz`) - blog title (html title)
* `UI_HEADING` (default `rKlotz`) - blog heading (index page header)
* `UI_INTRO` (default `simple golang-driven blog engine`) - blog intro (index page header)
* `UI_DATEFORMAT` (default `2 Jan 2006`) - post publishing date display format.
  Must be compatible with [`time.Format()`](http://golang.org/pkg/time/#Time.Format). See examples in
  [predefined time formats](https://golang.org/pkg/time/#pkg-constants).
* `UI_ABOUT_PATH` (default `/etc/rklotz/about.html`) - path to custom "about panel".
  If not found - `<WEB_TEMPLATES_PATH>/<UI_THEME>/partial/about.html` is used.

#### About panel

Template must have the following structure:

```html
{{ define "partial/about.html" }}
    Content goes here. html/template is used for rendering.
{{ end }}
```

See about panel example in [default theme](./templates/foundation/partial/about.html).

### Root URL settings

* `ROOT_URL_SCHEME` (default `http`) - blog absolute url scheme. Currently `https` si not supported on rKlotz web
  application level (in plans), so use `https` only if you have `SSL/TLS` certificate termination on the level before
  rKlotz (e.g. nginx as reverse proxy before your blog).
* `ROOT_URL_HOST` (default ``) - blog absolute url host. If empty - request host is used.
* `ROOT_URL_PATH` (default `/`) - blog absolute url path prefix. In case your blog is hosted on the second (or deeper)
  path level, e.g. `http://example.com/blog` (`ROOT_URL_PATH`=`/blog`)

### Plugins settings

### Plugins

rKlots supports plugins. Currently the following are implemented:

* [Disqus](https://disqus.com/) (`disqus`) - posts comments
* [Google Analytics](http://www.google.com/analytics/) (`ga`) - site visits analytics from Google
* [Google Tag Manager](https://tagmanager.google.com) (`gtm`) - tag management analytics from Google
* [highlight.js](https://highlightjs.org/) (`highlightjs`) - posts code highlighting
* [Yandex Metrika](https://metrika.yandex.ru/) (`yamka`) - site visits analytics from Yandex
* [Yandex Share](https://tech.yandex.ru/share/) (`yasha`) - share post buttons from Yandex

Plugins configuration available with the following settings:

* `PLUGINS_ENABLED` - comma-separated plugins list, e.g. `disqus,ga,highlightjs`
  to enable *Disqus*, *Google Analytics* and *highlight.js* plugins
* `PLUGINS_DISQUS` - *Disqus* plugin configuration in the format `<config1>:<value1>,<config2>:<value2>,...`
  The following configurations are available:
  * `shortname` (required) - account short name
* `PLUGINS_GA` - *Google Analytics* plugin configuration in the format `<config1>:<value1>,<config2>:<value2>,...`
  The following configurations are available:
  * `tracking_id` (required) - analytics tracking ID
* `PLUGINS_GTM` - *Google Tag Manager* plugin configuration in the format `<config1>:<value1>,<config2>:<value2>,...`
  The following configurations are available:
  * `id` (required) - tag manager ID
* `PLUGINS_HIGHLIGHTJS` - *highlight.js* plugin configuration in the format `<config1>:<value1>,<config2>:<value2>,...`
  The following configurations are available:
  * `version` (default `9.7.0`) - library version
  * `theme` (default `idea`) - colour scheme/theme
* `PLUGINS_YAMKA` - *Yandex Metrika* plugin configuration in the format `<config1>:<value1>,<config2>:<value2>,...`
  The following configurations are available:
  * `id` (required) - metrika ID 
* `PLUGINS_YASHA` - *Yandex Share* plugin configuration in the format `<config1>:<value1>,<config2>:<value2>,...`
  The following configurations are available (see fill list of values on plugin page):
  * `services` (default: `facebook twitter gplus`) - space-separated services list
  * `size` (default: `m`) - icons size: `m` - medium, `s` - small
  * `lang` (default `en`) - widget language, see [docs page](https://tech.yandex.ru/share/doc/dg/add-docpage/)
  for complete list of available languages

## TODO

- [x] Dockerize deployment
- [x] Get config values from os env
- [ ] Implement at least one more theme
- [x] Write some tests
- [x] Cover reindex logic with tests
- [x] Migrate to another Web Framework (maybe echo)
- [x] Get version from VERSION file (gb does not seem to inject ldflag into packages other than main)
- [x] SemVer versioning
- [x] SSL/TLS with Let's Encrypt
- [ ] Implement [`badger`](https://github.com/dgraph-io/badger) storage
- [x] Implement `memory` storage
- [ ] Implement `git` loader
