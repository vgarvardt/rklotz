package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/vgarvardt/rklotz/pkg/model"
)

var _ Storage = (*BadgerStorage)(nil)

type BadgerStorage struct {
	db *badger.DB

	path         string
	postsCount   int
	postsPerPage int
}

func NewBadgerStorage(path string, postsPerPage int) (*BadgerStorage, error) {
	var err error

	instance := &BadgerStorage{
		path:         path,
		postsPerPage: postsPerPage,
	}

	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}

	instance.db = db

	return instance, nil
}

func (b *BadgerStorage) Close() error {
	if err := b.db.Close(); err != nil {
		return err
	}
	return b.remove()
}

func (b *BadgerStorage) Save(post *model.Post) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		postData, err := json.Marshal(post)
		if err != nil {
			return err
		}

		// save key with batch index for pagination
		batchCount := 0
		if b.postsCount != 0 {
			if b.postsCount%b.postsPerPage == 0 {
				batchCount++
			}
		}

		key := fmt.Sprintf("pat:%d:%s", batchCount, post.Path) // iterate by page

		err = txn.Set([]byte(key), []byte(post.Path))
		if err != nil {
			return err
		}

		err = txn.Set([]byte(post.Path), postData) // find by path
		if err != nil {
			return err
		}

		for i := 0; i < len(post.Tags); i++ {
			tag := &model.Tag{}
			tagKey := []byte("tag_" + post.Tags[i])

			item, err := txn.Get(tagKey)
			if errors.Is(err, badger.ErrKeyNotFound) {
				tag = &model.Tag{Tag: post.Tags[i], Paths: []string{post.Path}}
			} else if err != nil {
				return err
			} else {
				err = item.Value(func(val []byte) error {
					return json.Unmarshal(val, tag)
				})
				if err != nil {
					return err
				}
				tag.Paths = append(tag.Paths, post.Path)
			}

			tagData, err := json.Marshal(tag)
			if err != nil {
				return err
			}

			err = txn.Set(tagKey, tagData)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	b.postsCount++

	return nil
}

func (b *BadgerStorage) Finalize() error {
	return nil
}

func (b *BadgerStorage) FindByPath(path string) (*model.Post, error) {
	existedPost := &model.Post{}

	err := b.db.View(func(txn *badger.Txn) error {
		post, err := txn.Get([]byte(path))
		if err != nil {
			return err
		}

		return post.Value(func(val []byte) error {
			return json.Unmarshal(val, existedPost)
		})
	})
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, ErrorNotFound
		}

		return nil, err
	}

	return existedPost, nil
}

func (b *BadgerStorage) ListAll(page int) ([]*model.Post, error) {
	offset := 0
	if page > 1 {
		offset = (page - 1) * b.postsPerPage
	}

	var posts []*model.Post

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(fmt.Sprintf("pat:%d", page))
		opts.PrefetchSize = b.postsPerPage
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		count := 0
		for it.Seek([]byte{}); it.Valid(); it.Next() {
			path := ""
			err := it.Item().Value(func(val []byte) error {
				path = string(val)
				return nil
			})

			if path == "" {
				continue
			}

			raw, err2 := txn.Get([]byte(path))
			if err2 != nil {
				return err2
			}

			post := &model.Post{}
			err = raw.Value(func(val []byte) error {
				return json.Unmarshal(val, post)
			})
			if err != nil {
				continue
			}

			if post.ID == "" {
				continue
			}

			if count >= offset && count < offset+b.postsPerPage {
				posts = append(posts, post)
			}

			count++
			if count >= offset+b.postsPerPage {
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (b *BadgerStorage) ListTag(tag string, page int) ([]*model.Post, error) {
	if page < 1 {
		page = 1
	}

	offset := (page - 1) * b.postsPerPage

	var posts []*model.Post

	err := b.db.View(func(txn *badger.Txn) error {
		tagKey := []byte("tag_" + tag)
		item, err := txn.Get(tagKey)
		if err != nil {
			return err
		}

		tag := &model.Tag{}
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, tag)
		})
		if err != nil {
			return err
		}

		// Implement pagination
		if offset >= len(tag.Paths) {
			return nil // No posts to return
		}
		end := offset + b.postsPerPage
		if end > len(tag.Paths) {
			end = len(tag.Paths)
		}

		for _, postPath := range tag.Paths[offset:end] {
			postItem, err := txn.Get([]byte(postPath))
			if err != nil {
				return err
			}

			post := &model.Post{}
			err = postItem.Value(func(val []byte) error {
				return json.Unmarshal(val, post)
			})
			if err != nil {
				return err
			}

			posts = append(posts, post)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (b *BadgerStorage) Meta() *model.Meta {
	return model.NewMeta(b.postsCount, b.postsPerPage)
}

func (b *BadgerStorage) TagMeta(tag string) *model.Meta {
	var meta *model.Meta

	err := b.db.View(func(txn *badger.Txn) error {
		tagKey := []byte("tag_" + tag)
		item, err := txn.Get(tagKey)
		if err != nil {
			return err
		}

		tag := &model.Tag{}
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, tag)
		})
		if err != nil {
			return err
		}

		meta = model.NewMeta(len(tag.Paths), b.postsPerPage)

		return nil
	})

	if err != nil {
		return nil
	}

	return meta
}

func (b *BadgerStorage) remove() error {
	_, err := os.Stat(b.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}

	return os.RemoveAll(b.path)
}
