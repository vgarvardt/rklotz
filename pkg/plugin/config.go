package plugin

import (
	"errors"
	"strings"

	"github.com/fatih/structs"
)

// Config is teh configuration for app plugins
type Config struct {
	Enabled  []string `envconfig:"PLUGINS_ENABLED"`
	Settings Settings
}

// SetUp applies configuration for enabled plugins
func (p Config) SetUp(instance Plugin) (map[string]string, error) {
	pluginName, err := GetName(instance)
	if err != nil {
		return nil, err
	}

	settingsMap := structs.Map(p.Settings)
	pluginSettings, ok := settingsMap[strings.Title(pluginName)]
	if !ok {
		return nil, errors.New("failed to get plugin settings")
	}

	pluginSettingsMap := pluginSettings.(map[string]string)
	if len(pluginSettingsMap) == 0 {
		return instance.SetUp(instance.Defaults())
	}

	return instance.SetUp(pluginSettingsMap)
}

// Settings is the configuration for available plugins
type Settings struct {
	Disqus      map[string]string `envconfig:"PLUGINS_DISQUS"`
	Ga          map[string]string `envconfig:"PLUGINS_GA"`
	Gtm         map[string]string `envconfig:"PLUGINS_GTM"`
	Yamka       map[string]string `envconfig:"PLUGINS_YAMKA"`
	Highlightjs map[string]string `envconfig:"PLUGINS_HIGHLIGHTJS"`
	Yasha       map[string]string `envconfig:"PLUGINS_YASHA"`
}
