package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/flosch/pongo2"

	"../cfg"
)

func pongo2Render(c *gin.Context, templatePath string, ctx pongo2.Context) {
	ctx["title"] = cfg.String("ui.title")
	ctx["heading"] = cfg.String("ui.heading")
	ctx["intro"] = cfg.String("ui.intro")

	tpl, _ := pongo2.FromFile(templatePath)
	out, _ := tpl.Execute(ctx)

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(200)
	c.Writer.Write([]byte(out))
}
