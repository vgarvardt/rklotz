package storage

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getBadgerPath(t *testing.T) string {
	t.Helper()

	return path.Join(os.TempDir(), fmt.Sprintf("rklotz-test.%s", getRandomHash(t, 5)))
}

func TestNewBadgerStorage(t *testing.T) {
	dbFilePath := getBadgerPath(t)
	storage, err := NewBadgerStorage(dbFilePath, 10)
	require.NoError(t, err)
	assert.Equal(t, dbFilePath, storage.path)
	assert.Equal(t, 10, storage.postsPerPage)

	assert.True(t, fileExists(dbFilePath))

	err = storage.Close()
	assert.NoError(t, err)
	assert.False(t, fileExists(dbFilePath))
}

func TestBadgerStorage_FindByPath(t *testing.T) {
	dbFilePath := getBadgerPath(t)
	storage, err := NewBadgerStorage(dbFilePath, 10)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := storage.Close()
		assert.NoError(t, err)
	})

	loadTestPosts(t, storage)

	_, err = storage.FindByPath("does-not-exist")
	assert.ErrorIs(t, err, ErrorNotFound)

	post, err := storage.FindByPath(filepath.FromSlash("/hello-world"))
	require.NoError(t, err)
	assert.Equal(t, filepath.FromSlash("/hello-world"), post.Path)
	assert.Equal(t, "Hello World Post Title\r", post.Title)

	post, err = storage.FindByPath(filepath.FromSlash("/nested/nested-path"))
	require.NoError(t, err)
	assert.Equal(t, filepath.FromSlash("/nested/nested-path"), post.Path)
	assert.Equal(t, "Nested Path Post Title\r", post.Title)
}

func TestBadgerStorage_ListAll_10(t *testing.T) {
	dbFilePath := getBadgerPath(t)
	storage, err := NewBadgerStorage(dbFilePath, 10)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := storage.Close()
		assert.NoError(t, err)
	})

	loadTestPosts(t, storage)
	assert.Equal(t, 1, storage.Meta().Pages)

	posts, err := storage.ListAll(0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(posts))

	posts, err = storage.ListAll(1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestBadgerStorage_ListAll_1(t *testing.T) {
	dbFilePath := getBadgerPath(t)
	storage, err := NewBadgerStorage(dbFilePath, 1)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := storage.Close()
		assert.NoError(t, err)
	})

	loadTestPosts(t, storage)
	assert.Equal(t, 2, storage.Meta().Pages)

	posts, err := storage.ListAll(0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(posts))

	posts, err = storage.ListAll(1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(posts))

	posts, err = storage.ListAll(2)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestBadgerStorage_ListTag_10(t *testing.T) {
	dbFilePath := getBadgerPath(t)
	storage, err := NewBadgerStorage(dbFilePath, 10)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := storage.Close()
		assert.NoError(t, err)
	})

	loadTestPosts(t, storage)

	tag := "test post"
	assert.Equal(t, 1, storage.TagMeta(tag).Pages)

	posts, err := storage.ListTag(tag, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(posts))

	// FIXME
	//posts, err = storage.ListTag(tag, 1)
	//require.NoError(t, err)
	//assert.Equal(t, 0, len(posts))
}
