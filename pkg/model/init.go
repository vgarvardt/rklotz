package model

import (
	"github.com/boltdb/bolt"
)

var DB *bolt.DB

//func init() {
//	dbPath := "/tmp/rklotz.db"
//	log.WithField("path", dbPath).Info("Openning DB")
//	os.MkdirAll(path.Dir(dbPath), 0755)
//
//	var err error
//	if DB, err = bolt.Open(dbPath, 0644, nil); err != nil {
//		panic(err)
//	}
//}
