package controller

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"../model"
)

func PostController(c *gin.Context) {
	if strings.Trim(fmt.Sprintf("%v", c.Request.URL), "/") == "@" {
		redirect(c, "/@/drafts")
		return
	}

	post := new(model.Post)
	post.LoadByPath(strings.Trim(fmt.Sprintf("%v", c.Request.URL), "/"))

	if len(post.UUID) > 0 {
		ctx := make(map[string]interface{})
		ctx["post"] = post
		render(c, "post.html", ctx)
	} else {
		c.Data(404, gin.MIMEPlain, []byte("404 page not found"))
	}
}

func AutoComplete(c *gin.Context) {
	q := c.Request.URL.Query()["q"][0]
	var tags []string
	tags = append(tags, q)
	c.JSON(200, gin.H{"tags": append(tags, model.AutoCompleteTags(q)...)})
}
