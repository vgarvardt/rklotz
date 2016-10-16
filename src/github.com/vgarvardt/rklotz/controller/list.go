package controller

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/vgarvardt/rklotz/model"
)

func DraftsController(ctx echo.Context) error {
	return ctx.Render(
		http.StatusOK,
		"@/drafts.html",
		map[string]interface{}{"posts": model.MustGetDraftPosts()},
	)
}

func PublishedController(ctx echo.Context) error {
	return ctx.Render(
		http.StatusOK,
		"@/published.html",
		map[string]interface{}{"posts": model.MustGetPublishedPosts()},
	)
}

func TagController(ctx echo.Context) error {
	tag := ctx.Param("tag")

	return ctx.Render(
		http.StatusOK,
		"tag.html",
		map[string]interface{}{"tag": tag, "posts": model.MustGetTagPosts(tag)},
	)
}
