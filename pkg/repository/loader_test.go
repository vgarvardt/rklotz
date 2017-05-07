package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoader(t *testing.T) {
	fileLoader, err := NewLoader("file:///path/to/posts")
	assert.NoError(t, err)
	assert.IsType(t, &FileLoader{}, fileLoader)
	assert.Equal(t, "/path/to/posts", fileLoader.(*FileLoader).path)

	_, err = NewLoader("unknown://")
	assert.Error(t, err)
	assert.Equal(t, ErrorUnknownLoaderType, err)

	_, err = NewLoader("~")
	assert.Error(t, err)
}
