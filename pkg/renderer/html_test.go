package renderer

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLRendererData(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.Contains(t, wd, "github.com/vgarvardt/rklotz")

	// .../github.com/vgarvardt/rklotz/pkg/renderer/../../templates
	templatesPath := filepath.Join(wd, "..", "..", "templates")
	theme := "foundation"
	expected := []string{
		path.Join(templatesPath, "plugins", "disqus.html"),
		path.Join(templatesPath, "plugins", "ga.html"),
		path.Join(templatesPath, "plugins", "highlightjs-css.html"),
		path.Join(templatesPath, "plugins", "highlightjs-js.html"),
		path.Join(templatesPath, "plugins", "yamka.html"),
		path.Join(templatesPath, "plugins", "yasha.html"),

		path.Join(templatesPath, theme, "partial", "alert.html"),
		path.Join(templatesPath, theme, "partial", "heading.html"),
		path.Join(templatesPath, theme, "partial", "info.html"),
		path.Join(templatesPath, theme, "partial", "pagination.html"),
		path.Join(templatesPath, theme, "partial", "posts.html"),
	}

	instance := &htmlRenderer{}

	// default about panel
	partials, err := instance.getPartials(templatesPath, theme, "")
	assert.NoError(t, err)
	assert.Equal(t, append(expected, path.Join(templatesPath, theme, "partial", "about.html")), partials)

	// custom about panel
	partials, err = instance.getPartials(templatesPath, theme, path.Join(templatesPath, theme, "partial", "alert.html"))
	assert.NoError(t, err)
	assert.Equal(t, append(expected, path.Join(templatesPath, theme, "partial", "alert.html")), partials)
}
