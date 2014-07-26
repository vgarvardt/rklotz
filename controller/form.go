package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/flosch/pongo2"

	"../model"
)

func FormController(c *gin.Context) {
	ctx := pongo2.Context{"formats": model.GetAvailableFormats()}

	post := new(model.Post)
	uuid := c.Params.ByName("uuid")
	if len(uuid) > 0 {
		post.Load(uuid)
		ctx["post"] = post
		//0dd9c29d-5fbc-472f-ab22-04562302dd28
	}

	if c.Request.Method == "POST" {
		auth := new(model.Auth)
		auth.Bind(c.Request)

		post.Bind(c.Request)
		ctx["post"] = post

		if c.Request.FormValue("op") == "preview" {
			post.ReFormat()
			ctx["preview"] = true
			ctx["post"] = post
		} else {
			if auth.IsValid() {
				post.Save(c.Request.FormValue("op") == "draft")
				redirect(c, "/edit/" + post.UUID)
				return
			} else {
				ctx["alert_warning"] = "Failed to authenticate with given values"
			}
		}
	}

	pongo2Render(c, "form.html", ctx)
}
