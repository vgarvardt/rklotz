package model

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestNewMetaCalculated(t *testing.T) {
	m1 := NewMeta(5, 10)
	assert.Equal(t, m1.Posts, 5)
	assert.Equal(t, m1.PerPage, 10)
	assert.Equal(t, m1.Pages, 1)

	m2 := NewMeta(0, 10)
	assert.Equal(t, m2.Posts, 0)
	assert.Equal(t, m2.PerPage, 10)
	assert.Equal(t, m2.Pages, 0)

	m3 := NewMeta(55, 10)
	assert.Equal(t, m3.Posts, 55)
	assert.Equal(t, m3.PerPage, 10)
	assert.Equal(t, m3.Pages, 6)
}
