package controller

import (
	"github.com/gin-gonic/gin"
)

func FrontController(c *gin.Context) {
	ctx := make(map[string]interface{})
	render(c, "index.html", ctx)
}
