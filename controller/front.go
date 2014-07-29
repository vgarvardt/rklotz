package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"../model"
)

func FrontController(c *gin.Context) {
	ctx := make(map[string]interface{})

	meta := new(model.Meta)
	meta.Load()
	ctx["meta"] = meta

	if meta.Posts > 0 {
		page := 0
		var err error

		pageParam := c.Request.URL.Query()["page"]
		if len(pageParam) > 0 {
			if page, err = strconv.Atoi(pageParam[0]); err != nil {
				page = 0
			}
		}

		posts, _ := model.GetPostsPage(page)
		ctx["posts"] = posts
		ctx["page"] = page
	}

	render(c, "index.html", ctx)
}
