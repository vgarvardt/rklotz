package model

import (
	"time"
	"net/http"
	"encoding/json"

	"github.com/russross/blackfriday"
	"github.com/gorilla/schema"
	"github.com/boltdb/bolt"
	"code.google.com/p/go-uuid/uuid"
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
				return string(blackfriday.MarkdownBasic([]byte(input)))
			}),
		},
	}
}

func bindFormToStruct(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	decoder := schema.NewDecoder()
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
	return bindFormToStruct(req, post);
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

	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("posts"))
		if err != nil {
			return err
		}

		if len(post.UUID) < 1 {
			post.UUID = uuid.New()
			post.CreatedAt = time.Now()
		}

		jsonPost, _ := json.Marshal(post)
		if err := bucket.Put([]byte(post.UUID), []byte(jsonPost)); err != nil {
			return err
		}

		return nil
	})
}

func (post *Post) Load(uuid string) error {
	return db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("posts"))
		if bucket == nil {
			panic("Bucket posts not found!")
		}

		jsonPost := bucket.Get([]byte(uuid))
		json.Unmarshal([]byte(jsonPost), &post)

		return nil
	})
}
