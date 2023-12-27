package loader

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/cappuccinotm/slogx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/vgarvardt/rklotz/pkg/formatter"
	"github.com/vgarvardt/rklotz/pkg/model"
)

type mockStorage struct {
	mock.Mock
}

func (s *mockStorage) Save(post *model.Post) error {
	args := s.Called(post)
	return args.Error(0)
}

func (s *mockStorage) Finalize() error {
	args := s.Called()
	return args.Error(0)
}

func (s *mockStorage) FindByPath(path string) (*model.Post, error) {
	args := s.Called(path)
	arg0 := args.Get(0)
	if arg0 == nil {
		return nil, args.Error(1)
	}

	return arg0.(*model.Post), args.Error(1)
}

func (s *mockStorage) ListAll(page int) ([]*model.Post, error) {
	args := s.Called(page)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (s *mockStorage) ListTag(tag string, page int) ([]*model.Post, error) {
	args := s.Called(tag, page)
	return args.Get(0).([]*model.Post), args.Error(1)
}

func (s *mockStorage) Close() error {
	args := s.Called()
	return args.Error(0)
}

func (s *mockStorage) Meta() *model.Meta {
	args := s.Called()
	return args.Get(0).(*model.Meta)
}

func (s *mockStorage) TagMeta(tag string) *model.Meta {
	args := s.Called(tag)
	return args.Get(0).(*model.Meta)
}

func TestFileLoader_Load(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	// .../pkg/loader/../../assets/posts
	postsBasePath := filepath.Join(wd, "..", "..", "assets", "posts")

	storage := new(mockStorage)

	var savedPosts []*model.Post
	storage.On("Save", mock.MatchedBy(func(post *model.Post) bool {
		savedPosts = append(savedPosts, post)
		return true
	})).Return(nil).Times(3)

	storage.On("Finalize").Return(nil).Once()

	logger := slog.New(slogx.TestHandler(t))
	f := formatter.New()

	fileLoader, err := NewFileLoader(postsBasePath, f, logger)
	require.NoError(t, err)

	err = fileLoader.Load(storage)
	require.NoError(t, err)

	storage.AssertExpectations(t)

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
			for _, post := range savedPosts {
				if post.Path == p.path {
					found = true
					assert.Equal(t, p.title, post.Title)
					break
				}
			}
			assert.True(t, found)
		})
	}
}
