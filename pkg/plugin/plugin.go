package plugin

import "errors"

var (
	// ErrorUnknownPlugin is the error returned when trying to get unknown plugin
	ErrorUnknownPlugin = errors.New("Unknown plugin")
)

var all map[string]Plugin

// Plugin is the interface for plugins
type Plugin interface {
	// Defaults returns maps of default plugin configurations
	Defaults() map[string]string
	// SetUp applies settings map to a plugin
	SetUp(settings map[string]string) (map[string]string, error)
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
	all["disqus"] = new(Disqus)
	all["ga"] = new(GoogleAnalytics)
	all["gtm"] = new(GoogleTagManager)
	all["yamka"] = new(YandexMetrika)
	all["highlightjs"] = new(HighlightJS)
	all["yasha"] = new(YandexShare)
}

func validateRequiredFields(settings map[string]string, fields []string) error {
	for _, field := range fields {
		if _, ok := settings[field]; !ok {
			return NewErrorConfiguring(field)
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
