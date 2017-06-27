package plugin

import (
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestYandexShare_Configure(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	p := &YandexShare{}
	_, err := p.Configure(map[string]string{})
	assert.NoError(t, err)

	settings, err := p.Configure(map[string]string{"lang": "de"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"services": "facebook,twitter,gplus", "size": "m", "lang": "de"}, settings)

	settings, err = p.Configure(map[string]string{"services": "facebook twitter"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"services": "facebook,twitter", "size": "m", "lang": "en"}, settings)
}
