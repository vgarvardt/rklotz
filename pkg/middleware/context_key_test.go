package middleware

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContextKey_String(t *testing.T) {
	key := time.Now().Format(time.RFC3339Nano)
	ctxKey := contextKey(key)
	assert.True(t, strings.HasPrefix(ctxKey.String(), prefix))
	assert.True(t, strings.HasSuffix(ctxKey.String(), key))
}
