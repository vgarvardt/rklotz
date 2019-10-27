package plugin

// HighlightJS is https://highlightjs.org/ Plugin implementation
type HighlightJS struct{}

// Defaults returns maps of default plugin configurations
func (p *HighlightJS) Defaults() map[string]string {
	return map[string]string{"version": "9.15.10", "theme": "idea"}
}

// SetUp applies settings map to a plugin
func (p *HighlightJS) SetUp(settings map[string]string) (map[string]string, error) {
	return mergeSettings(settings, p.Defaults()), nil
}
