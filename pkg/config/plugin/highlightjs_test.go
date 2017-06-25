package plugin

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestHighlightJS_Configure(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	p := &HighlightJS{}
	_, err := p.Configure(map[string]string{})
	assert.NoError(t, err)

	settings, err := p.Configure(map[string]string{"theme": "foo"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"theme": "foo", "version": "9.7.0"}, settings)
}
