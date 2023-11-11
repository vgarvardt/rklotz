package handler

import (
	"net/http"

	"github.com/vgarvardt/rklotz/pkg/server/renderer"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// Feed is the handler for RSS/Atom feeds
type Feed struct {
	storage  storage.Storage
	renderer renderer.Renderer
}

// NewFeed creates new Feed instance
func NewFeed(s storage.Storage, r renderer.Renderer) *Feed {
	return &Feed{s, r}
}

// Atom is the HTTP handler for Atom feed
func (h *Feed) Atom(w http.ResponseWriter, r *http.Request) {
	h.feed(w, r, "atom")
}

// Rss is the HTTP handler for RSS feed
func (h *Feed) Rss(w http.ResponseWriter, r *http.Request) {
	h.feed(w, r, "rss")
}

func (h *Feed) feed(w http.ResponseWriter, r *http.Request, template string) {
	posts, err := h.storage.ListAll(0)
	if err != nil {
		panic(err)
	}

	h.renderer.Render(w, http.StatusOK, renderer.NewData(r, template, renderer.D{"posts": posts}))
}
