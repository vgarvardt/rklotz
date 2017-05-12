package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/leekchan/gtf"
	"github.com/vgarvardt/rklotz/pkg/config"
)

const TemplateNameDateKey = "template_name"

type htmlRenderer struct {
	templates  map[string]*template.Template
	instanceId string
	uiSettings config.UISetting
}

func NewHTMLRenderer(templatesPath string, instanceId string, uiSettings config.UISetting) *htmlRenderer {
	partials := []string{
		fmt.Sprintf("%s/%s/partial/alert.html", templatesPath, uiSettings.Theme),
		fmt.Sprintf("%s/%s/partial/heading.html", templatesPath, uiSettings.Theme),
		fmt.Sprintf("%s/%s/partial/info.html", templatesPath, uiSettings.Theme),
		fmt.Sprintf("%s/%s/partial/pagination.html", templatesPath, uiSettings.Theme),
		fmt.Sprintf("%s/%s/partial/posts.html", templatesPath, uiSettings.Theme),

		fmt.Sprintf("%s/plugins/disqus.html", templatesPath),
		fmt.Sprintf("%s/plugins/ga.html", templatesPath),
		fmt.Sprintf("%s/plugins/highlightjs-css.html", templatesPath),
		fmt.Sprintf("%s/plugins/highlightjs-js.html", templatesPath),
		fmt.Sprintf("%s/plugins/yamka.html", templatesPath),
		fmt.Sprintf("%s/plugins/yasha.html", templatesPath),
	}

	uiAbout := uiSettings.AboutPath
	if _, err := os.Stat(uiAbout); os.IsNotExist(err) {
		log.Info("Loading default theme about panel")
		uiAbout = fmt.Sprintf("%s/%s/partial/about.html", templatesPath, uiSettings.Theme)
	} else {
		log.WithField("path", uiAbout).Info("Loading custom about panel")
	}

	partials = append(partials, uiAbout)
	baseFiles := append(partials, fmt.Sprintf("%s/%s/base.html", templatesPath, uiSettings.Theme))

	baseTemplate := template.Must(
		template.New("base.html").Funcs(getTmplFuncMap(uiSettings.DateFormat)).ParseFiles(baseFiles...))

	instance := &htmlRenderer{
		templates:  make(map[string]*template.Template),
		instanceId: instanceId,
		uiSettings: uiSettings,
	}

	for _, tmplName := range []string{"index.html", "post.html", "tag.html"} {
		tmplPath := fmt.Sprintf("%s/%s/%s", templatesPath, uiSettings.Theme, tmplName)

		log.WithFields(log.Fields{"name": tmplName, "path": tmplPath}).Debug("Initializing template")
		instance.templates[tmplName] = template.Must(template.Must(baseTemplate.Clone()).ParseFiles(tmplPath))
	}

	return instance
}

func (r *htmlRenderer) Render(w http.ResponseWriter, code int, data interface{}) {
	templateData := data.(map[string]interface{})

	templateData["lang"] = r.uiSettings.Language
	templateData["title"] = r.uiSettings.Title
	templateData["heading"] = r.uiSettings.Heading
	templateData["intro"] = r.uiSettings.Intro
	templateData["theme"] = r.uiSettings.Theme
	templateData["author"] = r.uiSettings.Author
	templateData["description"] = r.uiSettings.Description
	templateData["date_format"] = r.uiSettings.DateFormat

	templateData["instance_id"] = r.instanceId

	// enabled plugins and their settings
	//enabledPlugins := strings.Split(config.String("plugins"), " ")
	plugin := make(map[string]map[string]template.JS)
	plugins := make(map[string]bool)
	//for i := 0; i < len(enabledPlugins); i++ {
	//	plugin[enabledPlugins[i]] = make(map[string]template.JS)
	//	plugins[enabledPlugins[i]] = true
	//
	//	pluginParams := strings.Split(config.String(fmt.Sprintf("plugin.%s._", enabledPlugins[i])), " ")
	//	for j := 0; j < len(pluginParams); j++ {
	//		pluginCfgKey := fmt.Sprintf("plugin.%s.%s", enabledPlugins[i], pluginParams[j])
	//		plugin[enabledPlugins[i]][pluginParams[j]] = template.JS(config.String(pluginCfgKey))
	//	}
	//}
	templateData["plugins"] = plugins
	templateData["plugin"] = plugin

	templateName := templateData[TemplateNameDateKey].(string)
	err := r.templates[templateName].Execute(w, templateData)
	if nil != err {
		log.WithError(err).WithField("template", templateName).Error("Problems with rendering HTML template")
	}
}

func HTMLRendererData(r *http.Request, templateName string, data map[string]interface{}) interface{} {
	data[TemplateNameDateKey] = templateName
	data["url_path"] = r.URL.Path

	return data
}

func getTmplFuncMap(dateFormat string) template.FuncMap {
	funcs := gtf.GtfFuncMap

	funcs["format_date"] = func(value time.Time) string {
		return value.Format(dateFormat)
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
