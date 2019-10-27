package plugin

// YandexMetrika is https://metrika.yandex.ru/ Plugin implementation
type YandexMetrika struct{}

// Defaults returns maps of default plugin configurations
func (p *YandexMetrika) Defaults() map[string]string {
	return map[string]string{}
}

// SetUp applies settings map to a plugin
func (p *YandexMetrika) SetUp(settings map[string]string) (map[string]string, error) {
	err := validateRequiredFields(settings, []string{"id"})
	if nil != err {
		return nil, err
	}

	return mergeSettings(settings, p.Defaults()), nil
}
