package plugin

import (
	"errors"
	"reflect"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Config is the configuration for app plugins
type Config struct {
	Enabled []string `env:"PLUGINS_ENABLED"`
	Settings
}

// SetUp applies configuration for enabled plugins
func (c Config) SetUp(instance Plugin) (map[string]string, error) {
	pluginName, err := GetName(instance)
	if err != nil {
		return nil, err
	}

	pluginSettings, ok := c.Settings.Get(pluginName)
	if !ok {
		return nil, errors.New("failed to get plugin settings")
	}

	if len(pluginSettings) == 0 {
		return instance.SetUp(instance.Defaults())
	}

	return instance.SetUp(pluginSettings)
}

// Settings is the configuration for available plugins
type Settings struct {
	Disqus      map[string]string `env:"PLUGINS_DISQUS"`
	Ga          map[string]string `env:"PLUGINS_GA"`
	Gtm         map[string]string `env:"PLUGINS_GTM"`
	Yamka       map[string]string `env:"PLUGINS_YAMKA"`
	Highlightjs map[string]string `env:"PLUGINS_HIGHLIGHTJS"`
	Yasha       map[string]string `env:"PLUGINS_YASHA"`
}

// Get gets plugin settings by name
func (s Settings) Get(pluginName string) (map[string]string, bool) {
	r := reflect.ValueOf(s)
	f := reflect.Indirect(r).FieldByName(cases.Title(language.English).String(pluginName)).Interface()

	fMap, ok := f.(map[string]string)

	return fMap, ok
}
