package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/boltdb/bolt"

	"../cfg"
)

const (
	BUCKET_INDEX = "index"
	BUCKET_TAGS  = "tags"
)

type Meta struct {
	Posts     int
	PerPage   int
	Pages     int
	Drafts    int
	UpdatedAt time.Time
}

func (meta *Meta) init() {
	meta.Posts = 0
	meta.PerPage = cfg.Int("ui.per_page")
	meta.Pages = 0
	meta.Drafts = 0
	meta.UpdatedAt = time.Now()
}

func (meta *Meta) Load() {
	meta.init()
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte([]byte(BUCKET_INDEX)))
		if bucket != nil {
			jsonMeta := bucket.Get([]byte("meta"))
			json.Unmarshal(jsonMeta, &meta)
		}

		return nil
	})
}

func RebuildIndex() error {
	cfg.Log("Rebuilding index...")

	meta := new(Meta)
	meta.init()

	pathMap := make(map[string]string)
	pageMap := make(map[string][]*Post)
	tagMap := make(map[string][]*Post)
	var draftList []*Post

	if err := db.View(func(tx *bolt.Tx) error {
		bucketPosts := tx.Bucket([]byte(BUCKET_POSTS))
		c := bucketPosts.Cursor()

		currentPage := 0
		pageKey := fmt.Sprintf("page-%d", currentPage)
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			post := new(Post)
			json.Unmarshal(v, &post)

			if post.Draft {
				meta.Drafts++
				draftList = append(draftList, post)
			} else {
				meta.Posts++
				pathMap[post.Path] = post.UUID

				pageMap[pageKey] = append(pageMap[pageKey], post)
				if len(pageMap[pageKey]) >= meta.PerPage {
					currentPage++
					meta.Pages++
					pageKey = fmt.Sprintf("page-%d", currentPage)
				}

				for _, tag := range post.Tags {
					tagMap[strings.ToLower(tag)] = append(tagMap[strings.ToLower(tag)], post)
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if meta.Pages == 0 {
		meta.Pages = 1
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte(BUCKET_INDEX))
		bucketIndex, err := tx.CreateBucketIfNotExists([]byte(BUCKET_INDEX))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte(BUCKET_MAP))
		bucketMap, err := tx.CreateBucketIfNotExists([]byte(BUCKET_MAP))
		if err != nil {
			return err
		}
		tx.DeleteBucket([]byte(BUCKET_TAGS))
		bucketTags, err := tx.CreateBucketIfNotExists([]byte(BUCKET_TAGS))
		if err != nil {
			return err
		}

		jsonMeta, _ := json.Marshal(meta)
		cfg.Log(string(jsonMeta))
		if err := bucketIndex.Put([]byte("meta"), []byte(jsonMeta)); err != nil {
			return err
		}

		for pageKey, postsPage := range pageMap {
			jsonPage, _ := json.Marshal(postsPage)
			if err := bucketIndex.Put([]byte(pageKey), []byte(jsonPage)); err != nil {
				return err
			}
		}

		jsonDrafts, _ := json.Marshal(draftList)
		if err := bucketIndex.Put([]byte("drafts"), []byte(jsonDrafts)); err != nil {
			return err
		}

		for path, uuid := range pathMap {
			if err := bucketMap.Put([]byte(path), []byte(uuid)); err != nil {
				return err
			}
		}

		for tag, postsTag := range tagMap {
			jsonPosts, _ := json.Marshal(postsTag)
			if err := bucketTags.Put([]byte(tag), []byte(jsonPosts)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	cfg.Log("Index rebuilt!")
	return nil
}

func GetPostsPage(page int) ([]Post, error) {
	pageKey := fmt.Sprintf("page-%d", page)
	var posts []Post

	if err := db.View(func(tx *bolt.Tx) error {
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

	if err := db.View(func(tx *bolt.Tx) error {
		bucketIndex := tx.Bucket([]byte(BUCKET_TAGS))
		if bucketIndex == nil {
			panic("Bucket index not found!")
		}

		jsonPosts := bucketIndex.Get([]byte(strings.ToLower(tag)))
		json.Unmarshal(jsonPosts, &posts)

		return nil
	}); err != nil {
		return posts, err
	}

	return posts, nil
}

func GetDraftPosts() ([]Post, error) {
	var posts []Post

	if err := db.View(func(tx *bolt.Tx) error {
		bucketIndex := tx.Bucket([]byte(BUCKET_INDEX))
		if bucketIndex == nil {
			panic("Bucket index not found!")
		}

		jsonPosts := bucketIndex.Get([]byte("drafts"))
		json.Unmarshal(jsonPosts, &posts)

		return nil
	}); err != nil {
		return posts, err
	}

	return posts, nil
}
