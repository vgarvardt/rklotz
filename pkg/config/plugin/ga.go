package plugin

type GoogleAnalytics struct{}

func (p *GoogleAnalytics) Defaults() map[string]string {
	return map[string]string{}
}

func (p *GoogleAnalytics) Configure(settings map[string]string) (map[string]string, error) {
	return nil, ErrorConfiguring
}
