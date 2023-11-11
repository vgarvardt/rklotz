package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	require.NoError(t, err)

	dbFilePath := fmt.Sprintf("/tmp/rklotz-test.%s.db", hex.EncodeToString(hasher.Sum(nil))[:5])

	boltDBStorage, err := NewStorage("boltdb://"+dbFilePath, 10)
	assert.NoError(t, err)
	assert.IsType(t, &BoltDBStorage{}, boltDBStorage)
	assert.Equal(t, dbFilePath, boltDBStorage.(*BoltDBStorage).path)
	defer boltDBStorage.Close()

	_, err = NewStorage("unknown://", 10)
	assert.Error(t, err)
	assert.Equal(t, ErrorUnknownStorageType, err)

	_, err = NewStorage("~", 10)
	assert.Error(t, err)
}
