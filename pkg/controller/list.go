package controller

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/vgarvardt/rklotz/pkg/model"
)

func TagController(ctx echo.Context) error {
	tag := ctx.Param("tag")

	return ctx.Render(
		http.StatusOK,
		"tag.html",
		map[string]interface{}{"tag": tag, "posts": model.MustGetTagPosts(tag)},
	)
}
