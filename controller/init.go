package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/flosch/pongo2"

	"../cfg"
)

func pongo2Render(c *gin.Context, template string, ctx pongo2.Context) {
	ui := []string{"lang", "title", "heading", "intro", "theme"}
	for i := 0; i < len(ui); i++ {
		ctx[ui[i]] = cfg.String("ui." + ui[i])
	}

	tpl, _ := pongo2.FromFile(fmt.Sprintf("./templates/%v/%s", ctx["theme"], template))
	out, _ := tpl.Execute(ctx)

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(200)
	c.Writer.Write([]byte(out))
}
