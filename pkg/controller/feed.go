package controller

import (
	_ "encoding/xml"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/model"
)

type feedController struct {
	uiSettings config.UISetting
	rootUrl    config.RootURL
}

func NewFeedController(uiSettings config.UISetting, rootUrl config.RootURL) *feedController {
	return &feedController{uiSettings, rootUrl}
}

func (c *feedController) AtomHandler(ctx echo.Context) error {
	feed := feeds.Atom{Feed: c.getFeed(ctx.Request().(*standard.Request).Request)}
	atomFeed := feed.AtomFeed()
	if atom, err := feeds.ToXML(atomFeed); err != nil {
		return err
	} else {
		return ctx.Blob(http.StatusOK, echo.MIMEApplicationXMLCharsetUTF8, []byte(atom))
	}
}

func (c *feedController) RssHandler(ctx echo.Context) error {
	feed := feeds.Rss{Feed: c.getFeed(ctx.Request().(*standard.Request).Request)}
	rssFeed := feed.RssFeed()
	if rss, err := feeds.ToXML(rssFeed); err != nil {
		return err
	} else {
		return ctx.Blob(http.StatusOK, echo.MIMEApplicationXMLCharsetUTF8, []byte(rss))
	}
}

func (c *feedController) getFeed(r *http.Request) *feeds.Feed {
	rootUrl := c.rootUrl.URL(r)
	feed := &feeds.Feed{
		Title:       c.uiSettings.Title,
		Link:        &feeds.Link{Href: rootUrl.String()},
		Description: c.uiSettings.Description,
		Author:      &feeds.Author{Name: c.uiSettings.Author, Email: c.uiSettings.Email},
		Copyright:   "This work is copyright © " + c.uiSettings.Author,
	}

	meta := model.NewLoadedMeta(10)

	var items []*feeds.Item
	if meta.Posts > 0 {
		posts, _ := model.GetPostsPage(0)

		for _, post := range posts {
			rootUrl.Path = post.Path
			item := &feeds.Item{
				Id:          post.UUID,
				Title:       post.Title,
				Link:        &feeds.Link{Href: rootUrl.String()},
				Description: post.Body[0:int(math.Min(float64(len(post.Body)), 255))],
				Author:      &feeds.Author{Name: c.uiSettings.Author, Email: c.uiSettings.Email},
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
