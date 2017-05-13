package plugin

type HighlightJS struct{}

func (p *HighlightJS) Defaults() map[string]string {
	return map[string]string{}
}

func (p *HighlightJS) Configure(settings map[string]string) (map[string]string, error) {
	return nil, ErrorConfiguring
}
