package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoogleAnalytics_Configure(t *testing.T) {
	p := &GoogleAnalytics{}
	_, err := p.Configure(map[string]string{})
	require.Error(t, err)
	assert.IsType(t, &ErrorConfiguring{}, err)

	settings, err := p.Configure(map[string]string{"tracking_id": "foo"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"tracking_id": "foo"}, settings)
}
