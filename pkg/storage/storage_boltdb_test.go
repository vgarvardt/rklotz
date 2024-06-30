package storage

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoltDBStorage(t *testing.T) {
	dbFilePath := getFilePath(t)
	storage, err := NewBoltDBStorage(dbFilePath, 10)
	require.NoError(t, err)
	assert.Equal(t, dbFilePath, storage.path)
	assert.Equal(t, 10, storage.postsPerPage)

	assert.True(t, fileExists(dbFilePath))

	err = storage.Close()
	assert.NoError(t, err)
	assert.False(t, fileExists(dbFilePath))
}

func TestBoltDBStorage_Finalize(t *testing.T) {
	dbFilePath := getFilePath(t)
	storage, err := NewBoltDBStorage(dbFilePath, 10)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := storage.Close()
		assert.NoError(t, err)
	})

	err = storage.Finalize()
	require.NoError(t, err)
}

func TestBoltDBStorage_FindByPath(t *testing.T) {
	dbFilePath := getFilePath(t)
	storage, err := NewBoltDBStorage(dbFilePath, 10)
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

func TestBoltDBStorage_ListAll_10(t *testing.T) {
	dbFilePath := getFilePath(t)
	storage, err := NewBoltDBStorage(dbFilePath, 10)
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

	assert.Equal(t, filepath.FromSlash("/nested/nested-path"), posts[0].Path)
	assert.Equal(t, filepath.FromSlash("/hello-world"), posts[1].Path)

	posts, err = storage.ListAll(1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestBoltDBStorage_ListAll_1(t *testing.T) {
	dbFilePath := getFilePath(t)
	storage, err := NewBoltDBStorage(dbFilePath, 1)
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
	assert.Equal(t, filepath.FromSlash("/nested/nested-path"), posts[0].Path)

	posts, err = storage.ListAll(1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, filepath.FromSlash("/hello-world"), posts[0].Path)

	posts, err = storage.ListAll(2)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestBoltDBStorage_ListTag_10(t *testing.T) {
	dbFilePath := getFilePath(t)
	storage, err := NewBoltDBStorage(dbFilePath, 10)
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

	assert.Equal(t, filepath.FromSlash("/nested/nested-path"), posts[0].Path)
	assert.Equal(t, filepath.FromSlash("/hello-world"), posts[1].Path)

	posts, err = storage.ListTag(tag, 1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestBoltDBStorage_ListTag_1(t *testing.T) {
	dbFilePath := getFilePath(t)
	storage, err := NewBoltDBStorage(dbFilePath, 1)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := storage.Close()
		assert.NoError(t, err)
	})

	loadTestPosts(t, storage)

	tag := "test post"
	assert.Equal(t, 2, storage.TagMeta(tag).Pages)

	posts, err := storage.ListTag(tag, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, filepath.FromSlash("/nested/nested-path"), posts[0].Path)

	posts, err = storage.ListTag(tag, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, filepath.FromSlash("/hello-world"), posts[0].Path)

	posts, err = storage.ListTag(tag, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestBoltDBStorage_ListTag_ErrorNotFound(t *testing.T) {
	dbFilePath := getFilePath(t)
	storage, err := NewBoltDBStorage(dbFilePath, 1)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := storage.Close()
		assert.NoError(t, err)
	})

	loadTestPosts(t, storage)

	tag := getRandomHash(t, 10)
	_, err = storage.ListTag(tag, 0)
	require.Error(t, err)
	assert.ErrorIs(t, ErrorNotFound, err)
}
