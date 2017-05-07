package repository

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getFilePath() string {
	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	return fmt.Sprintf("/tmp/rklotz-test.%s.db", hex.EncodeToString(hasher.Sum(nil))[:5])
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func TestNewBoltDBStorage(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 10)
	assert.NoError(t, err)
	assert.Equal(t, dbFilePath, storage.path)

	assert.True(t, fileExists(dbFilePath))

	err = storage.Close()
	assert.NoError(t, err)
	assert.False(t, fileExists(dbFilePath))
}

func TestBoltDBStorage_Reindex(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 10)
	assert.NoError(t, err)
	defer storage.Close()

	err = storage.Reindex(10)
	assert.NoError(t, err)
}

func loadTestPosts(t *testing.T, storage Storage) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.Contains(t, wd, "github.com/vgarvardt/rklotz")

	// .../github.com/hellofresh/auth-service/pkg/model/../../assets/posts
	postsBasePath := filepath.Join(wd, "..", "..", "assets", "posts")

	fileLoader, err := NewFileLoader(postsBasePath)
	assert.NoError(t, err)

	err = fileLoader.Load(storage)
	assert.NoError(t, err)
}

func TestBoltDBStorage_FindByPath(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 10)
	assert.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)

	_, err = storage.FindByPath("does-not-exist")
	assert.Equal(t, err, ErrorNotFound)

	post, err := storage.FindByPath("/hello-world")
	assert.NoError(t, err)
	assert.Equal(t, "/hello-world", post.Path)
	assert.Equal(t, "Hello World Post Title", post.Title)

	post, err = storage.FindByPath("/nested/nested-path")
	assert.NoError(t, err)
	assert.Equal(t, "/nested/nested-path", post.Path)
	assert.Equal(t, "Nested Path Post Title", post.Title)
}

func TestBoltDBStorage_ListAll_10(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 10)
	assert.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)

	posts, err := storage.ListAll(0)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(posts))

	assert.Equal(t, "/nested/nested-path", posts[0].Path)
	assert.Equal(t, "/hello-world", posts[1].Path)

	posts, err = storage.ListAll(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestBoltDBStorage_ListAll_1(t *testing.T) {
	dbFilePath := getFilePath()
	storage, err := NewBoltDBStorage(dbFilePath, 1)
	assert.NoError(t, err)
	defer storage.Close()

	loadTestPosts(t, storage)

	posts, err := storage.ListAll(0)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, "/nested/nested-path", posts[0].Path)

	posts, err = storage.ListAll(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(posts))
	assert.Equal(t, "/hello-world", posts[0].Path)

	posts, err = storage.ListAll(2)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(posts))
}

func TestBoltDBStorage_ListTag_10(t *testing.T) {
	//dbFilePath := getFilePath()
	//storage, err := NewBoltDBStorage(dbFilePath, 10)
	//assert.NoError(t, err)
	//defer storage.Close()
	//
	//loadTestPosts(t, storage)
	//
	//posts, err := storage.ListTag("test post", 0)
	//assert.NoError(t, err)
	//assert.Equal(t, 2, len(posts))
	//
	//assert.Equal(t, "/nested/nested-path", posts[0].Path)
	//assert.Equal(t, "/hello-world", posts[1].Path)
	//
	//posts, err = storage.ListAll(1)
	//assert.NoError(t, err)
	//assert.Equal(t, 0, len(posts))
}
