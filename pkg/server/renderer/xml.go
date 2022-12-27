package renderer

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/feeds"

	"github.com/vgarvardt/rklotz/pkg/model"
)

const (
	charsetUTF8                   = "charset=utf-8"
	mimeApplicationXML            = "application/xml"
	mimeApplicationXMLCharsetUTF8 = mimeApplicationXML + "; " + charsetUTF8
	headerContentType             = "Content-Type"
)

// Feed implements Renderer for XML content
type Feed struct {
	cfgUI      UIConfig
	cfgRootURL RootURLConfig
}

// NewFeed creates new Feed instance
func NewFeed(cfgUI UIConfig, cfgRootURL RootURLConfig) *Feed {
	return &Feed{cfgUI, cfgRootURL}
}

// Render renders the data with response code to a HTTP response writer
func (r *Feed) Render(w http.ResponseWriter, code int, data *Data) {
	data.m.RLock()
	defer data.m.RUnlock()

	posts, ok := data.data["posts"].([]*model.Post)
	if !ok {
		panic(errors.New("unknown posts list type for feed"))
	}

	rootURL := r.cfgRootURL.URL(data.r)
	feed := &feeds.Feed{
		Title:       r.cfgUI.Title,
		Link:        &feeds.Link{Href: rootURL.String()},
		Description: r.cfgUI.Description,
		Author:      &feeds.Author{Name: r.cfgUI.Author, Email: r.cfgUI.Email},
		Copyright:   "This work is copyright © " + r.cfgUI.Author,
		Items:       make([]*feeds.Item, len(posts)),
	}

	for i, post := range posts {
		rootURL.Path = post.Path
		feed.Items[i] = &feeds.Item{
			Id:          post.ID,
			Title:       post.Title,
			Link:        &feeds.Link{Href: rootURL.String()},
			Description: post.Body[0:int(math.Min(float64(len(post.Body)), 255))],
			Author:      &feeds.Author{Name: r.cfgUI.Author, Email: r.cfgUI.Email},
			Created:     post.PublishedAt,
		}
	}

	if len(feed.Items) > 0 {
		feed.Created = feed.Items[0].Created
	} else {
		feed.Created = time.Now()
	}

	var xmlFeed feeds.XmlFeed
	switch data.template {
	case "atom":
		xmlFeed = (&feeds.Atom{Feed: feed}).AtomFeed()
	case "rss":
		xmlFeed = (&feeds.Rss{Feed: feed}).RssFeed()
	default:
		panic(fmt.Errorf("unknown feed type: %s", data.template))
	}

	xmlData, err := feeds.ToXML(xmlFeed)
	if err != nil {
		panic(err)
	}

	w.Header().Set(headerContentType, mimeApplicationXMLCharsetUTF8)
	w.WriteHeader(code)
	if _, err := w.Write([]byte(xmlData)); err != nil {
		panic(err)
	}
}

// Error renders error to a response writer
func (r *Feed) Error(_ *http.Request, w http.ResponseWriter, code int, err error) {
	http.Error(w, err.Error(), code)
}
