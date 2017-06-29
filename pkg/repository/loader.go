package repository

import (
	"errors"
	"net/url"
)

const (
	schemeFile = "file"
)

var (
	ErrorUnknownLoaderType = errors.New("Uknonwn loader type")
)

type Loader interface {
	Load(storage Storage) error
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
