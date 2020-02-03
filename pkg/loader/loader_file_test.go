package loader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/vgarvardt/rklotz/pkg/formatter"
	"github.com/vgarvardt/rklotz/pkg/model"
)

type mockStorage struct {
	saveCallCount  int
	saveCallParams []*model.Post
	saveCallResult []error
}

func (s *mockStorage) Save(post *model.Post) error {
	currentCall := s.saveCallCount
	s.saveCallCount++
	s.saveCallParams = append(s.saveCallParams, post)

	return s.saveCallResult[currentCall]
}

func (s *mockStorage) Finalize() error {
	return nil
}

func (s *mockStorage) FindByPath(path string) (*model.Post, error) {
	return nil, nil
}

func (s *mockStorage) ListAll(page int) ([]*model.Post, error) {
	return nil, nil
}

func (s *mockStorage) ListTag(tag string, page int) ([]*model.Post, error) {
	return nil, nil
}

func (s *mockStorage) Close() error {
	return nil
}

func (s *mockStorage) Meta() *model.Meta {
	return &model.Meta{
		Posts:   len(s.saveCallParams),
		PerPage: 0,
		Pages:   0,
	}
}

func (s *mockStorage) TagMeta(tag string) *model.Meta {
	return &model.Meta{
		Posts:   len(s.saveCallParams),
		PerPage: 0,
		Pages:   0,
	}
}

func TestFileLoader_Load(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	// .../github.com/vgarvardt/rklotz/pkg/model/../../assets/posts
	postsBasePath := filepath.Join(wd, "..", "..", "assets", "posts")

	storage := &mockStorage{saveCallResult: []error{nil, nil}}

	f := formatter.New()

	fileLoader, err := NewFileLoader(postsBasePath, f, zap.NewNop())
	require.NoError(t, err)

	err = fileLoader.Load(storage)
	require.NoError(t, err)

	assert.Equal(t, 2, storage.saveCallCount)
	assert.Equal(t, 2, len(storage.saveCallParams))

	found := false
	for _, post := range storage.saveCallParams {
		if post.Path == "/hello-world" {
			found = true
			assert.Equal(t, "Hello World Post Title", post.Title)
		}
	}
	assert.True(t, found)

	found = false
	for _, post := range storage.saveCallParams {
		if post.Path == "/nested/nested-path" {
			found = true
			assert.Equal(t, "Nested Path Post Title", post.Title)
		}
	}
	assert.True(t, found)
}
