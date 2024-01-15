package loader

import (
	"log/slog"
	"testing"

	"github.com/cappuccinotm/slogx/slogt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoader(t *testing.T) {
	logger := slog.New(slogt.Handler(t, slogt.SplitMultiline))

	fileLoader, err := New("file:///path/to/posts", logger)
	require.NoError(t, err)
	assert.IsType(t, &FileLoader{}, fileLoader)
	assert.Equal(t, "/path/to/posts", fileLoader.(*FileLoader).path)

	_, err = New("unknown://", logger)
	require.Error(t, err)
	assert.Equal(t, ErrorUnknownLoaderType, err)

	_, err = New("~", logger)
	require.Error(t, err)
}
