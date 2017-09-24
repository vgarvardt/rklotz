package storage

import (
	"os"
	"path/filepath"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/vgarvardt/rklotz/pkg/model"
)

const tagsNode = "__rklotz_tags"

type BoltDBStorage struct {
	db   *storm.DB
	path string
	tags storm.Node

	postsCount   int
	postsPerPage int
}

func NewBoltDBStorage(path string, postsPerPage int) (*BoltDBStorage, error) {
	var err error

	instance := &BoltDBStorage{path: path, postsPerPage: postsPerPage}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if nil != err {
		return nil, err
	}

	err = instance.remove()
	if nil != err {
		return nil, err
	}

	instance.db, err = storm.Open(path)
	if nil != err {
		return nil, err
	}

	instance.tags = instance.db.From(tagsNode)

	// Initialize buckets and indexes before saving an object
	// Useful when starting your application
	err = instance.db.Init(&model.Post{})
	if nil != err {
		return nil, err
	}
	err = instance.tags.Init(&model.Tag{})
	if nil != err {
		return nil, err
	}

	return instance, nil
}

func (s *BoltDBStorage) Save(post *model.Post) error {
	err := s.db.Save(post)
	if nil != err {
		return err
	}

	for i := range post.Tags {
		var tag model.Tag
		err := s.tags.One("Tag", post.Tags[i], &tag)
		if err == storm.ErrNotFound {
			tag = model.Tag{Tag: post.Tags[i], Paths: []string{}}
		} else if nil != err {
			return err
		}

		tag.Paths = append(tag.Paths, post.Path)
		err = s.tags.Save(&tag)
		if nil != err {
			return err
		}
	}

	s.postsCount++
	return nil
}

func (s *BoltDBStorage) Finalize() error {
	return nil
}

func (s *BoltDBStorage) FindByPath(path string) (*model.Post, error) {
	var post model.Post
	err := s.db.One("Path", path, &post)

	if err == storm.ErrNotFound {
		return nil, ErrorNotFound
	}

	return &post, err
}

func (s *BoltDBStorage) ListAll(page int) ([]*model.Post, error) {
	var posts []*model.Post
	offset := page * s.postsPerPage
	err := s.db.AllByIndex("PublishedAt", &posts, storm.Limit(s.postsPerPage), storm.Skip(offset), storm.Reverse())
	return posts, err
}

func (s *BoltDBStorage) ListTag(tag string, page int) ([]*model.Post, error) {
	var tagObject model.Tag

	err := s.tags.One("Tag", tag, &tagObject)
	if err == storm.ErrNotFound {
		return nil, ErrorNotFound
	}

	var posts []*model.Post
	offset := page * s.postsPerPage
	query := s.db.Select(q.In("Path", tagObject.Paths)).Limit(s.postsPerPage).Skip(offset).OrderBy("PublishedAt").Reverse()

	err = query.Find(&posts)
	if err == storm.ErrNotFound {
		return []*model.Post{}, nil
	}

	if nil != err {
		return nil, err
	}

	return posts, nil
}

func (s *BoltDBStorage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return s.remove()
}

func (s *BoltDBStorage) Meta() *model.Meta {
	return model.NewMeta(s.postsCount, s.postsPerPage)
}

func (s *BoltDBStorage) TagMeta(tag string) *model.Meta {
	var tagObject model.Tag

	err := s.tags.One("Tag", tag, &tagObject)
	if err != nil {
		return &model.Meta{}
	}

	return model.NewMeta(len(tagObject.Paths), s.postsPerPage)
}

func (s *BoltDBStorage) remove() error {
	_, err := os.Stat(s.path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err = os.Remove(s.path)
		if err != nil {
			return err
		}
	}

	return nil
}
