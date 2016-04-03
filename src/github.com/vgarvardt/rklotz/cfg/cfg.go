package cfg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/ini.v1"
)

const VERSION = "0.3.3"
const (
	COMMAND_RUN = "run"
	COMMAND_REBUILD = "rebuild"
	COMMAND_UPDATE = "update"
)

var (
	env = kingpin.Flag("env", "Application environment, defines config").Default("prod").String()
	rootDir = kingpin.Flag("root", "Force set root dir").Default(".").String()

	run = kingpin.Command(COMMAND_RUN, "Run application").Default()

	rebuild = kingpin.Command(COMMAND_REBUILD, "Rebuild index only, do not run web server")

	update = kingpin.Command(COMMAND_UPDATE, "Update post UUID field with new value")
	updateUUID = update.Arg("UUID", "Post UUID").Required().String()
	updateField = update.Arg("field", "Post field name").Required().String()
	updateValue = update.Arg("value", "Post field value").Required().String()

	version = kingpin.Version(VERSION)
)

type UpdateParams struct {
	UUID  string
	Field string
	Value string
}

var config *ini.File
var command string
var stdLogger *log.Logger
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

func Log(msg string) {
	stdLogger.Printf("[LOG] %v | %s\n", time.Now().Format("2006/01/02 - 15:04:05"), msg)
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

	stdLogger = log.New(os.Stdout, "", 0)
	Log(fmt.Sprintf("Initializing application ver %s", GetVersion()))

	Log(fmt.Sprintf("Root dir parameter value: %s", *rootDir))
	if *rootDir == "." {
		if rootDirPath, err := filepath.Abs(fmt.Sprintf("%s/../", filepath.Dir(os.Args[0]))); err != nil {
			panic(err)
		} else {
			rootDir = &rootDirPath
		}
		Log(fmt.Sprintf("Root dir absolute path: %s", *rootDir))
	}

	var err error
	configPath := fmt.Sprintf("%s/config.ini", GetRootDir())
	Log(fmt.Sprintf("Loading base config from %s", configPath))
	if config, err = ini.Load(configPath); err != nil {
		panic(err)
	}

	configPath = fmt.Sprintf("%s/%s.ini", GetRootDir(), *env)
	Log(fmt.Sprintf("Loading env config from %s", configPath))
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		Log("Env config not found")
	} else {
		if err := config.Append(configPath); err != nil {
			panic(err)
		}
	}

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format("2006/01/02 - 15:04:05")))
	instanceId = hex.EncodeToString(hasher.Sum(nil))[:5]
	Log(fmt.Sprintf("Initialized application instance ID %s", instanceId))
}
