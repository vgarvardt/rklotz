package model

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/boltdb/bolt"

	"../cfg"
)

var db *bolt.DB

func init() {
	// can not use cfg.GetRootDir() as it may be not initialized before this init
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	dbPath := fmt.Sprintf("%s/%s", rootDir, cfg.String("db.path"))
	cfg.Log(fmt.Sprintf("DB path: %s", dbPath))
	os.MkdirAll(path.Dir(dbPath), 0755)
	db, _ = bolt.Open(dbPath, 0644, nil)
}

func GetDB() *bolt.DB {
	return db
}
