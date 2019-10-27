package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// Posts is the handler for posts pages
type Posts struct {
	storage  storage.Storage
	renderer renderer.Renderer
}

// NewPosts creates new Posts instance
func NewPosts(storage storage.Storage, renderer renderer.Renderer) *Posts {
	return &Posts{storage, renderer}
}

// Front is the HTTP handler for the front page with post list
func (h *Posts) Front(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"meta": h.storage.Meta(),
	}

	page := h.getPageFromURL(r)
	posts, _ := h.storage.ListAll(page)
	data["posts"] = posts
	data["page"] = page

	tmplData := renderer.HTMLRendererData(r, "index.html", data)
	h.renderer.Render(w, http.StatusOK, tmplData)
}

// Tag is the HTTP handler for the tag page with post list for a tag
func (h *Posts) Tag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	data := map[string]interface{}{
		"meta": h.storage.TagMeta(tag),
	}

	page := h.getPageFromURL(r)
	posts, _ := h.storage.ListTag(tag, page)
	data["posts"] = posts
	data["page"] = page
	data["tag"] = tag

	tmplData := renderer.HTMLRendererData(r, "tag.html", data)
	h.renderer.Render(w, http.StatusOK, tmplData)
}

// Post is the HTTP handler for the post page
func (h *Posts) Post(w http.ResponseWriter, r *http.Request) {
	post, err := h.storage.FindByPath(r.URL.Path)

	if err != nil {
		code := map[bool]int{true: http.StatusNotFound, false: http.StatusInternalServerError}
		w.WriteHeader(code[err == storage.ErrorNotFound])
		w.Write([]byte(err.Error()))
	} else {
		tmplData := renderer.HTMLRendererData(r, "post.html", map[string]interface{}{"post": post})
		h.renderer.Render(w, http.StatusOK, tmplData)
	}
}

func (h *Posts) getPageFromURL(r *http.Request) int {
	var err error

	page := 0
	pageParam := r.URL.Query().Get("page")
	if len(pageParam) > 0 {
		if page, err = strconv.Atoi(pageParam); err != nil {
			page = 0
		}
	}

	return page
}
