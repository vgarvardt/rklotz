package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/vgarvardt/rklotz/model"
)

func DraftsController(c *gin.Context) {
	posts, _ := model.GetDraftPosts()

	render(c, "@/drafts.html", gin.H{"posts": posts})
}

func PublishedController(c *gin.Context) {
	posts, _ := model.GetPublishedPosts()

	render(c, "@/published.html", gin.H{"posts": posts})
}

func TagController(c *gin.Context) {
	tag := c.Params.ByName("tag")
	posts, _ := model.GetTagPosts(tag)

	render(c, "tag.html", gin.H{"tag": tag, "posts": posts})
}
