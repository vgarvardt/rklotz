package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYandexShare_Configure(t *testing.T) {
	p := &YandexShare{}
	_, err := p.Configure(map[string]string{})
	require.NoError(t, err)

	settings, err := p.Configure(map[string]string{"lang": "de"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"services": "facebook,twitter,gplus", "size": "m", "lang": "de"}, settings)

	settings, err = p.Configure(map[string]string{"services": "facebook twitter"})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"services": "facebook,twitter", "size": "m", "lang": "en"}, settings)
}
