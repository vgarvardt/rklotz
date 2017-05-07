package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pressly/chi"
	"github.com/vgarvardt/rklotz/pkg/model"
	"github.com/vgarvardt/rklotz/pkg/renderer"
)

type PostsHandler struct {
	renderer renderer.Renderer
}

func NewPostsHandler(renderer renderer.Renderer) *PostsHandler {
	return &PostsHandler{renderer}
}

func (h *PostsHandler) Front(w http.ResponseWriter, r *http.Request) {
	meta := model.NewLoadedMeta(10)
	data := map[string]interface{}{
		"meta": meta,
	}

	if meta.Posts > 0 {
		page := 0
		var err error

		pageParam := r.URL.Query().Get("page")
		if len(pageParam) > 0 {
			if page, err = strconv.Atoi(pageParam); err != nil {
				page = 0
			}
		}

		posts, _ := model.GetPostsPage(page)
		data["posts"] = posts
		data["page"] = page
	}

	tmplData := renderer.HTMLRendererData(r, "index.html", data)
	h.renderer.Render(w, http.StatusOK, tmplData)
}

func (h *PostsHandler) Tag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	tmplData := renderer.HTMLRendererData(
		r,
		"tag.html",
		map[string]interface{}{"tag": tag, "posts": model.MustGetTagPosts(tag)},
	)
	h.renderer.Render(w, http.StatusOK, tmplData)
}

func (h *PostsHandler) Post(w http.ResponseWriter, r *http.Request) {
	post := new(model.Post)
	post.LoadByPath(strings.Trim(fmt.Sprintf("%v", r.URL.Path), "/"))

	if len(post.Path) > 0 {
		tmplData := renderer.HTMLRendererData(r, "index.html", map[string]interface{}{"post": post})
		h.renderer.Render(w, http.StatusOK, tmplData)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
