package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoogleTagManager_Configure(t *testing.T) {
	p := &GoogleTagManager{}
	_, err := p.Configure(map[string]string{})
	require.Error(t, err)
	assert.IsType(t, &ErrorConfiguring{}, err)

	settings, err := p.Configure(map[string]string{"id": "foo"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"id": "foo"}, settings)
}
