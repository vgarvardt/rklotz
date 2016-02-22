package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/boltdb/bolt"

	"github.com/vgarvardt/rklotz/cfg"
)

const (
	BUCKET_INDEX = "index"
	BUCKET_TAGS = "tags"
	BUCKET_TAG_MAP = "tag_map"

	INDEX_META = "meta"
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
	meta.Pages = 1
	meta.Drafts = 0
	meta.UpdatedAt = time.Now()
}

func (meta *Meta) Load() {
	meta.init()
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte([]byte(BUCKET_INDEX)))
		if bucket != nil {
			jsonMeta := bucket.Get([]byte(INDEX_META))
			json.Unmarshal(jsonMeta, &meta)
		}

		return nil
	})
}

func NewLoadedMeta() *Meta {
	meta := new(Meta)
	meta.Load()
	return meta
}

func RebuildIndex() error {
	cfg.Log("Rebuilding index...")

	var publishedStamps []int
	postsMap := make(map[int]string)

	if err := db.View(func(tx *bolt.Tx) error {
		bucketPosts := tx.Bucket([]byte(BUCKET_POSTS))
		if bucketPosts != nil {
			c := bucketPosts.Cursor()
			for k, v := c.Last(); k != nil; k, v = c.Prev() {
				post := new(Post)
				json.Unmarshal(v, &post)

				unixPublished := int(post.PublishedAt.Unix())
				if _, exists := postsMap[unixPublished]; exists {
					return errors.New(fmt.Sprintf("Duplicate published at date %v", post.PublishedAt))
				}
				publishedStamps = append(publishedStamps, unixPublished)
				postsMap[unixPublished] = post.UUID
			}
		}

		return nil
	}); err != nil {
		return err
	}

	sort.Sort(sort.Reverse(sort.IntSlice(publishedStamps)))

	if err := db.Update(func(tx *bolt.Tx) error {
		meta := new(Meta)
		meta.init()

		pathMap := make(map[string]string)
		pageMap := make(map[string][]*Post)
		tags := make(map[string][]*Post)
		tagMap := make(map[string]string)

		currentPage := 0
		pageKey := fmt.Sprintf("page-%d", currentPage)

		bucketPosts := tx.Bucket([]byte(BUCKET_POSTS))
		if bucketPosts != nil {
			for _, unixPublished := range publishedStamps {
				post := new(Post)
				jsonPost := bucketPosts.Get([]byte(postsMap[unixPublished]))
				json.Unmarshal(jsonPost, &post)

				post.ReFormat()
				jsonPostReformatted, _ := json.Marshal(post)
				if err := bucketPosts.Put([]byte(post.UUID), []byte(jsonPostReformatted)); err != nil {
					return err
				}

				if post.Draft {
					meta.Drafts++
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
						_tag := strings.ToLower(tag)
						tags[_tag] = append(tags[_tag], post)
						tagMap[_tag] = tag
					}
				}
			}

			// fix situation, when last page is empty
			if len(pageMap[pageKey]) < 1 && meta.Pages > 1 {
				meta.Pages--
			}
		}

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
		tx.DeleteBucket([]byte(BUCKET_TAG_MAP))
		bucketTagMap, err := tx.CreateBucketIfNotExists([]byte(BUCKET_TAG_MAP))
		if err != nil {
			return err
		}

		jsonMeta, _ := json.Marshal(meta)
		cfg.Log(string(jsonMeta))
		if err := bucketIndex.Put([]byte(INDEX_META), []byte(jsonMeta)); err != nil {
			return err
		}

		for pageKey, postsPage := range pageMap {
			jsonPage, _ := json.Marshal(postsPage)
			if err := bucketIndex.Put([]byte(pageKey), []byte(jsonPage)); err != nil {
				return err
			}
		}

		for path, uuid := range pathMap {
			if err := bucketMap.Put([]byte(path), []byte(uuid)); err != nil {
				return err
			}
		}

		for tag, postsTag := range tags {
			jsonPosts, _ := json.Marshal(postsTag)
			if err := bucketTags.Put([]byte(tag), []byte(jsonPosts)); err != nil {
				return err
			}
		}

		for _tag, tag := range tagMap {
			if err := bucketTagMap.Put([]byte(_tag), []byte(tag)); err != nil {
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

func AutoCompleteTags(q string) []string {
	var tags []string
	db.View(func(tx *bolt.Tx) error {
		bucketTags := tx.Bucket([]byte(BUCKET_TAGS))
		if bucketTags == nil {
			return nil
		}
		bucketTagMap := tx.Bucket([]byte(BUCKET_TAG_MAP))
		if bucketTagMap == nil {
			return nil
		}
		c := bucketTags.Cursor()

		prefix := []byte(q)
		for k, _ := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			tags = append(tags, string(bucketTagMap.Get(k)))
		}

		return nil
	})

	return tags
}

func getIndexPosts(draft bool) ([]*Post, error) {
	var posts []*Post
	var result []*Post
	var createdStamps []int
	postsMap := make(map[int]int)

	if err := db.View(func(tx *bolt.Tx) error {
		bucketPosts := tx.Bucket([]byte(BUCKET_POSTS))
		if bucketPosts != nil {
			c := bucketPosts.Cursor()
			for k, v := c.Last(); k != nil; k, v = c.Prev() {
				post := new(Post)
				json.Unmarshal(v, &post)
				if post.Draft == draft {
					posts = append(posts, post)
					createdStamps = append(createdStamps, int(post.CreatedAt.Unix()))
					postsMap[createdStamps[len(createdStamps) - 1]] = len(posts) - 1
				}
			}
		}
		return nil
	}); err != nil {
		return posts, err
	}

	sort.Sort(sort.Reverse(sort.IntSlice(createdStamps)))
	for _, stamp := range createdStamps {
		result = append(result, posts[postsMap[stamp]])
	}

	return result, nil
}

func GetDraftPosts() ([]*Post, error) {
	return getIndexPosts(true)
}

func GetPublishedPosts() ([]*Post, error) {
	return getIndexPosts(false)
}
