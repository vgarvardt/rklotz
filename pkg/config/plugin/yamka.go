package plugin

type YandexMetrika struct{}

func (p *YandexMetrika) Defaults() map[string]string {
	return map[string]string{}
}

func (p *YandexMetrika) Configure(settings map[string]string) (map[string]string, error) {
	return nil, ErrorConfiguring
}
