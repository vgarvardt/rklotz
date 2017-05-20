# rKlotz

[![Build Status](https://travis-ci.org/vgarvardt/rklotz.svg?branch=master)](https://travis-ci.org/vgarvardt/rklotz)
[![Coverage Status](https://coveralls.io/repos/github/vgarvardt/rklotz/badge.svg?branch=master)](https://coveralls.io/github/vgarvardt/rklotz?branch=master)

## Simple golang-driven single-user blog engine on top of [Bolt DB](https://github.com/boltdb/bolt)

### Install and run in dev env with automatic server reload on files change

```sh
$ git clone git@github.com:vgarvardt/rklotz.git
$ cd rklotz
$ make deps
$ make build
$ ./dist/rklotz
```

Then open `http://127.0.0.1:8080` in your browser.
Admin area available at `http://127.0.0.1:8080/@`, login and password are both set to `q` by default,
so don't forget to override and change it on production.

### Base config and values overriding

`./db/config.ini` is the base config file loaded every time when rKlotz server is started.
Environment variables can be used to override its values. To override default config value
create env value with name `RKLOTZ_<config_section>`, e.g. `RKLOTZ_ui.title=My Blog`.

To override config value for development environment set values in `./env.dev.txt` - file
is loaded as env file to docker.

### Plugins

rKlots supports plugins. Currently the following are implemented:

* [Disqus](https://disqus.com/) (`disqus`) - posts comments
* [Google Analytics](http://www.google.com/analytics/) (`ga`) - site visits analytics from Google
* [Yandex Metrika](https://metrika.yandex.ru/) (`yamka`) - site visits analytics from Yandex
* [highlight.js](https://highlightjs.org/) (`highlightjs`) - posts code highlighting
* [Yandex Share](https://tech.yandex.ru/share/) (`yasha`) - share post buttons from Yandex

To enable some of them override `plugins` option. E.g., to enable comments, code highlighting
and share buttons for your blog set the following env variables:

```ini
RKLOTZ_plugins=disqus highlightjs yasha
RKLOTZ_plugin.disqus.shortname=<shortname>
```

Do not override plugin values like this `plugin.<plugin>._=...` - it lists all available plugin options
and used for internal plugin routines.

### About panel

About (author) panel can be overridden with `./var/about.html` template file.
Template must have the following structure:

```html
{{ define "partial/about.html" }}
    Content goes here. html/template is used for rendering.
{{ end }}
```

## TODO

- [ ] Dockerize deployment
- [x] Get config values from os env
- [ ] Implement Material Design Lite theme
- [x] Write some tests
- [x] Cover reindex logic with tests
- [x] Migrate to another Web Framework (maybe echo)
- [x] Get version from VERSION file (gb does not seem to inject ldflag into packages other than main)
- [ ] Post attachments (at least images) support
- [ ] Paths history with permanent redirects from old paths to new
- [x] SemVer versioning
