package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/flosch/pongo2"
)

func FrontController(c *gin.Context) {
	pongo2Render(c, "./view/index.html", pongo2.Context{})
}
