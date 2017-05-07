package repository

import (
	"os"
	"path/filepath"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/vgarvardt/rklotz/pkg/model"
)

type BoltDBStorage struct {
	db           *storm.DB
	path         string
	postsPerPage uint
}

func NewBoltDBStorage(path string, postsPerPage uint) (*BoltDBStorage, error) {
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

	// Initialize buckets and indexes before saving an object
	// Useful when starting your application
	err = instance.db.Init(&model.Post{})

	return instance, err
}

func (s *BoltDBStorage) Save(post *model.Post) error {
	return s.db.Save(post)
}

func (s *BoltDBStorage) Reindex(postsPerPage uint) error {
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

func (s *BoltDBStorage) ListAll(page uint) ([]*model.Post, error) {
	var posts []*model.Post
	offset := int(page * s.postsPerPage)
	err := s.db.AllByIndex("PublishedAt", &posts, storm.Limit(int(s.postsPerPage)), storm.Skip(offset), storm.Reverse())
	return posts, err
}

func (s *BoltDBStorage) ListTag(tag string, page uint) ([]*model.Post, error) {
	// TODO: does not work, need to implement tag-posts map in separate bucket
	var posts []*model.Post

	offset := int(page * s.postsPerPage)
	query := s.db.Select(q.Eq("Tags", tag)).Limit(int(s.postsPerPage)).Skip(offset).OrderBy("PublishedAt").Reverse()
	err := query.Find(&posts)

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
