package loader

import (
	"errors"
	"log/slog"
	"net/url"

	"github.com/vgarvardt/rklotz/pkg/formatter"
	"github.com/vgarvardt/rklotz/pkg/storage"
)

const (
	schemeFile = "file"
)

var (
	// ErrorUnknownLoaderType is the error returned when trying to instantiate a loader of unknown type
	ErrorUnknownLoaderType = errors.New("unknown loader type")
)

// Loader is the interface for posts loader
type Loader interface {
	// Load loads posts and saves them one by one in the storage
	Load(s storage.Storage) error
}

// New returns new loader instance by type
func New(dsn string, logger *slog.Logger) (Loader, error) {
	postsURL, err := url.Parse(dsn)
	if nil != err {
		return nil, err
	}

	f := formatter.New()

	if postsURL.Scheme == schemeFile {
		return NewFileLoader(postsURL.Path, f, logger)
	}

	return nil, ErrorUnknownLoaderType
}
