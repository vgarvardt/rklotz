package handler

import (
	"net/http"

	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// Feed is the handler for RSS/Atom feeds
type Feed struct {
	storage  storage.Storage
	renderer renderer.Renderer
}

// NewFeed creates new Feed instance
func NewFeed(storage storage.Storage, renderer renderer.Renderer) *Feed {
	return &Feed{storage, renderer}
}

// Atom is the HTTP handler for Atom feed
func (h *Feed) Atom(w http.ResponseWriter, r *http.Request) {
	posts, err := h.storage.ListAll(0)
	if err != nil {
		panic(err)
	}

	h.renderer.Render(w, http.StatusOK, renderer.NewData(r, "atom", map[string]interface{}{"posts": posts}))

}

// Rss is the HTTP handler for RSS feed
func (h *Feed) Rss(w http.ResponseWriter, r *http.Request) {
	posts, err := h.storage.ListAll(0)
	if err != nil {
		panic(err)
	}

	h.renderer.Render(w, http.StatusOK, renderer.NewData(r, "rss", map[string]interface{}{"posts": posts}))
}
