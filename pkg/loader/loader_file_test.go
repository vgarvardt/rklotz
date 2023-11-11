package loader

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/cappuccinotm/slogx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func (s *mockStorage) FindByPath(string) (*model.Post, error) {
	return nil, nil
}

func (s *mockStorage) ListAll(int) ([]*model.Post, error) {
	return nil, nil
}

func (s *mockStorage) ListTag(string, int) ([]*model.Post, error) {
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

func (s *mockStorage) TagMeta(string) *model.Meta {
	return &model.Meta{
		Posts:   len(s.saveCallParams),
		PerPage: 0,
		Pages:   0,
	}
}

func TestFileLoader_Load(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	// .../pkg/loader/../../assets/posts
	postsBasePath := filepath.Join(wd, "..", "..", "assets", "posts")

	storage := &mockStorage{saveCallResult: []error{nil, nil, nil}}

	logger := slog.New(slogx.TestHandler(t))
	f := formatter.New()

	fileLoader, err := NewFileLoader(postsBasePath, f, logger)
	require.NoError(t, err)

	err = fileLoader.Load(storage)
	require.NoError(t, err)

	require.Equal(t, 3, storage.saveCallCount)
	require.Equal(t, 3, len(storage.saveCallParams))

	for _, p := range []struct {
		path  string
		title string
	}{{
		path:  "/hello-world",
		title: "Hello World Post Title",
	}, {
		path:  "/nested/nested-path",
		title: "Nested Path Post Title",
	}, {
		path:  "/with-teaser",
		title: "Post With Teaser",
	}} {
		t.Run(p.path, func(t *testing.T) {
			var found bool
			for _, post := range storage.saveCallParams {
				if post.Path == p.path {
					found = true
					assert.Equal(t, p.title, post.Title)
				}
			}
			assert.True(t, found)
		})
	}
}
