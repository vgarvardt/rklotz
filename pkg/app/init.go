package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/vgarvardt/rklotz/pkg/svc"
	"gopkg.in/alecthomas/kingpin.v2"
)

var version = "0.0.0-dev"

var (
	rootDir = kingpin.Flag("root", "Force set root dir").Default(".").String()

	_ = kingpin.Version(version)
)

type UpdateParams struct {
	UUID  string
	Field string
	Value string
}

var instanceId string

func InstanceId() string {
	return instanceId
}

func Version() string {
	return version
}

func RootDir() string {
	return *rootDir
}

func RootUrl(r *http.Request) *url.URL {
	config := svc.Container.MustGet(svc.DI_CONFIG).(svc.Config)

	scheme := config.String("ui.root_url.scheme")
	host := config.String("ui.root_url.host")
	path := config.String("ui.root_url.path")
	if len(host) < 1 {
		host = r.Host
	}
	return &url.URL{Scheme: scheme, Host: host, Path: path}
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	kingpin.Parse()

	log.WithField("version", Version()).Info("Initializing application")
	log.WithField("path", *rootDir).Info("Root dir parameter value")

	if *rootDir == "." {
		if rootDirPath, err := filepath.Abs(fmt.Sprintf("%s/../", filepath.Dir(os.Args[0]))); err != nil {
			panic(err)
		} else {
			rootDir = &rootDirPath
		}
		log.WithField("path", rootDir).Info("Root dir absolute path")
	}

	config := svc.NewIniEnvConfig(fmt.Sprintf("%s/var/config.ini", RootDir()), svc.ENV_PREFIX)
	svc.Container.Set(svc.DI_CONFIG, config)

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format("2006/01/02 - 15:04:05")))
	instanceId = hex.EncodeToString(hasher.Sum(nil))[:5]
	log.WithField("ID", instanceId).Info("Initialized application instance")
}
