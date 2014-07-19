package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/flosch/pongo2"
)

func pongo2Render(c *gin.Context, templatePath string, ctx pongo2.Context) {
	tpl, _ := pongo2.FromFile(templatePath)
	out, _ := tpl.Execute(ctx)

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(200)
	c.Writer.Write([]byte(out))
}
