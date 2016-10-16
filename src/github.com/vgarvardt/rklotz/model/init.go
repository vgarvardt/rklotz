package model

import (
	"fmt"
	"os"
	"path"

	"github.com/boltdb/bolt"
	log "github.com/Sirupsen/logrus"

	"github.com/vgarvardt/rklotz/app"
	"github.com/vgarvardt/rklotz/svc"
)

var DB *bolt.DB

func init() {
	logger := svc.Container.MustGet(svc.DI_LOGGER).(*log.Logger)
	config := svc.Container.MustGet(svc.DI_CONFIG).(svc.Config)

	dbPath := fmt.Sprintf("%s/%s", app.RootDir(), config.String("db.path"))
	logger.WithField("path", dbPath).Info("Openning DB")
	os.MkdirAll(path.Dir(dbPath), 0755)

	var err error
	if DB, err = bolt.Open(dbPath, 0644, nil); err != nil {
		panic(err)
	}
}
