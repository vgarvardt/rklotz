package loader

import (
	"errors"
	"net/url"

	"github.com/vgarvardt/rklotz/pkg/storage"
)

const (
	schemeFile = "file"
)

var (
	// ErrorUnknownLoaderType is the error returned when trying to instantiate a loader of unknown type
	ErrorUnknownLoaderType = errors.New("Unknown loader type")
)

// Loader is the interface for posts loader
type Loader interface {
	// Load loads posts and saves them one by one in the storage
	Load(storage storage.Storage) error
}

// NewLoader returns new loader instance by type
func NewLoader(dsn string) (Loader, error) {
	postsURL, err := url.Parse(dsn)
	if nil != err {
		return nil, err
	}

	switch postsURL.Scheme {
	case schemeFile:
		return NewFileLoader(postsURL.Path)
	}

	return nil, ErrorUnknownLoaderType
}
