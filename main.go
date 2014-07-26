package main

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"./cfg"
	"./controller"
	"./model"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.GET("/", controller.FrontController)

	r.GET("/new", controller.FormController)
	r.POST("/new", controller.FormController)
	r.GET("/edit/:uuid", controller.FormController)
	r.POST("/edit/:uuid", controller.FormController)

	r.Static("/assets", "./assets")

	defer model.GetDB().Close()

	addr := cfg.String("addr")
	cfg.Log(fmt.Sprintf("Running @ %s", addr))
	r.Run(addr)
}
