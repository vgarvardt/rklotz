package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vgarvardt/rklotz/cfg"
)

func render(c *gin.Context, template string, ctx map[string]interface{}) {
	ui := []string{"lang", "title", "heading", "intro", "theme", "author", "description", "date_format"}
	for i := 0; i < len(ui); i++ {
		if _, ok := ctx[ui[i]]; !ok {
			ctx[ui[i]] = cfg.String("ui." + ui[i])
		}
	}

	ctx["instance_id"] = cfg.GetInstanceId()
	ctx["url_path"] = c.Request.URL.Path

	uiAbout := cfg.String("ui.about")
	if len(uiAbout) > 0 {
		ctx["about_path"] = uiAbout
	} else {
		ctx["about_path"] = "./partial/about.html"
	}

	// enabled plugins and their settings
	enabledPlugins := strings.Split(cfg.String("plugins"), " ")
	plugin := make(map[string]map[string]string)
	plugins := make(map[string]bool)
	for i := 0; i < len(enabledPlugins); i++ {
		plugin[enabledPlugins[i]] = make(map[string]string)
		plugins[enabledPlugins[i]] = true

		pluginParams := strings.Split(cfg.String(fmt.Sprintf("plugin.%s._", enabledPlugins[i])), " ")
		for j := 0; j < len(pluginParams); j++ {
			plugin[enabledPlugins[i]][pluginParams[j]] = cfg.String(fmt.Sprintf("plugin.%s.%s", enabledPlugins[i], pluginParams[j]))
		}
	}
	ctx["plugins"] = plugins
	ctx["plugin"] = plugin

	c.Set("template", fmt.Sprintf("%s/templates/%v/%s", cfg.GetRootDir(), ctx["theme"], template))
	c.Set("data", ctx)
	c.Writer.WriteHeader(200)
}

func redirect(c *gin.Context, url string) {
	http.Redirect(c.Writer, c.Request, url, http.StatusFound)
}
