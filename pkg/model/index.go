package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

const (
	BUCKET_INDEX = "index"
	BUCKET_TAGS  = "tags"

	INDEX_META = "meta"
)

type Meta struct {
	Posts     int
	PerPage   uint
	Pages     int
	Drafts    int
	UpdatedAt time.Time
}

func (meta *Meta) init(perPage uint) {
	meta.Posts = 0
	meta.PerPage = perPage
	meta.Pages = 1
	meta.Drafts = 0
	meta.UpdatedAt = time.Now()
}

func (meta *Meta) Load(perPage uint) {
	meta.init(perPage)
	DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte([]byte(BUCKET_INDEX)))
		if bucket != nil {
			jsonMeta := bucket.Get([]byte(INDEX_META))
			json.Unmarshal(jsonMeta, &meta)
		}

		return nil
	})
}

func NewLoadedMeta(perPage uint) *Meta {
	meta := new(Meta)
	meta.Load(perPage)
	return meta
}

func GetPostsPage(page int) ([]Post, error) {
	pageKey := fmt.Sprintf("page-%d", page)
	var posts []Post

	if err := DB.View(func(tx *bolt.Tx) error {
		bucketIndex := tx.Bucket([]byte(BUCKET_INDEX))
		if bucketIndex == nil {
			panic("Bucket index not found!")
		}

		jsonPosts := bucketIndex.Get([]byte(pageKey))
		json.Unmarshal(jsonPosts, &posts)

		return nil
	}); err != nil {
		return posts, err
	}

	return posts, nil
}

func GetTagPosts(tag string) ([]Post, error) {
	var posts []Post

	if err := DB.View(func(tx *bolt.Tx) error {
		bucketTags := tx.Bucket([]byte(BUCKET_TAGS))
		if bucketTags == nil {
			panic("Bucket index not found!")
		}

		jsonPosts := bucketTags.Get([]byte(strings.ToLower(tag)))
		json.Unmarshal(jsonPosts, &posts)

		return nil
	}); err != nil {
		return posts, err
	}

	return posts, nil
}

func MustGetTagPosts(tag string) []Post {
	if posts, err := GetTagPosts(tag); err != nil {
		panic(err)
	} else {
		return posts
	}
}
