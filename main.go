package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ngerakines/ginpongo2"

	"./cfg"
	"./controller"
	"./model"
)

func main() {
	defer model.GetDB().Close()

	if err := model.RebuildIndex(); err != nil {
		panic(err)
	}

	if _, ok := cfg.GetApp()["rebuild"]; !ok {
		r := gin.Default()
		r.Use(ginpongo2.Pongo2())

		r.GET("/", controller.FrontController)

		r.GET("/new", controller.FormController)
		r.POST("/new", controller.FormController)
		r.GET("/edit/:uuid", controller.FormController)
		r.POST("/edit/:uuid", controller.FormController)

		r.NoRoute(controller.PostController)

		r.Static("/assets", "./assets")

		addr := cfg.String("addr")
		cfg.Log(fmt.Sprintf("Running @ %s", addr))
		r.Run(addr)
	}
}
