package model

import (
	"time"
	"strings"
	"net/http"
	"encoding/json"

	"github.com/russross/blackfriday"
	"github.com/gorilla/schema"
	"github.com/boltdb/bolt"
	"code.google.com/p/go-uuid/uuid"
)

const (
	BUCKET_POSTS = "posts"
	BUCKET_MAP = "path_map"
)

type Format struct {
	Name    string
	Title   string
	Handler func(string) string
}

func GetAvailableFormats() []Format {
	return []Format{
		Format{
			Name:  "md",
			Title: "MarkDown",
			Handler: (func(input string) string {
				return string(blackfriday.MarkdownCommon([]byte(input)))
			}),
		},
	}
}

func bindFormToStruct(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(obj, req.PostForm); err != nil {
		return err
	}

	return nil
}

type Post struct {
	UUID		string
	Path		string
	Title		string
	Body		string
	Format		string
	HTML		string		`schema:"-"`
	Tags		[]string
	Draft		bool		`schema:"-"`
	CreatedAt	time.Time	`schema:"-"`
	UpdatedAt	time.Time	`schema:"-"`
}

func (post *Post) Bind(req *http.Request) error {
	if err := bindFormToStruct(req, post); err != nil {
		return err
	}
	post.Path = strings.Trim(post.Path, "/")
	if len(post.Tags) > 0 {
		post.Tags = strings.Split(post.Tags[0], ",")
	}

	return nil
}

func (post *Post) ReFormat() string {
	post.HTML = post.Body

	formats := GetAvailableFormats()
	for i := 0; i < len(formats); i++ {
		if formats[i].Name == post.Format {
			post.HTML = formats[i].Handler(post.Body)
		}
	}
	return post.HTML
}

func (post *Post) Save(draft bool) error {
	post.Draft = draft
	post.UpdatedAt = time.Now()
	post.ReFormat()

	if err := db.Update(func(tx *bolt.Tx) error {
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

	go RebuildIndex()
	return nil
}

func (post *Post) Load(uuid string) error {
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_POSTS))
		if bucket == nil {
			panic("Bucket posts not found!")
		}

		jsonPost := bucket.Get([]byte(uuid))
		json.Unmarshal(jsonPost, &post)

		return nil
	})
}

func (post *Post) LoadByPath(path string) error {
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_MAP))
		if bucket == nil {
			panic("Bucket path_map not found!")
		}

		uuid := bucket.Get([]byte(path))
		return post.Load(string(uuid))
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
