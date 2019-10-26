package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vgarvardt/rklotz/pkg/model"
)

func getRandomHash(length int) string {
	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	return hex.EncodeToString(hasher.Sum(nil))[:length]
}

func getFilePath() string {
	return fmt.Sprintf("/tmp/rklotz-test.%s.db", getRandomHash(5))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func TestNewBoltDBStorage(t *testing.T) {
	dbFilePath := getFilePath()
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
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 10)
	require.NoError(t, err)
	defer storage.Close()

	err = storage.Finalize()
	require.NoError(t, err)
}

func loadTestPosts(t *testing.T, storage Storage) {
	t.Helper()

	wd, err := os.Getwd()
	require.NoError(t, err)

	// .../github.com/vgarvardt/rklotz/pkg/repository/../../assets/posts
	postsBasePath := filepath.Join(wd, "..", "..", "assets", "posts")

	post1, err := model.NewPostFromFile(
		postsBasePath,
		filepath.Join(postsBasePath, "hello-world.md"),
	)
	require.NoError(t, err)
	err = storage.Save(post1)
	require.NoError(t, err)

	post2, err := model.NewPostFromFile(
		postsBasePath,
		filepath.Join(postsBasePath, "nested/nested-path.md"),
	)
	require.NoError(t, err)
	err = storage.Save(post2)
	require.NoError(t, err)

	err = storage.Finalize()
	require.NoError(t, err)

	require.Equal(t, 2, storage.Meta().Posts)
}

func TestBoltDBStorage_FindByPath(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 10)
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

func TestBoltDBStorage_ListAll_10(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 10)
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

func TestBoltDBStorage_ListAll_1(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 1)
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

func TestBoltDBStorage_ListTag_10(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 10)
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

func TestBoltDBStorage_ListTag_1(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 1)
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

func TestBoltDBStorage_ListTag_ErrorNotFound(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 1)
	require.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)

	tag := getRandomHash(10)
	_, err = storage.ListTag(tag, 0)
	require.Error(t, err)
	assert.Equal(t, ErrorNotFound, err)
}
