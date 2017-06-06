package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/leekchan/gtf"
	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/config/plugin"
)

const TemplateNameDateKey = "template_name"

type htmlRenderer struct {
	templates  map[string]*template.Template
	instanceId string
	uiSettings config.UISetting
	plugins    config.Plugins

	enabledPluginsMap map[string]bool
	pluginsSettings   map[string]map[string]template.JS
}

func NewHTMLRenderer(templatesPath string, instanceId string, uiSettings config.UISetting, plugins config.Plugins) (*htmlRenderer, error) {
	instance := &htmlRenderer{
		templates:  make(map[string]*template.Template),
		instanceId: instanceId,
		uiSettings: uiSettings,
		plugins:    plugins,
	}

	partials, err := instance.getPartials(templatesPath, uiSettings.Theme, uiSettings.AboutPath)
	if nil != err {
		return nil, err
	}

	baseFiles := append(partials, fmt.Sprintf("%s/%s/base.html", templatesPath, uiSettings.Theme))
	baseTemplate := template.Must(
		template.New("base.html").Funcs(getTmplFuncMap(uiSettings.DateFormat)).ParseFiles(baseFiles...))

	for _, tmplName := range []string{"index.html", "post.html", "tag.html"} {
		tmplPath := fmt.Sprintf("%s/%s/%s", templatesPath, uiSettings.Theme, tmplName)

		log.WithFields(log.Fields{"name": tmplName, "path": tmplPath}).Debug("Initializing template")
		instance.templates[tmplName] = template.Must(template.Must(baseTemplate.Clone()).ParseFiles(tmplPath))
	}

	err = instance.initPlugins()
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (r *htmlRenderer) getPartials(templatesPath, theme, uiAbout string) ([]string, error) {
	var partials []string

	walkFn := func(path string, f os.FileInfo, err error) error {
		if nil == err && !f.IsDir() && !strings.HasSuffix(path, "about.html") {
			partials = append(partials, path)
		}
		return err
	}

	pluginsPath := path.Join(templatesPath, "plugins")
	err := filepath.Walk(pluginsPath, walkFn)
	if err != nil {
		return nil, err
	}

	themePartialsPath := path.Join(templatesPath, theme, "partial")
	err = filepath.Walk(themePartialsPath, walkFn)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(uiAbout)
	if os.IsNotExist(err) {
		log.WithField("path", uiAbout).Info("Custom about panel not found, loading default theme about panel")
		uiAbout = fmt.Sprintf("%s/%s/partial/about.html", templatesPath, theme)
	} else if nil != err {
		log.WithError(err).WithField("path", uiAbout).Error("Failed to load custom about panel")
		return nil, err
	} else {
		log.WithField("path", uiAbout).Info("Loading custom about panel")
	}

	partials = append(partials, uiAbout)

	return partials, nil
}

func (r *htmlRenderer) initPlugins() error {
	r.enabledPluginsMap = make(map[string]bool, len(r.plugins.Enabled))
	r.pluginsSettings = make(map[string]map[string]template.JS, len(r.plugins.Enabled))

	for i := range r.plugins.Enabled {
		r.enabledPluginsMap[r.plugins.Enabled[i]] = true

		log.WithField("name", r.plugins.Enabled[i]).Info("Loading plugin")
		p, err := plugin.GetByName(r.plugins.Enabled[i])
		if err != nil {
			return err
		}

		log.WithField("name", r.plugins.Enabled[i]).Info("Configuring plugin")
		settings, err := r.plugins.Configure(p)
		if err != nil {
			return err
		}

		r.pluginsSettings[r.plugins.Enabled[i]] = make(map[string]template.JS)
		for settingName, settingValue := range settings {
			r.pluginsSettings[r.plugins.Enabled[i]][settingName] = template.JS(settingValue)
		}
	}

	return nil
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

	templateData["plugins"] = r.enabledPluginsMap
	templateData["plugin"] = r.pluginsSettings

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
