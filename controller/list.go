package controller

import (
	"github.com/gin-gonic/gin"

	"../model"
)

func DraftsController(c *gin.Context) {
	ctx := make(map[string]interface{})

	posts, _ := model.GetDraftPosts()
	ctx["posts"] = posts

	render(c, "drafts.html", ctx)
}

func TagController(c *gin.Context) {
	ctx := make(map[string]interface{})

	tag := c.Params.ByName("tag")
	posts, _ := model.GetTagPosts(tag)
	ctx["tag"] = tag
	ctx["posts"] = posts

	render(c, "tag.html", ctx)
}
