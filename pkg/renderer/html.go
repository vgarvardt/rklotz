package renderer

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/leekchan/gtf"
	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/config/plugin"
	"go.uber.org/zap"
)

const (
	templateNameDateKey = "template_name"
	dataRequestKey      = "__request"
)

// HTMLRendererConfig is configuration for HTML renderer
type HTMLRendererConfig struct {
	TemplatesPath string
	InstanceID    string
	UICfg         config.UI
	PluginsCfg    config.Plugins
	RootURLCfg    config.RootURL
}

// HTMLRenderer implements Renderer for HTML content
type HTMLRenderer struct {
	templates map[string]*template.Template
	config    HTMLRendererConfig
	logger    *zap.Logger

	enabledPluginsMap map[string]bool
	pluginsSettings   map[string]map[string]template.JS
}

// NewHTMLRenderer creates new HTMLRenderer instance
func NewHTMLRenderer(config HTMLRendererConfig, logger *zap.Logger) (*HTMLRenderer, error) {
	instance := &HTMLRenderer{
		templates: make(map[string]*template.Template),
		config:    config,
		logger:    logger,
	}

	partials, err := instance.getPartials(config.TemplatesPath, config.UICfg.Theme, config.UICfg.AboutPath)
	if nil != err {
		return nil, err
	}

	baseFiles := append(partials, fmt.Sprintf("%s/%s/base.html", config.TemplatesPath, config.UICfg.Theme))
	baseTemplate := template.Must(
		template.New("base.html").Funcs(getTmplFuncMap(config.UICfg.DateFormat)).ParseFiles(baseFiles...))

	for _, tmplName := range []string{"index.html", "post.html", "tag.html"} {
		tmplPath := fmt.Sprintf("%s/%s/%s", config.TemplatesPath, config.UICfg.Theme, tmplName)

		instance.logger.Debug("Initializing template", zap.String("name", tmplName), zap.String("path", tmplPath))
		instance.templates[tmplName] = template.Must(template.Must(baseTemplate.Clone()).ParseFiles(tmplPath))
	}

	err = instance.initPlugins()
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (r *HTMLRenderer) getPartials(templatesPath, theme, uiAbout string) ([]string, error) {
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
		r.logger.Info("Custom about panel not found, loading default theme about panel", zap.String("path", uiAbout))
		uiAbout = fmt.Sprintf("%s/%s/partial/about.html", templatesPath, theme)
	} else if nil != err {
		r.logger.Error("Failed to load custom about panel", zap.Error(err), zap.String("path", uiAbout))
		return nil, err
	} else {
		r.logger.Info("Loading custom about panel", zap.String("path", uiAbout))
	}

	partials = append(partials, uiAbout)

	return partials, nil
}

func (r *HTMLRenderer) initPlugins() error {
	r.enabledPluginsMap = make(map[string]bool, len(r.config.PluginsCfg.Enabled))
	r.pluginsSettings = make(map[string]map[string]template.JS, len(r.config.PluginsCfg.Enabled))

	for i := range r.config.PluginsCfg.Enabled {
		r.enabledPluginsMap[r.config.PluginsCfg.Enabled[i]] = true

		r.logger.Info("Loading plugin", zap.String("name", r.config.PluginsCfg.Enabled[i]))
		p, err := plugin.GetByName(r.config.PluginsCfg.Enabled[i])
		if err != nil {
			return err
		}

		r.logger.Info("Configuring plugin", zap.String("name", r.config.PluginsCfg.Enabled[i]))
		settings, err := r.config.PluginsCfg.Configure(p)
		if err != nil {
			switch e := err.(type) {
			case *plugin.ErrorConfiguring:
				r.logger.Error("Failed to configure plugin", zap.Error(err), zap.String("field", e.Field()))
			}
			return err
		}

		r.pluginsSettings[r.config.PluginsCfg.Enabled[i]] = make(map[string]template.JS)
		for settingName, settingValue := range settings {
			r.pluginsSettings[r.config.PluginsCfg.Enabled[i]][settingName] = template.JS(settingValue)
		}
	}

	return nil
}

// Render renders the data with response code to a HTTP response writer
func (r *HTMLRenderer) Render(w http.ResponseWriter, code int, data interface{}) {
	templateData := data.(map[string]interface{})

	templateData["lang"] = r.config.UICfg.Language
	templateData["title"] = r.config.UICfg.Title
	templateData["heading"] = r.config.UICfg.Heading
	templateData["intro"] = r.config.UICfg.Intro
	templateData["theme"] = r.config.UICfg.Theme
	templateData["author"] = r.config.UICfg.Author
	templateData["description"] = r.config.UICfg.Description
	templateData["date_format"] = r.config.UICfg.DateFormat

	templateData["instance_id"] = r.config.InstanceID

	templateData["plugins"] = r.enabledPluginsMap
	templateData["plugin"] = r.pluginsSettings

	rq, ok := templateData[dataRequestKey].(*http.Request)
	if ok {
		templateData["url_path"] = rq.URL.Path
		templateData["root_url"] = r.config.RootURLCfg.URL(rq).String()

		currentURL := &url.URL{}
		*currentURL = *r.config.RootURLCfg.URL(rq)
		currentURL.Path = rq.URL.Path
		templateData["current_url"] = currentURL.String()
	}

	templateName := templateData[templateNameDateKey].(string)
	err := r.templates[templateName].Execute(w, templateData)
	if nil != err {
		r.logger.Error("Problems with rendering HTML template", zap.Error(err), zap.String("template", templateName))
	}
}

// HTMLRendererData sets service fields for HTML renderer data
func HTMLRendererData(r *http.Request, templateName string, data map[string]interface{}) interface{} {
	data[templateNameDateKey] = templateName
	data[dataRequestKey] = r

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
