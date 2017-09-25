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
	// ErrorUnknownStorageType is the error returned when trying to instantiate a storage of unknown type
	ErrorUnknownStorageType = errors.New("Unknown storage type")
	// ErrorNotFound is the error returned when trying to find a post by non-existent path
	ErrorNotFound = errors.New("Record not found")
)

// Storage is the interface for posts storage
type Storage interface {
	// Save persists new post in the storage
	Save(post *model.Post) error
	// Finalize is called after all posts are persisted in the storage
	Finalize() error
	// FindByPath searches for a post by path
	FindByPath(path string) (*model.Post, error)
	// ListAll returns ordered by date posts page
	ListAll(page int) ([]*model.Post, error)
	// ListTag returns ordered by date posts page for a tag
	ListTag(tag string, page int) ([]*model.Post, error)
	// Close closes the storage and frees all resources
	Close() error

	// Meta returns metadata for all persisted posts
	Meta() *model.Meta
	// TagMeta returns metadata for all persisted posts for a tag
	TagMeta(tag string) *model.Meta
}

// NewStorage returns new storage instance by type
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
