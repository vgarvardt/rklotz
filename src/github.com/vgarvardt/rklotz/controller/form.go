package controller

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vgarvardt/rklotz/model"
)

func FormController(c *gin.Context) {
	ctx := make(map[string]interface{})

	post := new(model.Post)
	uuid := c.Params.ByName("uuid")
	if len(uuid) > 0 {
		post.Load(uuid)
	} else {
		post.PublishedAt = time.Now()
	}
	ctx["post"] = post

	ctx["formats"] = model.GetAvailableFormats()
	if c.Request.Method == "POST" {
		if c.Request.FormValue("op") == "delete" {
			post.Delete()
			if post.Draft {
				redirect(c, "/@/drafts")
			} else {
				redirect(c, "/@/published")
			}
			return
		} else {
			post.Bind(c)
			ctx["post"] = post

			if c.Request.FormValue("op") == "preview" {
				post.ReFormat()
				ctx["preview"] = true
				ctx["post"] = post
			} else {
				if formErrors := post.Validate(); len(formErrors) > 0 {
					ctx["alert_warning"] = "Please fix form values"
					ctx["errors"] = formErrors
				} else {
					post.Save(c.Request.FormValue("op") == "draft")
					redirect(c, "/@/edit/"+post.UUID)
					return
				}
			}
		}
	}

	render(c, "@/form.html", ctx)
}
