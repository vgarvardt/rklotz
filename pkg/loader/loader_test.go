package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewLoader(t *testing.T) {
	fileLoader, err := New("file:///path/to/posts", zap.NewNop())
	assert.NoError(t, err)
	assert.IsType(t, &FileLoader{}, fileLoader)
	assert.Equal(t, "/path/to/posts", fileLoader.(*FileLoader).path)

	_, err = New("unknown://", zap.NewNop())
	assert.Error(t, err)
	assert.Equal(t, ErrorUnknownLoaderType, err)

	_, err = New("~", zap.NewNop())
	assert.Error(t, err)
}
