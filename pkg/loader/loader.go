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
	ErrorUnknownLoaderType = errors.New("Unknown loader type")
)

type Loader interface {
	Load(storage storage.Storage) error
}

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
