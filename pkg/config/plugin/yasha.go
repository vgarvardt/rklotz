package plugin

type YandexShare struct{}

func (p *YandexShare) Defaults() map[string]string {
	return map[string]string{}
}

func (p *YandexShare) Configure(settings map[string]string) (map[string]string, error) {
	return nil, ErrorConfiguring
}
