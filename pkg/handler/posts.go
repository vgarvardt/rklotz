package handler

import (
	"net/http"
	"strconv"

	"github.com/pressly/chi"
	"github.com/vgarvardt/rklotz/pkg/model"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/repository"
)

type PostsHandler struct {
	storage  repository.Storage
	renderer renderer.Renderer
}

func NewPostsHandler(storage repository.Storage, renderer renderer.Renderer) *PostsHandler {
	return &PostsHandler{storage, renderer}
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
	post, err := h.storage.FindByPath(r.URL.Path)

	if err != nil {
		code := map[bool]int{true: http.StatusNotFound, false: http.StatusInternalServerError}
		w.WriteHeader(code[err == repository.ErrorNotFound])
		w.Write([]byte(err.Error()))
	} else {
		tmplData := renderer.HTMLRendererData(r, "index.html", map[string]interface{}{"post": post})
		h.renderer.Render(w, http.StatusOK, tmplData)
	}
}
