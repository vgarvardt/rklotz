package controller

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leekchan/gtf"

	"github.com/vgarvardt/rklotz/cfg"
)

var templates map[string]*template.Template

func render(c *gin.Context, templateName string, ctx gin.H) {
	ui := []string{"lang", "title", "heading", "intro", "theme", "author", "description", "date_format"}
	for i := 0; i < len(ui); i++ {
		if _, ok := ctx[ui[i]]; !ok {
			ctx[ui[i]] = cfg.String("ui." + ui[i])
		}
	}

	ctx["instance_id"] = cfg.GetInstanceId()
	ctx["url_path"] = c.Request.URL.Path

	// enabled plugins and their settings
	enabledPlugins := strings.Split(cfg.String("plugins"), " ")
	plugin := make(map[string]map[string]template.JS)
	plugins := make(map[string]bool)
	for i := 0; i < len(enabledPlugins); i++ {
		plugin[enabledPlugins[i]] = make(map[string]template.JS)
		plugins[enabledPlugins[i]] = true

		pluginParams := strings.Split(cfg.String(fmt.Sprintf("plugin.%s._", enabledPlugins[i])), " ")
		for j := 0; j < len(pluginParams); j++ {
			pluginCfgKey := fmt.Sprintf("plugin.%s.%s", enabledPlugins[i], pluginParams[j])
			plugin[enabledPlugins[i]][pluginParams[j]] = template.JS(cfg.String(pluginCfgKey))
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
	funcs := gtf.GtfFuncMap

	funcs["format_date"] = func(value time.Time) string {
		return value.Format(cfg.String("ui.date_format"))
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
	partials := []string{
		fmt.Sprintf("%s/templates/%s/partial/alert.html", cfg.GetRootDir(), cfg.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/heading.html", cfg.GetRootDir(), cfg.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/info.html", cfg.GetRootDir(), cfg.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/nav.html", cfg.GetRootDir(), cfg.String("ui.theme")),
		fmt.Sprintf("%s/templates/%s/partial/posts.html", cfg.GetRootDir(), cfg.String("ui.theme")),

		fmt.Sprintf("%s/templates/plugins/disqus.html", cfg.GetRootDir()),
		fmt.Sprintf("%s/templates/plugins/ga.html", cfg.GetRootDir()),
		fmt.Sprintf("%s/templates/plugins/highlightjs.html", cfg.GetRootDir()),
		fmt.Sprintf("%s/templates/plugins/yamka.html", cfg.GetRootDir()),
		fmt.Sprintf("%s/templates/plugins/yasha.html", cfg.GetRootDir()),
	}

	uiAbout := strings.TrimSpace(cfg.String("ui.about"))
	if len(uiAbout) < 1 {
		cfg.Log("Loading default theme about panel")
		uiAbout = fmt.Sprintf("%s/templates/%s/partial/about.html", cfg.GetRootDir(), cfg.String("ui.theme"))
	} else {
		cfg.Log(fmt.Sprintf("Loading custom about panel @ %s", uiAbout))
	}

	partials = append(partials, uiAbout)
	baseFiles := append(partials, fmt.Sprintf("%s/templates/%s/base.html", cfg.GetRootDir(), cfg.String("ui.theme")))

	baseTemplate := template.Must(template.New("base.html").Funcs(getTmplFuncMap()).ParseFiles(baseFiles...))

	templates = make(map[string]*template.Template)
	for _, tmplName := range []string{
		"index.html", "post.html", "tag.html",
		"@/drafts.html", "@/published.html", "@/form.html",
	} {
		templates[tmplName] = template.Must(template.Must(baseTemplate.Clone()).ParseFiles(
			fmt.Sprintf("%s/templates/%s/%s", cfg.GetRootDir(), cfg.String("ui.theme"), tmplName),
		))
	}
}
