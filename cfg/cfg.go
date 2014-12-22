package cfg

import (
	"fmt"
	"log"
	"os"
	"time"
	"crypto/md5"
	"encoding/hex"

	"github.com/IoSync/goptimist"
	"github.com/robfig/config"
)

const VERSION = "0.2"

var app map[string]interface{}
var reader *config.Config
var env string
var stdlogger *log.Logger
var instanceId string

func String(key string) string {
	val, _ := reader.String(env, key)
	return val
}

func Int(key string) int {
	val, _ := reader.Int(env, key)
	return val
}

func Bool(key string) bool {
	val, _ := reader.Bool(env, key)
	return val
}

func Log(msg string) {
	stdlogger.Printf("[LOG] %v | %s\n", time.Now().Format("2006/01/02 - 15:04:05"), msg)
}

func GetApp() map[string]interface{} {
	return app
}

func GetInstanceId() string {
	return instanceId
}

func GetVersion() string {
	return VERSION
}

func init() {
	app = goptimist.CliApp(os.Args)
	if val, ok := app["env"]; !ok {
		env = "prod"
	} else {
		env = fmt.Sprintf("%v", val)
	}

	stdlogger = log.New(os.Stdout, "", 0)
	Log(fmt.Sprintf("Loading config with env set to %s", env))

	filePath := fmt.Sprintf("./%s.ini", env)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Log(fmt.Sprintf("Loading config from ./config.ini"))
		reader, _ = config.ReadDefault("./config.ini")
	} else {
		Log(fmt.Sprintf("Loading config from %s", filePath))
		reader, _ = config.ReadDefault(filePath)
	}

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format("2006/01/02 - 15:04:05")))
	instanceId = hex.EncodeToString(hasher.Sum(nil))[:5]
	Log(fmt.Sprintf("Initialized App Instance ID %s", instanceId))
}
