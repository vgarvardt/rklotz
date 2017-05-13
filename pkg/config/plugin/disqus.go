package plugin

type Disqus struct{}

func (p *Disqus) Defaults() map[string]string {
	return map[string]string{}
}

func (p *Disqus) Configure(settings map[string]string) (map[string]string, error) {
	return nil, ErrorConfiguring
}
