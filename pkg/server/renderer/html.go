package renderer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/cappuccinotm/slogx"
	"github.com/leekchan/gtf"

	"github.com/vgarvardt/rklotz/pkg/server/plugin"
	"github.com/vgarvardt/rklotz/pkg/server/rqctx"
)

// HTMLConfig is configuration for HTML renderer
type HTMLConfig struct {
	Debug         bool
	TemplatesPath string
	UICfg         UIConfig
	RootURLCfg    RootURLConfig
	PluginsCfg    plugin.Config
}

// HTML implements Renderer for HTML content
type HTML struct {
	templates  map[string]*template.Template
	config     HTMLConfig
	logger     *slog.Logger
	instanceID string

	enabledPluginsMap map[string]bool
	pluginsSettings   map[string]map[string]template.JS
}

// NewHTML creates new HTML instance
func NewHTML(config HTMLConfig, logger *slog.Logger) (*HTML, error) {
	instance := &HTML{
		templates: make(map[string]*template.Template),
		config:    config,
		logger:    logger,
	}

	return instance, instance.initTemplates()
}

func (r *HTML) newID() string {
	hash := sha256.Sum256([]byte(time.Now().Format(time.RFC3339Nano)))
	return hex.EncodeToString(hash[:])[:6]
}

func (r *HTML) initTemplates() error {
	r.instanceID = r.newID()

	baseFiles, err := r.getPartials(r.config.TemplatesPath, r.config.UICfg.Theme, r.config.UICfg.AboutPath)
	if nil != err {
		return err
	}

	baseFiles = append(baseFiles, fmt.Sprintf("%s/%s/base.tpl", r.config.TemplatesPath, r.config.UICfg.Theme))
	baseTemplate := template.Must(
		template.New("base.tpl").
			Funcs(getTmplFuncMap(r.config.UICfg.DateFormat)).
			ParseFiles(baseFiles...),
	)

	for _, tmplName := range []string{"404.tpl", "500.tpl", "index.tpl", "post.tpl", "tag.tpl"} {
		tmplPath := fmt.Sprintf("%s/%s/%s", r.config.TemplatesPath, r.config.UICfg.Theme, tmplName)

		r.logger.Debug("Initializing template", slog.String("name", tmplName), slog.String("path", tmplPath))
		r.templates[tmplName] = template.Must(template.Must(baseTemplate.Clone()).ParseFiles(tmplPath))
	}

	err = r.initPlugins()
	if err != nil {
		return err
	}

	return nil
}

func (r *HTML) getPartials(templatesPath, theme, uiAbout string) ([]string, error) {
	var partials []string

	walkFn := func(path string, f os.FileInfo, err error) error {
		if nil == err && !f.IsDir() && !strings.HasSuffix(path, "about.tpl") {
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
	switch {
	case os.IsNotExist(err):
		r.logger.Info("Custom about panel not found, loading default theme about panel", slog.String("path", uiAbout))
		uiAbout = fmt.Sprintf("%s/%s/partial/about.tpl", templatesPath, theme)
	case err != nil:
		r.logger.Error("Failed to load custom about panel", slogx.Error(err), slog.String("path", uiAbout))
		return nil, err
	default:
		r.logger.Info("Loading custom about panel", slog.String("path", uiAbout))
	}

	partials = append(partials, uiAbout)

	return partials, nil
}

func (r *HTML) initPlugins() error {
	r.enabledPluginsMap = make(map[string]bool, len(r.config.PluginsCfg.Enabled))
	r.pluginsSettings = make(map[string]map[string]template.JS, len(r.config.PluginsCfg.Enabled))

	for i := range r.config.PluginsCfg.Enabled {
		r.enabledPluginsMap[r.config.PluginsCfg.Enabled[i]] = true

		r.logger.Info("Loading plugin", slog.String("name", r.config.PluginsCfg.Enabled[i]))
		p, err := plugin.GetByName(r.config.PluginsCfg.Enabled[i])
		if err != nil {
			return err
		}

		r.logger.Info("Configuring plugin", slog.String("name", r.config.PluginsCfg.Enabled[i]))
		settings, err := r.config.PluginsCfg.SetUp(p)
		if err != nil {
			var e *plugin.ErrorConfiguring
			if errors.As(err, &e) {
				r.logger.Error("Failed to configure plugin", slogx.Error(err), slog.String("field", e.Field()))
			}
			return err
		}

		r.pluginsSettings[r.config.PluginsCfg.Enabled[i]] = make(map[string]template.JS)
		for settingName, settingValue := range settings {
			/* #nosec G203 -- plugins are supposed to be non-safe and contain JS */
			r.pluginsSettings[r.config.PluginsCfg.Enabled[i]][settingName] = template.JS(settingValue)
		}
	}

	return nil
}

// Render renders the data with response code to a HTTP response writer
func (r *HTML) Render(w http.ResponseWriter, code int, data *Data) {
	data.m.RLock()
	defer data.m.RUnlock()

	logger := rqctx.GetLogger(data.r.Context())

	if r.config.Debug {
		logger.Warn("HTML renderer is in the debug mode, reloading all templates")
		if err := r.initTemplates(); err != nil {
			logger.Error(
				"Problems with reloading HTML templates",
				slogx.Error(err),
			)
			return
		}
	}

	templateData := data.data

	templateData["lang"] = r.config.UICfg.Language
	templateData["title"] = r.config.UICfg.Title
	templateData["heading"] = r.config.UICfg.Heading
	templateData["intro"] = r.config.UICfg.Intro
	templateData["theme"] = r.config.UICfg.Theme
	templateData["author"] = r.config.UICfg.Author
	templateData["description"] = r.config.UICfg.Description
	templateData["date_format"] = r.config.UICfg.DateFormat

	templateData["instance_id"] = r.instanceID
	templateData["plugins"] = r.enabledPluginsMap
	templateData["plugin"] = r.pluginsSettings

	templateData["url_path"] = data.r.URL.Path
	templateData["root_url"] = r.config.RootURLCfg.URL(data.r).String()

	currentURL := &url.URL{}
	*currentURL = *r.config.RootURLCfg.URL(data.r)
	currentURL.Path = data.r.URL.Path
	templateData["current_url"] = currentURL.String()

	tmpl, found := r.templates[data.template]
	if !found {
		logger.Error("Template is not found in the templates registry", slog.String("template", data.template))
		panic(fmt.Errorf("template is not found in the templates registry %q", data.template))
	}

	w.WriteHeader(code)

	if err := tmpl.Execute(w, templateData); nil != err {
		logger.Error("Problems with rendering HTML template", slogx.Error(err), slog.String("template", data.template))
	}
}

// Error renders error to a response writer
func (r *HTML) Error(rq *http.Request, w http.ResponseWriter, code int, err error) {
	r.Render(w, code, NewData(rq, fmt.Sprintf("%d.tpl", code), D{"error": err.Error()}))
}

func getTmplFuncMap(dateFormat string) template.FuncMap {
	funcs := gtf.GtfFuncMap

	funcs["format_date"] = func(value time.Time) string {
		return value.Format(dateFormat)
	}
	funcs["add"] = func(arg int, value int) int {
		return value + arg
	}
	funcs["noescape"] = func(value string) template.HTML {
		/* #nosec G203 -- function is supposed to be non-safe and contain JS */
		return template.HTML(value)
	}

	return funcs
}
