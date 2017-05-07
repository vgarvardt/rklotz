package model

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPostFromFile(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.Contains(t, wd, "github.com/vgarvardt/rklotz")

	// .../github.com/hellofresh/auth-service/pkg/model/../../assets/posts
	postsBasePath := filepath.Join(wd, "..", "..", "assets", "posts")
	postPath := filepath.Join(postsBasePath, "hello-world.md")

	publishedAt, _ := time.Parse(time.RFC3339, "2017-05-06T16:34:00+02:00")

	post, err := NewPostFromFile(postsBasePath, postPath)
	assert.NoError(t, err)

	assert.Equal(t, "/hello-world", post.Path)
	assert.Equal(t, "Hello World Post Title", post.Title)
	assert.Equal(t, []string{"hello world", "test post", "foobar"}, post.Tags)
	assert.Equal(t, "md", post.Format)
	assert.Equal(t, publishedAt, post.PublishedAt)

	body := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
	assert.Equal(t, body, post.Body[:len(body)])

	html := "<p>" + body
	assert.Equal(t, html, post.HTML[:len(html)])
}
