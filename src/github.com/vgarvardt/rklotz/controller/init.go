package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leekchan/gtf"
	log "github.com/Sirupsen/logrus"

	"github.com/vgarvardt/rklotz/app"
	"github.com/vgarvardt/rklotz/svc"
)

var templates map[string]*template.Template

func render(c *gin.Context, templateName string, ctx gin.H) {
	config := svc.Container.MustGet(svc.DI_CONFIG).(svc.Config)

	ui := []string{"lang", "title", "heading", "intro", "theme", "author", "description", "date_format"}
	for i := 0; i < len(ui); i++ {
		if _, ok := ctx[ui[i]]; !ok {
			ctx[ui[i]] = config.String("ui." + ui[i])
		}
	}

	ctx["instance_id"] = app.InstanceId()
	ctx["url_path"] = c.Request.URL.Path

	// enabled plugins and their settings
	enabledPlugins := strings.Split(config.String("plugins"), " ")
	plugin := make(map[string]map[string]template.JS)
	plugins := make(map[string]bool)
	for i := 0; i < len(enabledPlugins); i++ {
		plugin[enabledPlugins[i]] = make(map[string]template.JS)
		plugins[enabledPlugins[i]] = true

		pluginParams := strings.Split(config.String(fmt.Sprintf("plugin.%s._", enabledPlugins[i])), " ")
		for j := 0; j < len(pluginParams); j++ {
			pluginCfgKey := fmt.Sprintf("plugin.%s.%s", enabledPlugins[i], pluginParams[j])
			plugin[enabledPlugins[i]][pluginParams[j]] = template.JS(config.String(pluginCfgKey))
		}
	}
	ctx["plugins"] = plugins
	ctx["plugin"] = plugin

	c.Writer.WriteHeader(http.StatusOK)
	if err := templates[templateName].Execute(c.Writer, ctx); err != nil {
		panic(err)
	}
}

func redirect(c *gin.Context, url string) {
	http.Redirect(c.Writer, c.Request, url, http.StatusFound)
}

func getTmplFuncMap() template.FuncMap {
	config := svc.Container.MustGet(svc.DI_CONFIG).(svc.Config)

	funcs := gtf.GtfFuncMap

	funcs["format_date"] = func(value time.Time) string {
		return value.Format(config.String("ui.date_format"))
	}
	funcs["add"] = func(arg int, value int) int {
		return value + arg
	}
	funcs["safe"] = func(value string) template.HTML {
		return template.HTML(value)
	}
	funcs["date"] = func(format string, value time.Time) string {
		return value.Format(format)
	}

	return funcs
}

func init() {
	config := svc.Container.MustGet(svc.DI_CONFIG).(svc.Config)

	partials := []string{
		fmt.Sprintf("%s/templates/%s/partial/alert.html", app.RootDir(), config.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/heading.html", app.RootDir(), config.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/info.html", app.RootDir(), config.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/nav.html", app.RootDir(), config.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/posts.html", app.RootDir(), config.String("ui.theme")),

		fmt.Sprintf("%s/templates/plugins/disqus.html", app.RootDir()),
		fmt.Sprintf("%s/templates/plugins/ga.html", app.RootDir()),
		fmt.Sprintf("%s/templates/plugins/highlightjs.html", app.RootDir()),
		fmt.Sprintf("%s/templates/plugins/yamka.html", app.RootDir()),
		fmt.Sprintf("%s/templates/plugins/yasha.html", app.RootDir()),
	}

	logger := svc.Container.MustGet(svc.DI_LOGGER).(*log.Logger)

	uiAbout := strings.TrimSpace(config.String("ui.about"))
	if len(uiAbout) < 1 {
		logger.Info("Loading default theme about panel")
		uiAbout = fmt.Sprintf("%s/templates/%s/partial/about.html", app.RootDir(), config.String("ui.theme"))
	} else {
		logger.WithField("path", uiAbout).Info("Loading custom about panel")
	}

	partials = append(partials, uiAbout)
	baseFiles := append(partials, fmt.Sprintf("%s/templates/%s/base.html", app.RootDir(), config.String("ui.theme")))

	baseTemplate := template.Must(template.New("base.html").Funcs(getTmplFuncMap()).ParseFiles(baseFiles...))

	templates = make(map[string]*template.Template)
	for _, tmplName := range []string{
		"index.html", "post.html", "tag.html",
		"@/drafts.html", "@/published.html", "@/form.html",
	} {
		templates[tmplName] = template.Must(template.Must(baseTemplate.Clone()).ParseFiles(
			fmt.Sprintf("%s/templates/%s/%s", app.RootDir(), config.String("ui.theme"), tmplName),
		))
	}
}
