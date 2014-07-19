package main

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"./cfg"
	"./controller"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/", controller.FrontController)

	r.Static("/static", "./static")

	addr := cfg.String("addr")
	cfg.Log(fmt.Sprintf("Running @ %s", addr))
	r.Run(addr)
}
