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

	settings, err := p.Configure(map[string]string{"l10n": "de"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"services": "vkontakte,facebook,twitter,gplus", "type": "icon", "l10n": "de"}, settings)

	settings, err = p.Configure(map[string]string{"services": "facebook twitter"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"services": "facebook,twitter", "type": "icon", "l10n": "en"}, settings)
}
