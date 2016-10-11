package model

import (
	"fmt"
	"os"
	"path"

	"github.com/boltdb/bolt"
	log "github.com/Sirupsen/logrus"

	"github.com/vgarvardt/rklotz/cfg"
	"github.com/vgarvardt/rklotz/svc"
)

var db *bolt.DB

func init() {
	logger := svc.Container.MustGet(svc.DI_LOGGER).(*log.Logger)

	dbPath := fmt.Sprintf("%s/%s", cfg.GetRootDir(), cfg.String("db.path"))
	logger.WithField("path", dbPath).Info("Openning DB")
	os.MkdirAll(path.Dir(dbPath), 0755)
	db, _ = bolt.Open(dbPath, 0644, nil)
}

func GetDB() *bolt.DB {
	return db
}
