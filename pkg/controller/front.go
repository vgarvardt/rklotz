package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/vgarvardt/rklotz/pkg/model"
)

func FrontController(ctx echo.Context) error {
	meta := model.NewLoadedMeta(10)
	data := map[string]interface{}{
		"meta": meta,
	}

	if meta.Posts > 0 {
		page := 0
		var err error

		pageParam := ctx.QueryParam("page")
		if len(pageParam) > 0 {
			if page, err = strconv.Atoi(pageParam); err != nil {
				page = 0
			}
		}

		posts, _ := model.GetPostsPage(page)
		data["posts"] = posts
		data["page"] = page
	}

	return ctx.Render(http.StatusOK, "index.html", data)
}
