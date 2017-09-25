package plugin

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrorUnknownPlugin is the error returned when trying to get unknown plugin
	ErrorUnknownPlugin = errors.New("Unknown plugin")
	// ErrorConfiguring is the error returned when plugin configuring fails
	ErrorConfiguring = errors.New("Failed to configure plugin")
)

var all map[string]Plugin

// Plugin is the interface for plugins
type Plugin interface {
	// Defaults returns maps of default plugin configurations
	Defaults() map[string]string
	// Configure applies settings map to a plugin
	Configure(settings map[string]string) (map[string]string, error)
}

// GetByName returns plugin instance by name
func GetByName(name string) (Plugin, error) {
	p, ok := all[name]
	if ok {
		return p, nil
	}
	return nil, ErrorUnknownPlugin
}

// GetName returns the name for a loaded plugin
func GetName(p Plugin) (string, error) {
	for name, instance := range all {
		if p == instance {
			return name, nil
		}
	}
	return "", ErrorUnknownPlugin
}

// GetAll returns loaded plugins map with the key as plugin name and value as plugin instance
func GetAll() map[string]Plugin {
	return all
}

func init() {
	all = make(map[string]Plugin)
	all["disqus"] = &Disqus{}
	all["ga"] = &GoogleAnalytics{}
	all["gtm"] = &GoogleTagManager{}
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
