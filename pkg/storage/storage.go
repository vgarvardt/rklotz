package storage

import (
	"errors"
	"net/url"

	"github.com/vgarvardt/rklotz/pkg/model"
)

const (
	schemeBoldDB = "boltdb"
	schemeMemory = "memory"
)

var (
	ErrorUnknownStorageType = errors.New("Unknown storage type")
	ErrorNotFound           = errors.New("Record not found")
)

type Storage interface {
	Save(post *model.Post) error
	Finalize() error
	FindByPath(path string) (*model.Post, error)
	ListAll(page int) ([]*model.Post, error)
	ListTag(tag string, page int) ([]*model.Post, error)
	Close() error

	Meta() *model.Meta
	TagMeta(tag string) *model.Meta
}

func NewStorage(dsn string, postsPerPage int) (Storage, error) {
	dsnURL, err := url.Parse(dsn)
	if nil != err {
		return nil, err
	}

	switch dsnURL.Scheme {
	case schemeBoldDB:
		return NewBoltDBStorage(dsnURL.Path, postsPerPage)
	case schemeMemory:
		return NewMemoryStorage(postsPerPage)
	}

	return nil, ErrorUnknownStorageType
}
