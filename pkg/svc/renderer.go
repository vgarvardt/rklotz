package svc

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/leekchan/gtf"
)

type renderable struct {
	templates  map[string]*template.Template
	instanceId string
}

func (r *renderable) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	config := Container.MustGet(DI_CONFIG).(Config)

	templateData := data.(map[string]interface{})

	ui := []string{"lang", "title", "heading", "intro", "theme", "author", "description", "date_format"}
	for i := 0; i < len(ui); i++ {
		if _, ok := templateData[ui[i]]; !ok {
			templateData[ui[i]] = config.String("ui." + ui[i])
		}
	}

	templateData["instance_id"] = r.instanceId
	templateData["url_path"] = ctx.Request().URI()

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
	templateData["plugins"] = plugins
	templateData["plugin"] = plugin

	return r.templates[name].Execute(w, templateData)
}

func Renderer(rootDir, instanceId string) *renderable {
	config := Container.MustGet(DI_CONFIG).(Config)
	logger := Container.MustGet(DI_LOGGER).(*log.Logger)

	partials := []string{
		fmt.Sprintf("%s/templates/%s/partial/alert.html", rootDir, config.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/heading.html", rootDir, config.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/info.html", rootDir, config.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/nav.html", rootDir, config.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/posts.html", rootDir, config.String("ui.theme")),

		fmt.Sprintf("%s/templates/plugins/disqus.html", rootDir),
		fmt.Sprintf("%s/templates/plugins/ga.html", rootDir),
		fmt.Sprintf("%s/templates/plugins/highlightjs-css.html", rootDir),
		fmt.Sprintf("%s/templates/plugins/highlightjs-js.html", rootDir),
		fmt.Sprintf("%s/templates/plugins/yamka.html", rootDir),
		fmt.Sprintf("%s/templates/plugins/yasha.html", rootDir),
	}

	uiAbout := fmt.Sprintf("%s/var/about.html", rootDir)
	if _, err := os.Stat(uiAbout); os.IsNotExist(err) {
		logger.Info("Loading default theme about panel")
		uiAbout = fmt.Sprintf("%s/templates/%s/partial/about.html", rootDir, config.String("ui.theme"))
	} else {
		logger.WithField("path", uiAbout).Info("Loading custom about panel")
	}

	partials = append(partials, uiAbout)
	baseFiles := append(partials, fmt.Sprintf("%s/templates/%s/base.html", rootDir, config.String("ui.theme")))

	baseTemplate := template.Must(template.New("base.html").Funcs(getTmplFuncMap()).ParseFiles(baseFiles...))

	renderer := &renderable{
		templates:  make(map[string]*template.Template),
		instanceId: instanceId,
	}

	for _, tmplName := range []string{
		"index.html", "post.html", "tag.html",
		"@/drafts.html", "@/published.html", "@/form.html",
	} {
		tmplPath := fmt.Sprintf("%s/templates/%s/%s", rootDir, config.String("ui.theme"), tmplName)

		logger.WithFields(
			log.Fields{"name": tmplName, "path": tmplPath},
		).Debug("Initializing template")
		renderer.templates[tmplName] = template.Must(template.Must(baseTemplate.Clone()).ParseFiles(tmplPath))
	}

	return renderer
}

func getTmplFuncMap() template.FuncMap {
	config := Container.MustGet(DI_CONFIG).(Config)

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
