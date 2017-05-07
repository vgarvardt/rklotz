package repository

import (
	"errors"
	"net/url"

	"github.com/vgarvardt/rklotz/pkg/model"
)

const (
	schemeBoldDB = "boltdb"
)

var (
	ErrorUnknownStorageType = errors.New("Uknonwn storage type")
	ErrorNotFound           = errors.New("Post not found")
)

type Storage interface {
	Save(post *model.Post) error
	Reindex(postsPerPage uint) error
	FindByPath(path string) (*model.Post, error)
	ListAll(page uint) ([]*model.Post, error)
	ListTag(tag string, page uint) ([]*model.Post, error)
	Close() error
}

func NewStorage(dsn string, postsPerPage uint) (Storage, error) {
	storageURL, err := url.Parse(dsn)
	if nil != err {
		return nil, err
	}

	switch storageURL.Scheme {
	case schemeBoldDB:
		return NewBoltDBStorage(storageURL.Path, postsPerPage)
	}

	return nil, ErrorUnknownStorageType
}
