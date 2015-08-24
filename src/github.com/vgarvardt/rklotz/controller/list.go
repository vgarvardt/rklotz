package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/vgarvardt/rklotz/model"
)

func DraftsController(c *gin.Context) {
	ctx := make(map[string]interface{})

	posts, _ := model.GetDraftPosts()
	ctx["posts"] = posts

	render(c, "@/drafts.html", ctx)
}

func PublishedController(c *gin.Context) {
	ctx := make(map[string]interface{})

	posts, _ := model.GetPublishedPosts()
	ctx["posts"] = posts

	render(c, "@/published.html", ctx)
}

func TagController(c *gin.Context) {
	ctx := make(map[string]interface{})

	tag := c.Params.ByName("tag")
	posts, _ := model.GetTagPosts(tag)
	ctx["tag"] = tag
	ctx["posts"] = posts

	render(c, "tag.html", ctx)
}
