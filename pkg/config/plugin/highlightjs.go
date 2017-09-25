package plugin

// HighlightJS is https://highlightjs.org/ Plugin implementation
type HighlightJS struct{}

// Defaults returns maps of default plugin configurations
func (p *HighlightJS) Defaults() map[string]string {
	return map[string]string{"version": "9.7.0", "theme": "idea"}
}

// Configure applies settings map to a plugin
func (p *HighlightJS) Configure(settings map[string]string) (map[string]string, error) {
	return mergeSettings(settings, p.Defaults()), nil
}
