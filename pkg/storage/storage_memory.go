package storage

import (
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/vgarvardt/rklotz/pkg/model"
)

type postSlice []*model.Post

// Len is the number of elements in the collection.
func (p postSlice) Len() int { return len(p) }

// Less reports whether the element with
// index i should sort before the element with index j.
func (p postSlice) Less(i, j int) bool { return p[i].PublishedAt.Before(p[j].PublishedAt) }

// Swap swaps the elements with indexes i and j.
func (p postSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// MemoryStorage is the in memory Storage implementation
type MemoryStorage struct {
	posts *sync.Map
	tags  *sync.Map

	postsList postSlice
	tagsList  *sync.Map

	postsCount   int
	postsPerPage int
}

// NewMemoryStorage creates new MemoryStorage instance
func NewMemoryStorage(postsPerPage int) (*MemoryStorage, error) {
	instance := &MemoryStorage{
		posts:        new(sync.Map),
		tags:         new(sync.Map),
		tagsList:     new(sync.Map),
		postsPerPage: postsPerPage,
	}

	return instance, nil
}

// Save persists new post in the storage
func (s *MemoryStorage) Save(post *model.Post) error {
	s.posts.Store(post.Path, post)
	s.postsCount++

	for i := range post.Tags {
		tag := &model.Tag{Tag: post.Tags[i], Paths: []string{post.Path}}
		tagSlice := postSlice{post}

		loadedTag, okTag := s.tags.LoadOrStore(strings.ToLower(post.Tags[i]), tag)
		loadedTagSlice, okTagSlice := s.tagsList.LoadOrStore(strings.ToLower(post.Tags[i]), tagSlice)

		if okTag && okTagSlice {
			loadedTag.(*model.Tag).Paths = append(loadedTag.(*model.Tag).Paths, post.Path)
			loadedTagSlice = append(loadedTagSlice.(postSlice), post)

			s.tags.Store(strings.ToLower(post.Tags[i]), loadedTag)
			s.tagsList.Store(strings.ToLower(post.Tags[i]), loadedTagSlice)
		}
	}

	return nil
}

// Finalize is called after all posts are persisted in the storage
func (s *MemoryStorage) Finalize() error {
	s.postsList = make([]*model.Post, s.postsCount)
	i := 0
	s.posts.Range(func(path, post interface{}) bool {
		s.postsList[i] = post.(*model.Post)
		i++
		return true
	})
	sort.Sort(sort.Reverse(s.postsList))

	s.tagsList.Range(func(tag, tagSlice interface{}) bool {
		sort.Sort(sort.Reverse(tagSlice.(postSlice)))
		return true
	})

	return nil
}

// FindByPath searches for a post by path
func (s *MemoryStorage) FindByPath(path string) (*model.Post, error) {
	post, ok := s.posts.Load(path)
	if !ok {
		return nil, ErrorNotFound
	}
	return post.(*model.Post), nil
}

// ListAll returns ordered by date posts page
func (s *MemoryStorage) ListAll(page int) ([]*model.Post, error) {
	return s.slicePage(s.postsList, page)
}

// ListTag returns ordered by date posts page for a tag
func (s *MemoryStorage) ListTag(tag string, page int) ([]*model.Post, error) {
	tagSlice, ok := s.tagsList.Load(strings.ToLower(tag))
	if !ok {
		return nil, ErrorNotFound
	}
	return s.slicePage(tagSlice.(postSlice), page)
}

func (s *MemoryStorage) slicePage(slice []*model.Post, page int) ([]*model.Post, error) {
	offset := page * s.postsPerPage
	if offset > len(slice) {
		return []*model.Post{}, nil
	}

	offsetBound := int(math.Min(float64(len(slice)), float64(offset+s.postsPerPage)))

	return slice[offset:offsetBound], nil
}

// Close closes the storage and frees all resources
func (s *MemoryStorage) Close() error {
	return nil
}

// Meta returns metadata for all persisted posts
func (s *MemoryStorage) Meta() *model.Meta {
	return model.NewMeta(s.postsCount, s.postsPerPage)
}

// TagMeta returns metadata for all persisted posts for a tag
func (s *MemoryStorage) TagMeta(tag string) *model.Meta {
	tagModel, ok := s.tags.Load(strings.ToLower(tag))
	if !ok {
		return &model.Meta{}
	}
	return model.NewMeta(len(tagModel.(*model.Tag).Paths), s.postsPerPage)
}
