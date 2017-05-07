package handler

import (
	"net/http"
	"strconv"

	"github.com/pressly/chi"
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
	data := map[string]interface{}{
		"meta": h.storage.Meta(),
	}

	page := h.getPageFromURL(r)
	posts, _ := h.storage.ListAll(uint(page))
	data["posts"] = posts
	data["page"] = page

	tmplData := renderer.HTMLRendererData(r, "index.html", data)
	h.renderer.Render(w, http.StatusOK, tmplData)
}

func (h *PostsHandler) Tag(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"meta": h.storage.Meta(),
	}

	tag := chi.URLParam(r, "tag")

	page := h.getPageFromURL(r)
	posts, _ := h.storage.ListTag(tag, page)
	data["posts"] = posts
	data["page"] = page
	data["tag"] = tag

	tmplData := renderer.HTMLRendererData(r, "tag.html", data)
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

func (h *PostsHandler) getPageFromURL(r *http.Request) uint {
	var err error

	page := 0
	pageParam := r.URL.Query().Get("page")
	if len(pageParam) > 0 {
		if page, err = strconv.Atoi(pageParam); err != nil {
			page = 0
		}
	}

	return uint(page)
}
