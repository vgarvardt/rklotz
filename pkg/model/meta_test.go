package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetaCalculated(t *testing.T) {
	m1 := NewMeta(5, 10)
	assert.Equal(t, 5, m1.Posts)
	assert.Equal(t, 10, m1.PerPage)
	assert.Equal(t, 1, m1.Pages)

	m2 := NewMeta(0, 10)
	assert.Equal(t, 0, m2.Posts)
	assert.Equal(t, 10, m2.PerPage)
	assert.Equal(t, 0, m2.Pages)

	m3 := NewMeta(55, 10)
	assert.Equal(t, 55, m3.Posts)
	assert.Equal(t, 10, m3.PerPage)
	assert.Equal(t, 6, m3.Pages)
}
