package handler

import (
	"math"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

// Feed is the handler for RSS/Atom feeds
type Feed struct {
	storage    storage.Storage
	renderer   renderer.Renderer
	cfgUI      config.UI
	cfgRootURL config.RootURL
}

// NewFeed creates new Feed instance
func NewFeed(storage storage.Storage, renderer renderer.Renderer, cfgUI config.UI, cfgRootURL config.RootURL) *Feed {
	return &Feed{storage, renderer, cfgUI, cfgRootURL}
}

// Atom is the HTTP handler for Atom feed
func (h *Feed) Atom(w http.ResponseWriter, r *http.Request) {
	feed := feeds.Atom{Feed: h.getFeed(r)}
	atomFeed := feed.AtomFeed()
	if atom, err := feeds.ToXML(atomFeed); err != nil {
		panic(err)
	} else {
		h.renderer.Render(w, http.StatusOK, atom)
	}
}

// Rss is the HTTP handler for RSS feed
func (h *Feed) Rss(w http.ResponseWriter, r *http.Request) {
	feed := feeds.Rss{Feed: h.getFeed(r)}
	rssFeed := feed.RssFeed()
	if rss, err := feeds.ToXML(rssFeed); err != nil {
		panic(err)
	} else {
		h.renderer.Render(w, http.StatusOK, rss)
	}
}

func (h *Feed) getFeed(r *http.Request) *feeds.Feed {
	rootURL := h.cfgRootURL.URL(r)
	feed := &feeds.Feed{
		Title:       h.cfgUI.Title,
		Link:        &feeds.Link{Href: rootURL.String()},
		Description: h.cfgUI.Description,
		Author:      &feeds.Author{Name: h.cfgUI.Author, Email: h.cfgUI.Email},
		Copyright:   "This work is copyright © " + h.cfgUI.Author,
	}

	var items []*feeds.Item
	posts, _ := h.storage.ListAll(0)

	for _, post := range posts {
		rootURL.Path = post.Path
		item := &feeds.Item{
			Id:          post.ID,
			Title:       post.Title,
			Link:        &feeds.Link{Href: rootURL.String()},
			Description: post.Body[0:int(math.Min(float64(len(post.Body)), 255))],
			Author:      &feeds.Author{Name: h.cfgUI.Author, Email: h.cfgUI.Email},
			Created:     post.PublishedAt,
		}
		items = append(items, item)
	}

	if len(items) > 0 {
		feed.Created = items[0].Created
	} else {
		feed.Created = time.Now()
	}

	feed.Items = items
	return feed
}
