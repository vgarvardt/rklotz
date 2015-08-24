package model

import (
	"fmt"
	"os"
	"path"

	"github.com/boltdb/bolt"

	"github.com/vgarvardt/rklotz/cfg"
)

var db *bolt.DB

func init() {
	dbPath := fmt.Sprintf("%s/%s", cfg.GetRootDir(), cfg.String("db.path"))
	cfg.Log(fmt.Sprintf("DB path: %s", dbPath))
	os.MkdirAll(path.Dir(dbPath), 0755)
	db, _ = bolt.Open(dbPath, 0644, nil)
}

func GetDB() *bolt.DB {
	return db
}
