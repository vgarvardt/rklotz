package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/vgarvardt/rklotz/pkg/server/renderer"
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
	page := h.getPageFromURL(r)
	posts, _ := h.storage.ListAll(page)

	h.renderer.Render(w, http.StatusOK, renderer.NewData(r, "index.html", renderer.D{
		"meta":  h.storage.Meta(),
		"posts": posts,
		"page":  page,
	}))
}

// Tag is the HTTP handler for the tag page with post list for a tag
func (h *Posts) Tag(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	page := h.getPageFromURL(r)
	posts, _ := h.storage.ListTag(tag, page)

	h.renderer.Render(w, http.StatusOK, renderer.NewData(r, "tag.html", renderer.D{
		"meta":  h.storage.TagMeta(tag),
		"posts": posts,
		"page":  page,
		"tag":   tag,
	}))
}

// Post is the HTTP handler for the post page
func (h *Posts) Post(w http.ResponseWriter, r *http.Request) {
	post, err := h.storage.FindByPath(r.URL.Path)

	if err != nil {
		status := map[bool]int{
			true:  http.StatusNotFound,
			false: http.StatusInternalServerError,
		}[err == storage.ErrorNotFound]

		http.Error(w, err.Error(), status)
		return
	}

	h.renderer.Render(w, http.StatusOK, renderer.NewData(r, "post.html", renderer.D{"post": post}))

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
