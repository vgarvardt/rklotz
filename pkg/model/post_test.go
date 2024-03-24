package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vgarvardt/rklotz/pkg/formatter"
)

func TestNewPostFromFile(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	// .../pkg/model/../../assets/posts
	postsBasePath := filepath.Join(wd, "..", "..", "assets", "posts")
	testCases := []struct {
		path           string
		publishedAt    string
		postPath       string
		postTitle      string
		postTags       []string
		posyFormat     string
		postBody       string
		postBodyHTML   string
		postTeaser     string
		postTeaserHTML string
	}{{
		path:           "hello-world.md",
		publishedAt:    "2017-05-06T16:34:00+02:00",
		postPath:       "/hello-world",
		postTitle:      "Hello World Post Title",
		postTags:       []string{"hello world", "test post", "foobar"},
		posyFormat:     "md",
		postBody:       strings.TrimSpace(helloWorldBody),
		postBodyHTML:   helloWorldBodyHTML,
		postTeaser:     "",
		postTeaserHTML: "",
	}, {
		path:           "with-teaser.md",
		publishedAt:    "2020-02-23T15:00:00+02:00",
		postPath:       "/with-teaser",
		postTitle:      "Post With Teaser",
		postTags:       []string{"test post", "foobar", "teaser"},
		posyFormat:     "md",
		postBody:       strings.TrimSpace(withTeaserBody),
		postBodyHTML:   withTeaserBodyHTML,
		postTeaser:     strings.TrimSpace(withTeaserTeaser),
		postTeaserHTML: withTeaserTeaserHTML,
	}}

	for _, tt := range testCases {
		t.Run(tt.path, func(t *testing.T) {
			postPath := filepath.Join(postsBasePath, tt.path)

			publishedAt, err := time.Parse(time.RFC3339, tt.publishedAt)
			require.NoError(t, err)

			post, err := NewPostFromFile(postsBasePath, postPath, formatter.New())
			require.NoError(t, err)

			assert.Equal(t, tt.postPath, post.Path)
			assert.Equal(t, tt.postTitle, post.Title)
			assert.Equal(t, tt.postTags, post.Tags)
			assert.Equal(t, tt.posyFormat, post.Format)
			assert.Equal(t, publishedAt, post.PublishedAt)

			assert.Equal(t, tt.postBody, post.Body)
			assert.Equal(t, tt.postBodyHTML, post.BodyHTML)
			assert.Equal(t, tt.postTeaser, post.Teaser)
			assert.Equal(t, tt.postTeaserHTML, post.TeaserHTML)
		})
	}
}

const (
	helloWorldBody = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eget risus in lorem convallis semper.
Sed posuere vehicula feugiat. Maecenas facilisis nunc nisl, sit amet ornare quam scelerisque vel.
Vestibulum non nunc justo. Donec vitae justo ipsum. Cras tempor nec tortor vitae suscipit.
In vulputate lorem id quam tincidunt, non pulvinar dui various. Sed a imperdiet orci.
Aliquam et sem in tellus dapibus lobortis. Quisque auctor laoreet massa, in tincidunt lectus rutrum vitae.

Vestibulum hendrerit massa libero, et sagittis felis luctus ut. Nunc condimentum aliquet lectus,
id posuere risus rhoncus et. Vivamus sed diam aliquam, gravida neque ut, luctus purus.
Mauris fringilla sagittis pretium. In egestas urna lectus, semper vehicula libero eleifend vitae.
Duis vitae dolor sit amet purus eleifend venenatis in vitae ligula. In quis est libero.

Pellentesque ultrices massa blandit, pellentesque tortor eu, sagittis orci. Aliquam erat volutpat.
Duis pharetra malesuada nisi, eu semper est luctus vel. Quisque ac nisl sapien. Etiam eros lorem,
auctor ac placerat sit amet, egestas sed lectus. Curabitur dolor odio, bibendum vitae ex sed, viverra commodo ligula.
Donec eget sem ex. Sed eleifend hendrerit purus id euismod. Vestibulum mauris elit, egestas non risus sed,
volutpat iaculis lorem. Fusce rutrum quam et lacus iaculis iaculis. Etiam dictum neque a justo finibus,
sit amet hendrerit sem placerat. Aliquam bibendum ex sit amet nisi pharetra condimentum.
`
	helloWorldBodyHTML = `<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eget risus in lorem convallis semper.
Sed posuere vehicula feugiat. Maecenas facilisis nunc nisl, sit amet ornare quam scelerisque vel.
Vestibulum non nunc justo. Donec vitae justo ipsum. Cras tempor nec tortor vitae suscipit.
In vulputate lorem id quam tincidunt, non pulvinar dui various. Sed a imperdiet orci.
Aliquam et sem in tellus dapibus lobortis. Quisque auctor laoreet massa, in tincidunt lectus rutrum vitae.</p>
<p>Vestibulum hendrerit massa libero, et sagittis felis luctus ut. Nunc condimentum aliquet lectus,
id posuere risus rhoncus et. Vivamus sed diam aliquam, gravida neque ut, luctus purus.
Mauris fringilla sagittis pretium. In egestas urna lectus, semper vehicula libero eleifend vitae.
Duis vitae dolor sit amet purus eleifend venenatis in vitae ligula. In quis est libero.</p>
<p>Pellentesque ultrices massa blandit, pellentesque tortor eu, sagittis orci. Aliquam erat volutpat.
Duis pharetra malesuada nisi, eu semper est luctus vel. Quisque ac nisl sapien. Etiam eros lorem,
auctor ac placerat sit amet, egestas sed lectus. Curabitur dolor odio, bibendum vitae ex sed, viverra commodo ligula.
Donec eget sem ex. Sed eleifend hendrerit purus id euismod. Vestibulum mauris elit, egestas non risus sed,
volutpat iaculis lorem. Fusce rutrum quam et lacus iaculis iaculis. Etiam dictum neque a justo finibus,
sit amet hendrerit sem placerat. Aliquam bibendum ex sit amet nisi pharetra condimentum.</p>
`
	withTeaserBody = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eget risus in lorem convallis semper.
Sed posuere vehicula feugiat. Maecenas facilisis nunc nisl, sit amet ornare quam scelerisque vel.
Vestibulum non nunc justo. Donec vitae justo ipsum. Cras tempor nec tortor vitae suscipit.
In vulputate lorem id quam tincidunt, non pulvinar dui various. Sed a imperdiet orci.
Aliquam et sem in tellus dapibus lobortis. Quisque auctor laoreet massa, in tincidunt lectus rutrum vitae.

Vestibulum hendrerit massa libero, et sagittis felis luctus ut. Nunc condimentum aliquet lectus,
id posuere risus rhoncus et. Vivamus sed diam aliquam, gravida neque ut, luctus purus.
Mauris fringilla sagittis pretium. In egestas urna lectus, semper vehicula libero eleifend vitae.
Duis vitae dolor sit amet purus eleifend venenatis in vitae ligula. In quis est libero.

Pellentesque ultrices massa blandit, pellentesque tortor eu, sagittis orci. Aliquam erat volutpat.
Duis pharetra malesuada nisi, eu semper est luctus vel. Quisque ac nisl sapien. Etiam eros lorem,
auctor ac placerat sit amet, egestas sed lectus. Curabitur dolor odio, bibendum vitae ex sed, viverra commodo ligula.
Donec eget sem ex. Sed eleifend hendrerit purus id euismod. Vestibulum mauris elit, egestas non risus sed,
volutpat iaculis lorem. Fusce rutrum quam et lacus iaculis iaculis. Etiam dictum neque a justo finibus,
sit amet hendrerit sem placerat. Aliquam bibendum ex sit amet nisi pharetra condimentum.
`
	withTeaserBodyHTML = `<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eget risus in lorem convallis semper.
Sed posuere vehicula feugiat. Maecenas facilisis nunc nisl, sit amet ornare quam scelerisque vel.
Vestibulum non nunc justo. Donec vitae justo ipsum. Cras tempor nec tortor vitae suscipit.
In vulputate lorem id quam tincidunt, non pulvinar dui various. Sed a imperdiet orci.
Aliquam et sem in tellus dapibus lobortis. Quisque auctor laoreet massa, in tincidunt lectus rutrum vitae.</p>
<p>Vestibulum hendrerit massa libero, et sagittis felis luctus ut. Nunc condimentum aliquet lectus,
id posuere risus rhoncus et. Vivamus sed diam aliquam, gravida neque ut, luctus purus.
Mauris fringilla sagittis pretium. In egestas urna lectus, semper vehicula libero eleifend vitae.
Duis vitae dolor sit amet purus eleifend venenatis in vitae ligula. In quis est libero.</p>
<p>Pellentesque ultrices massa blandit, pellentesque tortor eu, sagittis orci. Aliquam erat volutpat.
Duis pharetra malesuada nisi, eu semper est luctus vel. Quisque ac nisl sapien. Etiam eros lorem,
auctor ac placerat sit amet, egestas sed lectus. Curabitur dolor odio, bibendum vitae ex sed, viverra commodo ligula.
Donec eget sem ex. Sed eleifend hendrerit purus id euismod. Vestibulum mauris elit, egestas non risus sed,
volutpat iaculis lorem. Fusce rutrum quam et lacus iaculis iaculis. Etiam dictum neque a justo finibus,
sit amet hendrerit sem placerat. Aliquam bibendum ex sit amet nisi pharetra condimentum.</p>
`
	withTeaserTeaser = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eget risus in lorem convallis semper.
Sed posuere vehicula feugiat. Maecenas facilisis nunc nisl, sit amet ornare quam scelerisque vel.
Vestibulum non nunc justo. Donec vitae justo ipsum. Cras tempor nec tortor vitae suscipit.
In vulputate lorem id quam tincidunt, non pulvinar dui various. Sed a imperdiet orci.
Aliquam et sem in tellus dapibus lobortis. Quisque auctor laoreet massa, in tincidunt lectus rutrum vitae.
`
	withTeaserTeaserHTML = `<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eget risus in lorem convallis semper.
Sed posuere vehicula feugiat. Maecenas facilisis nunc nisl, sit amet ornare quam scelerisque vel.
Vestibulum non nunc justo. Donec vitae justo ipsum. Cras tempor nec tortor vitae suscipit.
In vulputate lorem id quam tincidunt, non pulvinar dui various. Sed a imperdiet orci.
Aliquam et sem in tellus dapibus lobortis. Quisque auctor laoreet massa, in tincidunt lectus rutrum vitae.</p>
`
)
