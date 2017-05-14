package plugin

// YandexMetrika is https://metrika.yandex.ru/ Plugin implementation
type YandexMetrika struct{}

func (p *YandexMetrika) Defaults() map[string]string {
	return map[string]string{"clickmap": "true", "trackLinks": "true", "accurateTrackBounce": "true"}
}

func (p *YandexMetrika) Configure(settings map[string]string) (map[string]string, error) {
	err := validateRequiredFields(settings, []string{"id"})
	if nil != err {
		return nil, err
	}

	return mergeSettings(settings, p.Defaults()), nil
}
