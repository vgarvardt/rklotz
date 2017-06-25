package plugin

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDisqus_Configure(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	p := &Disqus{}
	_, err := p.Configure(map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, err, ErrorConfiguring)

	settings, err := p.Configure(map[string]string{"shortname": "foo"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"shortname": "foo"}, settings)
}
