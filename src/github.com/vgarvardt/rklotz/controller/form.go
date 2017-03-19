package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/vgarvardt/rklotz/model"
)

const (
	FORM_VALUE_DELETE  = "delete"
	FORM_VALUE_PREVIEW = "preview"
	FORM_VALUE_DRAFT   = "draft"
	FORM_VALUE_PUBLISH = "publish"
)

func FormController(ctx echo.Context) error {
	post := new(model.Post)
	uuid := ctx.Param("uuid")
	if len(uuid) > 0 {
		post.Load(uuid)
	} else {
		post.PublishedAt = time.Now()
	}
	data := map[string]interface{}{
		"post":             post,
		"formValueDelete":  FORM_VALUE_DELETE,
		"formValuePreview": FORM_VALUE_PREVIEW,
		"formValueDraft":   FORM_VALUE_DRAFT,
		"formValuePublish": FORM_VALUE_PUBLISH,
	}

	data["formats"] = model.GetAvailableFormats()
	if ctx.Request().Method() == http.MethodPost {
		if ctx.FormValue("op") == FORM_VALUE_DELETE {
			post.Delete()
			if post.Draft {
				return ctx.Redirect(http.StatusFound, "/@/drafts")
			} else {
				return ctx.Redirect(http.StatusFound, "/@/published")
			}
		} else {
			if err := post.Bind(ctx); err != nil {
				return err
			}
			data["post"] = post

			if ctx.FormValue("op") == FORM_VALUE_PREVIEW {
				post.ReFormat()
				data["preview"] = true
				data["post"] = post
			} else {
				if formErrors := post.Validate(); len(formErrors) > 0 {
					data["alert_warning"] = "Please fix form values"
					data["errors"] = formErrors
				} else {
					post.Save(ctx.FormValue("op") == FORM_VALUE_DRAFT)
					return ctx.Redirect(http.StatusFound, "/@/edit/"+post.UUID)
				}
			}
		}
	}

	return ctx.Render(http.StatusOK, "@/form.html", data)
}
