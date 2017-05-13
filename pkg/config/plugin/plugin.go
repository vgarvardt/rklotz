package plugin

import "errors"

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
