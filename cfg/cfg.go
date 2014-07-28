package cfg

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IoSync/goptimist"
	"github.com/robfig/config"
)

var app map[string]interface{}
var reader *config.Config
var env string
var stdlogger *log.Logger

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

func init() {
	app = goptimist.CliApp(os.Args)
	if val, ok := app["env"]; !ok {
		env = "prod"
	} else {
		env = fmt.Sprintf("%v", val)
	}

	stdlogger = log.New(os.Stdout, "", 0)
	Log(fmt.Sprintf("Loading config with env set to %s", env))

	filePath := fmt.Sprintf("./cfg/%s.ini", env)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		Log(fmt.Sprintf("Loading config from ./cfg/config.ini"))
		reader, _ = config.ReadDefault("./cfg/config.ini")
	} else {
		Log(fmt.Sprintf("Loading config from %s", filePath))
		reader, _ = config.ReadDefault(filePath)
	}
}
