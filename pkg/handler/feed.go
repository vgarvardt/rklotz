package handler

import (
	"math"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/model"
	"github.com/vgarvardt/rklotz/pkg/renderer"
)

type FeedHandler struct {
	renderer   renderer.Renderer
	uiSettings config.UISetting
	rootUrl    config.RootURL
}

func NewFeedHandler(renderer renderer.Renderer, uiSettings config.UISetting, rootUrl config.RootURL) *FeedHandler {
	return &FeedHandler{renderer, uiSettings, rootUrl}
}

func (h *FeedHandler) Atom(w http.ResponseWriter, r *http.Request) {
	feed := feeds.Atom{Feed: h.getFeed(r)}
	atomFeed := feed.AtomFeed()
	if atom, err := feeds.ToXML(atomFeed); err != nil {
		panic(err)
	} else {
		h.renderer.Render(w, http.StatusOK, atom)
	}
}

func (h *FeedHandler) Rss(w http.ResponseWriter, r *http.Request) {
	feed := feeds.Rss{Feed: h.getFeed(r)}
	rssFeed := feed.RssFeed()
	if rss, err := feeds.ToXML(rssFeed); err != nil {
		panic(err)
	} else {
		h.renderer.Render(w, http.StatusOK, rss)
	}
}

func (h *FeedHandler) getFeed(r *http.Request) *feeds.Feed {
	rootUrl := h.rootUrl.URL(r)
	feed := &feeds.Feed{
		Title:       h.uiSettings.Title,
		Link:        &feeds.Link{Href: rootUrl.String()},
		Description: h.uiSettings.Description,
		Author:      &feeds.Author{Name: h.uiSettings.Author, Email: h.uiSettings.Email},
		Copyright:   "This work is copyright © " + h.uiSettings.Author,
	}

	meta := model.NewLoadedMeta(10)

	var items []*feeds.Item
	if meta.Posts > 0 {
		posts, _ := model.GetPostsPage(0)

		for _, post := range posts {
			rootUrl.Path = post.Path
			item := &feeds.Item{
				Id:          post.ID,
				Title:       post.Title,
				Link:        &feeds.Link{Href: rootUrl.String()},
				Description: post.Body[0:int(math.Min(float64(len(post.Body)), 255))],
				Author:      &feeds.Author{Name: h.uiSettings.Author, Email: h.uiSettings.Email},
				Created:     post.PublishedAt,
			}
			items = append(items, item)
		}
	}

	if len(items) > 0 {
		feed.Created = items[0].Created
	} else {
		feed.Created = time.Now()
	}

	feed.Items = items
	return feed
}
