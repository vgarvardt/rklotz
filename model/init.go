package model

import (
	"os"
	"path"
	"github.com/boltdb/bolt"

	"../cfg"
)

var db *bolt.DB

func init() {
	os.MkdirAll(path.Dir(cfg.String("db.path")), 0755)
	db, _ = bolt.Open(cfg.String("db.path"), 0644, nil)
}

func GetDB() *bolt.DB {
	return db
}
