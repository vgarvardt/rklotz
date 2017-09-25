package plugin

import "strings"

// YandexShare is https://tech.yandex.ru/share/ Plugin implementation
type YandexShare struct{}

// Defaults returns maps of default plugin configurations
func (p *YandexShare) Defaults() map[string]string {
	return map[string]string{"services": "facebook,twitter,gplus", "size": "m", "lang": "en"}
}

// Configure applies settings map to a plugin
func (p *YandexShare) Configure(settings map[string]string) (map[string]string, error) {
	// convert spaces to commas in services list as this is how this setting come from environment settings
	if val, ok := settings["services"]; ok {
		settings["services"] = strings.Replace(val, " ", ",", -1)
	}

	return mergeSettings(settings, p.Defaults()), nil
}
