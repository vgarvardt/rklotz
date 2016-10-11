package cfg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/ini.v1"
	log "github.com/Sirupsen/logrus"

	"github.com/vgarvardt/rklotz/svc"
)

const VERSION = "0.3.5"
const (
	COMMAND_RUN = "run"
	COMMAND_REBUILD = "rebuild"
	COMMAND_UPDATE = "update"
)

var (
	env = kingpin.Flag("env", "Application environment, defines config").Default("prod").String()
	rootDir = kingpin.Flag("root", "Force set root dir").Default(".").String()

	_ = kingpin.Command(COMMAND_RUN, "Run application").Default()

	_ = kingpin.Command(COMMAND_REBUILD, "Rebuild index only, do not run web server")

	update = kingpin.Command(COMMAND_UPDATE, "Update post UUID field with new value")
	updateUUID = update.Arg("UUID", "Post UUID").Required().String()
	updateField = update.Arg("field", "Post field name").Required().String()
	updateValue = update.Arg("value", "Post field value").Required().String()

	_ = kingpin.Version(VERSION)
)

type UpdateParams struct {
	UUID  string
	Field string
	Value string
}

var config *ini.File
var command string
var instanceId string

func configKey(key string) *ini.Key {
	return config.Section("").Key(key)
}

func String(key string) string {
	return configKey(key).String()
}

func Int(key string) int {
	if val, err := configKey(key).Int(); err != nil {
		panic(err)
	} else {
		return val
	}
}

func Bool(key string) bool {
	if val, err := configKey(key).Bool(); err != nil {
		panic(err)
	} else {
		return val
	}
}

func GetInstanceId() string {
	return instanceId
}

func GetVersion() string {
	return VERSION
}

func GetRootDir() string {
	return *rootDir
}

func GetCommand() string {
	return command
}

func GetUpdateParams() UpdateParams {
	return UpdateParams{UUID: *updateUUID, Field: *updateField, Value: *updateValue}
}

func GetRootUrl(r *http.Request) *url.URL {
	scheme := String("ui.root_url.scheme")
	host := String("ui.root_url.host")
	path := String("ui.root_url.path")
	if len(host) < 1 {
		host = r.Host
	}
	return &url.URL{Scheme: scheme, Host: host, Path: path}
}

func init() {
	command = kingpin.Parse()

	logger := svc.Container.MustGet(svc.DI_LOGGER).(*log.Logger)

	logger.WithField("version", GetVersion()).Info("Initializing application")

	logger.WithField("path", *rootDir).Info("Root dir parameter value")
	if *rootDir == "." {
		if rootDirPath, err := filepath.Abs(fmt.Sprintf("%s/../", filepath.Dir(os.Args[0]))); err != nil {
			panic(err)
		} else {
			rootDir = &rootDirPath
		}
		logger.WithField("path", *rootDir).Info("Root dir absolute path")
	}

	var err error
	configPath := fmt.Sprintf("%s/config.ini", GetRootDir())
	logger.WithField("path", configPath).Info("Loading base config")
	if config, err = ini.Load(configPath); err != nil {
		panic(err)
	}

	configPath = fmt.Sprintf("%s/%s.ini", GetRootDir(), *env)
	logger.WithField("path", configPath).Info("Loading env config")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logger.WithField("path", configPath).Warn("Env config not found")
	} else {
		if err := config.Append(configPath); err != nil {
			logger.WithField("err", err).Fatal("Failed to append env config")
		}
	}

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format("2006/01/02 - 15:04:05")))
	instanceId = hex.EncodeToString(hasher.Sum(nil))[:5]
	logger.WithField("ID", instanceId).Info("Initialized application instance")
}
