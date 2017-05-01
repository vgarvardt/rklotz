package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/vgarvardt/rklotz/pkg/model"
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
