package plugin

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestYandexMetrika_Configure(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	p := &YandexMetrika{}
	_, err := p.Configure(map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, err, ErrorConfiguring)

	settings, err := p.Configure(map[string]string{"id": "foo"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"id": "foo"}, settings)
}
