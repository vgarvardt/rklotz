package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"../cfg"
)

func render(c *gin.Context, template string, ctx map[string]interface{}) {
	ui := []string{"lang", "title", "heading", "intro", "theme", "author", "description"}
	for i := 0; i < len(ui); i++ {
		if _, ok := ctx[ui[i]]; !ok {
			ctx[ui[i]] = cfg.String("ui." + ui[i])
		}
	}

	ctx["instanceId"] = cfg.GetInstanceId()
	ctx["url_path"] = c.Request.URL.Path

	c.Set("template", fmt.Sprintf("./templates/%v/%s", ctx["theme"], template))
	c.Set("data", ctx)
	c.Writer.WriteHeader(200)
}

func redirect(c *gin.Context, url string) {
	http.Redirect(c.Writer, c.Request, url, http.StatusFound)
}
