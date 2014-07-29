package controller

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"../model"
)

func PostController(c *gin.Context) {
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