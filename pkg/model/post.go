package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/labstack/echo"
	"github.com/pborman/uuid"
	"github.com/russross/blackfriday"
	"github.com/vgarvardt/rklotz/pkg/svc"
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
	UUID        string    `form:"uuid"`
	Path        string    `form:"path"`
	Title       string    `form:"title"`
	Body        string    `form:"body"`
	Format      string    `form:"format"`
	Tags        []string  `form:"tags"`
	HTML        string    `form:"-"`
	Draft       bool      `form:"-"`
	CreatedAt   time.Time `form:"-"`
	UpdatedAt   time.Time `form:"-"`
	PublishedAt time.Time `form:"-"`
}

func (post *Post) Bind(ctx echo.Context) error {
	var err error
	if err = ctx.Bind(post); err != nil {
		return err
	}

	post.Path = strings.Trim(post.Path, "/")
	if post.PublishedAt, err = time.Parse(time.RFC3339, ctx.FormValue("published_at")); err != nil {
		post.PublishedAt = time.Now()
	}

	return nil
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

func (post *Post) Save(draft bool) error {
	post.Draft = draft
	post.UpdatedAt = time.Now()
	post.ReFormat()

	if err := DB.Update(func(tx *bolt.Tx) error {
		bucketPosts, err := tx.CreateBucketIfNotExists([]byte(BUCKET_POSTS))
		if err != nil {
			return err
		}

		if len(post.UUID) < 1 {
			post.UUID = uuid.New()
			post.CreatedAt = time.Now()
		}

		jsonPost, _ := json.Marshal(post)
		if err := bucketPosts.Put([]byte(post.UUID), []byte(jsonPost)); err != nil {
			return err
		}

		bucketMap, err := tx.CreateBucketIfNotExists([]byte(BUCKET_MAP))
		if err != nil {
			return err
		}

		if err := bucketMap.Put([]byte(post.Path), []byte(post.UUID)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	logger := svc.Container.MustGet(svc.DI_LOGGER).(*log.Logger)
	logger.WithField("UUID", post.UUID).Info("Saved post")
	go RebuildIndex()
	return nil
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

func UpdatePostField(uuid, field, value string) error {
	post := new(Post)
	if err := post.Load(uuid); err != nil {
		return err
	}

	if len(post.UUID) < 1 {
		return errors.New("Post could not be loaded")
	}

	switch {
	case field == "PublishedAt":
		var t time.Time
		var err error
		if t, err = time.Parse(time.RFC3339, value); err != nil {
			return errors.New(fmt.Sprintf("Invalid value '%s' for '%s': %v (must be in '%s' format)", value, field, err, time.RFC3339))
		}
		post.PublishedAt = t
		break
	default:
		return errors.New(fmt.Sprintf("Unknown field '%s'", field))
	}

	return post.Save(post.Draft)
}

func (post *Post) Delete() error {
	if err := DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_POSTS))
		if bucket == nil {
			panic(fmt.Sprintf("Bucket %s not found!", BUCKET_POSTS))
		}

		if err := bucket.Delete([]byte(post.UUID)); err != nil {
			panic(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	logger := svc.Container.MustGet(svc.DI_LOGGER).(*log.Logger)

	logger.WithField("UUID", post.UUID).Info("Removed post")
	go RebuildIndex()
	return nil
}
