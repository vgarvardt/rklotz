package controller

import (
	"github.com/gin-gonic/gin"

	"../model"
)

func FormController(c *gin.Context) {
	ctx := make(map[string]interface{})

	post := new(model.Post)
	uuid := c.Params.ByName("uuid")
	if len(uuid) > 0 {
		post.Load(uuid)
		ctx["post"] = post
	}

	ctx["formats"] = model.GetAvailableFormats()
	if c.Request.Method == "POST" {
		auth := new(model.Auth)
		auth.Bind(c.Request)

		if err := post.Bind(c.Request); err != nil {
			panic(err)
		}
		ctx["post"] = post

		if c.Request.FormValue("op") == "preview" {
			post.ReFormat()
			ctx["preview"] = true
			ctx["post"] = post
		} else {
			if auth.IsValid() {
				if formErrors := post.Validate(); len(formErrors) > 0 {
					ctx["alert_warning"] = "Please fix form values"
					ctx["errors"] = formErrors
				} else {
					post.Save(c.Request.FormValue("op") == "draft")
					redirect(c, "/edit/" + post.UUID)
					return
				}
			} else {
				ctx["alert_warning"] = "Failed to authenticate with given values"
			}
		}
	}

	render(c, "form.html", ctx)
}
