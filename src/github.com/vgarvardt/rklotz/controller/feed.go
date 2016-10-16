package controller

import (
	_ "encoding/xml"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/gorilla/feeds"

	"github.com/vgarvardt/rklotz/app"
	"github.com/vgarvardt/rklotz/model"
	"github.com/vgarvardt/rklotz/svc"
	"github.com/vgarvardt/rklotz/utils"
)

func AtomController(ctx echo.Context) error {
	feed := feeds.Atom{Feed: getFeed(ctx.Request().(*standard.Request).Request)}
	atomFeed := feed.AtomFeed()
	if atom, err := feeds.ToXML(atomFeed); err != nil {
		return err
	} else {
		return ctx.Blob(http.StatusOK, echo.MIMEApplicationXMLCharsetUTF8, []byte(atom))
	}
}

func RssController(ctx echo.Context) error {
	feed := feeds.Rss{Feed: getFeed(ctx.Request().(*standard.Request).Request)}
	rssFeed := feed.RssFeed()
	if rss, err := feeds.ToXML(rssFeed); err != nil {
		return err
	} else {
		return ctx.Blob(http.StatusOK, echo.MIMEApplicationXMLCharsetUTF8, []byte(rss))
	}
}

func getFeed(r *http.Request) *feeds.Feed {
	config := svc.Container.MustGet(svc.DI_CONFIG).(svc.Config)

	rootUrl := app.RootUrl(r)
	feed := &feeds.Feed{
		Title:       config.String("ui.title"),
		Link:        &feeds.Link{Href: rootUrl.String()},
		Description: config.String("ui.description"),
		Author:      &feeds.Author{Name: config.String("ui.author"), Email: config.String("ui.email")},
		Copyright:   "This work is copyright © " + config.String("ui.author"),
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
				Description: post.Body[0:utils.Min(len(post.Body), 255)],
				Author:      &feeds.Author{Name: config.String("ui.author"), Email: config.String("ui.email")},
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
