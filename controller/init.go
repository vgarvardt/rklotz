package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/flosch/pongo2"

	"../cfg"
)

func pongo2Render(c *gin.Context, template string, ctx pongo2.Context) {
	ui := []string{"lang", "title", "heading", "intro", "theme", "author", "description"}
	for i := 0; i < len(ui); i++ {
		if _, ok := ctx[ui[i]]; !ok {
			ctx[ui[i]] = cfg.String("ui." + ui[i])
		}
	}

	tpl, err := pongo2.FromFile(fmt.Sprintf("./templates/%v/%s", ctx["theme"], template))
	if err != nil {
		panic(err)
	}

	out, err := tpl.Execute(ctx)
	if err != nil {
		panic(err)
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	c.Writer.WriteHeader(200)
	c.Writer.Write([]byte(out))
}

func redirect(c *gin.Context, url string) {
	http.Redirect(c.Writer, c.Request, url, http.StatusFound)
}
