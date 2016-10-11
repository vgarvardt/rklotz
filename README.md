# rKlotz

## Simple golang-driven single-user blog engine on top of [Bolt DB](https://github.com/boltdb/bolt)

### Requirements

As of version 0.3 rKlotz requires Go >= 1.6

##### Install and run in dev env with automatic server reload on files change

```sh
$ git clone git@github.com:vgarvardt/rklotz.git
$ cd rklotz
$ go get github.com/constabulary/gb/...
$ brew install fswatch npm
$ npm install -g bower
$ make vendor
$ make serve
```

Then open `http://127.0.0.1:8080` in your browser.
Admin area available at `http://127.0.0.1:8080/@`, login and password are both set to `q` by default,
so don't forget to override and change it on production.

### Env and config values overriding

`./config.ini` is the base config file loaded every time when rKlotz server is started.
Env can be used to override its values. E.g. for `dev` env create file `./dev.ini`, put
required parameters there and then start server with `env` parameter set to `dev`:

```sh
$ make build && ./bin/rklotz --env dev
```

### Plugins

rKlots supports plugins. Currently the following are implemented:

* [Disqus](https://disqus.com/) (`disqus`) - posts comments
* [Google Analytics](http://www.google.com/analytics/) (`ga`) - site visits analytics from Google
* [Yandex Metrika](https://metrika.yandex.ru/) (`yamka`) - site visits analytics from Yandex
* [highlight.js](https://highlightjs.org/) (`highlightjs`) - posts code highlighting
* [Yandex Share](https://tech.yandex.ru/share/) (`yasha`) - share post buttons from Yandex

To enable some of them override `plugins` option. E.g., to enable comments, code highlighting
and share buttons for your blog put the following lines into your env config:

```ini
plugins=disqus highlightjs yasha
plugin.disqus.shortname=<shortname>
```

Do not override plugin values like this `plugin.<plugin>._=...` - it lists all available plugin options
and used for internal plugin routines.

### About panel

About (author) panel can be overridden with `ui.about` option. File must have the following structure:

```html
{{ define "partial/about.html" }}
    Content goes here. html/template is used for rendering.
{{ end }}
```

## TODO

- [x] Config loader as interface in DI container
- [ ] Dockerize development (move from gb to glide)
- [ ] Dockerize deployment
- [ ] Get config values from os env 
- [ ] Implement Material Design Lite theme
- [x] Write some tests
- [ ] Cover reindex logic with tests
- [ ] Migrate to another Web Framework (maybe echo)
- [ ] Get version from VERSION file (gb does not seem to inject ldflag into packages other than main)
- [ ] Replace raw validation with govalidator - https://github.com/asaskevich/govalidator
- [ ] Post attachments (at least images) support
- [ ] Paths history with permanent redirects from old paths to new
- [ ] SemVer versioning
