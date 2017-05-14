package plugin

import (
	"errors"

	log "github.com/Sirupsen/logrus"
)

var (
	ErrorUnknownPlugin = errors.New("Unknown plugin")
	ErrorConfiguring   = errors.New("Failed to configure plugin")
)

var all map[string]Plugin

type Plugin interface {
	Defaults() map[string]string
	Configure(settings map[string]string) (map[string]string, error)
}

func GetByName(name string) (Plugin, error) {
	p, ok := all[name]
	if ok {
		return p, nil
	}
	return nil, ErrorUnknownPlugin
}

func GetName(p Plugin) (string, error) {
	for name, instance := range all {
		if p == instance {
			return name, nil
		}
	}
	return "", ErrorUnknownPlugin
}

func GetAll() map[string]Plugin {
	return all
}

func init() {
	all = make(map[string]Plugin)
	all["disqus"] = &Disqus{}
	all["ga"] = &GoogleAnalytics{}
	all["yamka"] = &YandexMetrika{}
	all["highlightjs"] = &HighlightJS{}
	all["yasha"] = &YandexShare{}
}

func validateRequiredFields(settings map[string]string, fields []string) error {
	for _, field := range fields {
		if _, ok := settings[field]; !ok {
			log.WithField("field", field).Error("Required field missing")
			return ErrorConfiguring
		}
	}
	return nil
}

func mergeSettings(settings map[string]string, defaults map[string]string) map[string]string {
	for name, value := range defaults {
		if _, ok := settings[name]; !ok {
			settings[name] = value
		}
	}
	return settings
}
