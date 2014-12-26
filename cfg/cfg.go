package cfg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/voxelbrain/goptions"
	"gopkg.in/ini.v0"
)

const VERSION = "0.1"

type Options struct {
	Env     string        `goptions:"-e, --env, description='<env> Application environment, defines config'"`
	Rebuild bool          `goptions:"-r, --rebuild, description='Rebuild index only, do not run web server'"`
	RootDir string        `goptions:"-d, --root, description='<dir> Force set root dir'"`
	Update  string        `goptions:"-u, --update, description='<uuid> Update post UUID field with new value (works with --field and --value set), do not run web server'"`
	Field   string        `goptions:"-f, --field, description='<field> Update post UUID field with new value (works with --update and --value set), do not run web server'"`
	Value   string        `goptions:"-v, --value, description='<value> Update post UUID field with new value (works with --update and --field set), do not run web server'"`
	Help    goptions.Help `goptions:"-h, --help, description='Show this help'"`
}

var config *ini.File
var options Options
var stdlogger *log.Logger
var instanceId string
var rootDir string
var runWebServer bool

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
	stdlogger.Printf("[LOG] %v | %s\n", time.Now().Format("2006/01/02 - 15:04:05"), msg)
}

func GetOptions() Options {
	return options
}

func GetInstanceId() string {
	return instanceId
}

func GetVersion() string {
	return VERSION
}

func GetRootDir() string {
	return rootDir
}

func GetRunWebServer() bool {
	return runWebServer
}

func init() {
	options = Options{
		Env:     "prod",
		Rebuild: false,
	}
	goptions.ParseAndFail(&options)

	runWebServer = !options.Rebuild && len(options.Update) < 1

	stdlogger = log.New(os.Stdout, "", 0)
	Log(fmt.Sprintf("Initializing application ver %s", VERSION))

	if len(options.RootDir) > 0 {
		rootDir = options.RootDir
	} else {
		var err error
		if rootDir, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			panic(err)
		}
		options.RootDir = rootDir
	}

	var err error
	configPath := fmt.Sprintf("%s/config.ini", GetRootDir())
	Log(fmt.Sprintf("Loading base config from %s", configPath))
	if config, err = ini.Load(configPath); err != nil {
		panic(err)
	}

	configPath = fmt.Sprintf("%s/%s.ini", GetRootDir(), options.Env)
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
