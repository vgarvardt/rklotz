package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/vgarvardt/rklotz/model"
)

func PostController(ctx echo.Context) error {
	post := new(model.Post)
	post.LoadByPath(strings.Trim(fmt.Sprintf("%v", ctx.Request().URL().Path()), "/"))

	if len(post.UUID) > 0 {
		return ctx.Render(http.StatusOK, "post.html", map[string]interface{}{"post": post})
	} else {
		return ctx.NoContent(http.StatusNotFound)
	}
}

func AutoComplete(ctx echo.Context) error {
	q := ctx.Request().URL().QueryParam("q")
	var tags []string
	tags = append(tags, q)
	return ctx.JSON(http.StatusOK, map[string]interface{}{"tags": append(tags, model.AutoCompleteTags(q)...)})
}
