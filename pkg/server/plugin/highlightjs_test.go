package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHighlightJS_Configure(t *testing.T) {
	p := &HighlightJS{}
	_, err := p.SetUp(map[string]string{})
	require.NoError(t, err)

	settings, err := p.SetUp(map[string]string{"theme": "foo"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"theme": "foo", "version": "9.15.10"}, settings)
}
