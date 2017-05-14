package plugin

import (
	"io/ioutil"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGoogleAnalytics_Configure(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	p := &GoogleAnalytics{}
	_, err := p.Configure(map[string]string{})
	assert.Error(t, err)
	assert.Equal(t, err, ErrorConfiguring)

	settings, err := p.Configure(map[string]string{"tracking_id": "foo"})
	assert.NoError(t, err)
	assert.Equal(t, map[string]string{"tracking_id": "foo"}, settings)
}
