package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisqus_Configure(t *testing.T) {
	p := &Disqus{}
	_, err := p.SetUp(map[string]string{})
	require.Error(t, err)
	assert.IsType(t, &ErrorConfiguring{}, err)

	settings, err := p.SetUp(map[string]string{"shortname": "foo"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"shortname": "foo"}, settings)
}
