package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemoryStorage(t *testing.T) {
	storage, err := NewMemoryStorage(10)
	require.NoError(t, err)
	assert.Equal(t, 10, storage.postsPerPage)
}

func TestMemoryStorage_Finalize(t *testing.T) {
	storage, err := NewMemoryStorage(10)
	require.NoError(t, err)
	defer storage.Close()

	err = storage.Finalize()
	require.NoError(t, err)
}

func TestMemoryStorage_FindByPath(t *testing.T) {
	storage, err := NewMemoryStorage(10)
	require.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)

	_, err = storage.FindByPath("does-not-exist")
	assert.Equal(t, err, ErrorNotFound)

	post, err := storage.FindByPath("/hello-world")
	require.NoError(t, err)
	assert.Equal(t, "/hello-world", post.Path)
	assert.Equal(t, "Hello World Post Title", post.Title)

	post, err = storage.FindByPath("/nested/nested-path")
	require.NoError(t, err)
	assert.Equal(t, "/nested/nested-path", post.Path)
	assert.Equal(t, "Nested Path Post Title", post.Title)
}

func TestMemoryStorage_ListAll_10(t *testing.T) {
	storage, err := NewMemoryStorage(10)
	require.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)
	assert.Equal(t, 1, storage.Meta().Pages)

	posts, err := storage.ListAll(0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(posts))

	assert.Equal(t, "/nested/nested-path", posts[0].Path)
	assert.Equal(t, "/hello-world", posts[1].Path)

	posts, err = storage.ListAll(1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestMemoryStorage_ListAll_1(t *testing.T) {
	storage, err := NewMemoryStorage(1)
	require.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)
	assert.Equal(t, 2, storage.Meta().Pages)

	posts, err := storage.ListAll(0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, "/nested/nested-path", posts[0].Path)

	posts, err = storage.ListAll(1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, "/hello-world", posts[0].Path)

	posts, err = storage.ListAll(2)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestMemoryStorage_ListTag_10(t *testing.T) {
	storage, err := NewMemoryStorage(10)
	require.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)

	tag := "test post"
	assert.Equal(t, 1, storage.TagMeta(tag).Pages)

	posts, err := storage.ListTag(tag, 0)
	require.NoError(t, err)
	assert.Equal(t, 2, len(posts))

	assert.Equal(t, "/nested/nested-path", posts[0].Path)
	assert.Equal(t, "/hello-world", posts[1].Path)

	posts, err = storage.ListTag(tag, 1)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestMemoryStorage_ListTag_1(t *testing.T) {
	storage, err := NewMemoryStorage(1)
	require.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)

	tag := "test post"
	assert.Equal(t, 2, storage.TagMeta(tag).Pages)

	posts, err := storage.ListTag(tag, 0)
	require.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, "/nested/nested-path", posts[0].Path)

	posts, err = storage.ListTag(tag, 1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, "/hello-world", posts[0].Path)

	posts, err = storage.ListTag(tag, 2)
	require.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestMemoryStorage_ListTag_ErrorNotFound(t *testing.T) {
	storage, err := NewMemoryStorage(1)
	require.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)

	tag := getRandomHash(10)
	_, err = storage.ListTag(tag, 0)
	require.Error(t, err)
	assert.Equal(t, ErrorNotFound, err)
}
