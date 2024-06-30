package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/vgarvardt/rklotz/pkg/formatter"
	"github.com/vgarvardt/rklotz/pkg/model"
)

func getRandomHash(t *testing.T, length int) string {
	t.Helper()

	hasher := md5.New()
	_, err := hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	require.NoError(t, err)

	return hex.EncodeToString(hasher.Sum(nil))[:length]
}

func getFilePath(t *testing.T) string {
	t.Helper()

	return fmt.Sprintf("/tmp/rklotz-test.%s.db", getRandomHash(t, 5))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func loadTestPosts(t *testing.T, storage Storage) {
	t.Helper()

	wd, err := os.Getwd()
	require.NoError(t, err)

	f := formatter.New()

	// ../../assets/posts
	postsBasePath := filepath.Join(wd, "..", "..", "assets", "posts")

	post1, err := model.NewPostFromFile(
		postsBasePath,
		filepath.Join(postsBasePath, "hello-world.md"),
		f,
	)
	require.NoError(t, err)
	err = storage.Save(post1)
	require.NoError(t, err)

	post2, err := model.NewPostFromFile(
		postsBasePath,
		filepath.Join(postsBasePath, "nested/nested-path.md"),
		f,
	)
	require.NoError(t, err)
	err = storage.Save(post2)
	require.NoError(t, err)

	err = storage.Finalize()
	require.NoError(t, err)

	require.Equal(t, 2, storage.Meta().Posts)
}
