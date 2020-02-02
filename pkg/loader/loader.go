package loader

import (
	"errors"
	"net/url"

	"github.com/vgarvardt/rklotz/pkg/storage"
	"go.uber.org/zap"
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

// New returns new loader instance by type
func New(dsn string, logger *zap.Logger) (Loader, error) {
	postsURL, err := url.Parse(dsn)
	if nil != err {
		return nil, err
	}

	switch postsURL.Scheme {
	case schemeFile:
		return NewFileLoader(postsURL.Path, logger)
	}

	return nil, ErrorUnknownLoaderType
}
