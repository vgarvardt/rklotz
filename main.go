package main

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"./cfg"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	addr := cfg.String("addr")
	cfg.Log(fmt.Sprintf("Running @ %s", addr))
	r.Run(addr)
}
