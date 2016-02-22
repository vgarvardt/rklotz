package controller

import (
	_ "encoding/xml"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"

	"github.com/vgarvardt/rklotz/cfg"
	"github.com/vgarvardt/rklotz/model"
)

func AtomController(c *gin.Context) {
	feed := feeds.Atom{getFeed(c)}
	atomFeed := feed.AtomFeed()
	if atom, err := feeds.ToXML(atomFeed); err != nil {
		c.Abort()
	} else {
		// tried c.XML(...) but browser detect it as XML,
		// with this custom code browser detects it as feed
		c.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(atom))
	}
}

func RssController(c *gin.Context) {
	feed := feeds.Rss{getFeed(c)}
	rssFeed := feed.RssFeed()
	if rss, err := feeds.ToXML(rssFeed); err != nil {
		c.Abort()
	} else {
		c.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(rss))
	}
}

func getFeed(c *gin.Context) *feeds.Feed {
	rootUrl := cfg.GetRootUrl(c.Request)
	feed := &feeds.Feed{
		Title:       cfg.String("ui.title"),
		Link:        &feeds.Link{Href: rootUrl.String()},
		Description: cfg.String("ui.description"),
		Author:      &feeds.Author{cfg.String("ui.author"), cfg.String("ui.email")},
		Copyright:   "This work is copyright © " + cfg.String("ui.author"),
	}

	meta := model.NewLoadedMeta()

	var items []*feeds.Item
	if meta.Posts > 0 {
		posts, _ := model.GetPostsPage(0)

		for _, post := range posts {
			rootUrl.Path = post.Path
			item := &feeds.Item{
				Id:          post.UUID,
				Title:       post.Title,
				Link:        &feeds.Link{Href: rootUrl.String()},
				Description: post.Body[0:cfg.Min(len(post.Body), 255)],
				Author:      &feeds.Author{cfg.String("ui.author"), cfg.String("ui.email")},
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
