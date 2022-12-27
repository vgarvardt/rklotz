package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByName(t *testing.T) {
	p, err := GetByName("ga")
	require.NoError(t, err)
	assert.IsType(t, &GoogleAnalytics{}, p)

	_, err = GetByName("!!!")
	require.Error(t, err)
	assert.Equal(t, ErrorUnknownPlugin, err)
}

type mockPlugin struct{}

func (p *mockPlugin) Defaults() map[string]string {
	return map[string]string{}
}

func (p *mockPlugin) SetUp(settings map[string]string) (map[string]string, error) {
	return mergeSettings(settings, p.Defaults()), nil
}

func TestGetName(t *testing.T) {
	p, err := GetByName("ga")
	require.NoError(t, err)
	assert.IsType(t, &GoogleAnalytics{}, p)

	name, err := GetName(p)
	require.NoError(t, err)
	assert.Equal(t, "ga", name)

	_, err = GetName(&mockPlugin{})
	require.Error(t, err)
	assert.Equal(t, ErrorUnknownPlugin, err)
}

func TestGetAll(t *testing.T) {
	assert.Equal(t, all, GetAll())
	assert.True(t, len(GetAll()) > 0)
}
