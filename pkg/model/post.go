package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/russross/blackfriday"
)

const (
	BUCKET_POSTS = "posts"
	BUCKET_MAP   = "path_map"
)

type Format struct {
	Name    string
	Title   string
	Handler func(string) string
}

func GetAvailableFormats() []Format {
	return []Format{
		{
			Name:  "md",
			Title: "MarkDown",
			Handler: func(input string) string {
				return string(blackfriday.MarkdownCommon([]byte(input)))
			},
		},
	}
}

type Post struct {
	UUID        string
	Path        string
	Title       string
	Body        string
	Format      string
	Tags        []string
	HTML        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt time.Time
}

func (post *Post) ReFormat() string {
	post.HTML = post.Body

	formats := GetAvailableFormats()
	for i := 0; i < len(formats); i++ {
		if formats[i].Name == post.Format {
			post.HTML = formats[i].Handler(post.Body)
			break
		}
	}

	// open all links in new tab
	post.HTML = strings.Replace(post.HTML, `<a href=`, `<a target="_blank" href=`, -1)
	// fix code class to make highlight.js work
	re := regexp.MustCompile(`<code class="language-(\w+)">`)
	post.HTML = re.ReplaceAllString(post.HTML, "<code class=\"$1\">")

	return post.HTML
}

func (post *Post) Load(uuid string) error {
	return DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_POSTS))
		if bucket == nil {
			panic(fmt.Sprintf("Bucket %s not found!", BUCKET_POSTS))
		}

		jsonPost := bucket.Get([]byte(uuid))
		json.Unmarshal(jsonPost, &post)

		return nil
	})
}

func (post *Post) LoadByPath(path string) error {
	return DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_MAP))
		if bucket == nil {
			panic(fmt.Sprintf("Bucket %s not found!", BUCKET_MAP))
		}

		postUUID := bucket.Get([]byte(path))
		return post.Load(string(postUUID))
	})
}

func (post *Post) Validate() map[string]string {
	var err map[string]string = make(map[string]string)

	if strings.TrimSpace(post.Title) == "" {
		err["Title"] = "Title can not be empty"
	}

	if strings.TrimSpace(post.Body) == "" {
		err["Body"] = "Body can not be empty"
	}

	if strings.TrimSpace(post.Format) == "" {
		err["Format"] = "Format can not be empty"
	}

	if strings.TrimSpace(post.Path) == "" {
		err["Path"] = "Path can not be empty"
	}

	return err
}
