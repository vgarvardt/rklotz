package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/vgarvardt/rklotz/model"
)

func FrontController(c *gin.Context) {
	meta := model.NewLoadedMeta()
	ctx := gin.H{"meta": meta}

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
