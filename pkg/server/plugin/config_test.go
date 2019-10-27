package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlugins_SetUp(t *testing.T) {
	p := Config{
		Settings: Settings{
			Disqus:      map[string]string{"shortname": "foo"},
			Ga:          map[string]string{"tracking_id": "foo"},
			Gtm:         map[string]string{"id": "foo"},
			Yamka:       map[string]string{"id": "foo"},
			Highlightjs: map[string]string{},
			Yasha:       map[string]string{},
		},
	}
	instance, _ := GetByName("ga")

	config, err := p.SetUp(instance)
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"tracking_id": "foo"}, config)

	_, err = p.SetUp(&mockPlugin{})
	require.Error(t, err)
	assert.Equal(t, ErrorUnknownPlugin, err)

	for _, instance := range GetAll() {
		_, err = p.SetUp(instance)
		assert.NoError(t, err)
	}
}
