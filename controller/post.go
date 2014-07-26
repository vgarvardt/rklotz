package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	_ "../model"
)

func PostController(c *gin.Context) {
	c.String(200, fmt.Sprintf("%v", c.Request.URL))
}
