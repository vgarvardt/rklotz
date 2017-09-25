package plugin

// Disqus is https://disqus.com/ Plugin implementation
type Disqus struct{}

// Defaults returns maps of default plugin configurations
func (p *Disqus) Defaults() map[string]string {
	return map[string]string{}
}

// Configure applies settings map to a plugin
func (p *Disqus) Configure(settings map[string]string) (map[string]string, error) {
	err := validateRequiredFields(settings, []string{"shortname"})
	if nil != err {
		return nil, err
	}

	return mergeSettings(settings, p.Defaults()), nil
}
