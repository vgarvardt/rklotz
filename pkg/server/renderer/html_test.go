package renderer

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewHTML(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	// .../rklotz/pkg/server/renderer/../../../templates
	templatesPath := filepath.Join(wd, "..", "..", "..", "templates")
	theme := "foundation6"
	expected := []string{
		path.Join(templatesPath, "plugins", "disqus.html"),
		path.Join(templatesPath, "plugins", "ga.html"),
		path.Join(templatesPath, "plugins", "gtm-body.html"),
		path.Join(templatesPath, "plugins", "gtm-head.html"),
		path.Join(templatesPath, "plugins", "highlightjs-css.html"),
		path.Join(templatesPath, "plugins", "highlightjs-js.html"),
		path.Join(templatesPath, "plugins", "yamka.html"),
		path.Join(templatesPath, "plugins", "yasha.html"),

		path.Join(templatesPath, theme, "partial", "heading.html"),
		path.Join(templatesPath, theme, "partial", "info.html"),
		path.Join(templatesPath, theme, "partial", "pagination.html"),
		path.Join(templatesPath, theme, "partial", "posts.html"),
	}

	instance := &HTML{logger: zap.NewNop()}

	// default about panel
	partials, err := instance.getPartials(templatesPath, theme, "")
	require.NoError(t, err)
	assert.Equal(t, append(expected, path.Join(templatesPath, theme, "partial", "about.html")), partials)

	// custom about panel
	partials, err = instance.getPartials(templatesPath, theme, path.Join(templatesPath, theme, "partial", "pagination.html"))
	require.NoError(t, err)
	assert.Equal(t, append(expected, path.Join(templatesPath, theme, "partial", "pagination.html")), partials)
}
