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
		path.Join(templatesPath, "plugins", "disqus.tpl"),
		path.Join(templatesPath, "plugins", "ga.tpl"),
		path.Join(templatesPath, "plugins", "gtm-body.tpl"),
		path.Join(templatesPath, "plugins", "gtm-head.tpl"),
		path.Join(templatesPath, "plugins", "highlightjs-css.tpl"),
		path.Join(templatesPath, "plugins", "highlightjs-js.tpl"),
		path.Join(templatesPath, "plugins", "yamka.tpl"),
		path.Join(templatesPath, "plugins", "yasha.tpl"),

		path.Join(templatesPath, theme, "partial", "heading.tpl"),
		path.Join(templatesPath, theme, "partial", "info.tpl"),
		path.Join(templatesPath, theme, "partial", "pagination.tpl"),
		path.Join(templatesPath, theme, "partial", "posts.tpl"),
	}

	instance := &HTML{logger: zap.NewNop()}

	// default about panel
	partials, err := instance.getPartials(templatesPath, theme, "")
	require.NoError(t, err)
	assert.Equal(t, append(expected, path.Join(templatesPath, theme, "partial", "about.tpl")), partials)

	// custom about panel
	partials, err = instance.getPartials(templatesPath, theme, path.Join(templatesPath, theme, "partial", "pagination.tpl"))
	require.NoError(t, err)
	assert.Equal(t, append(expected, path.Join(templatesPath, theme, "partial", "pagination.tpl")), partials)
}
