package plugin

// GoogleTagManager is https://tagmanager.google.com Plugin implementation
type GoogleTagManager struct{}

// Defaults returns maps of default plugin configurations
func (p *GoogleTagManager) Defaults() map[string]string {
	return map[string]string{}
}

// Configure applies settings map to a plugin
func (p *GoogleTagManager) Configure(settings map[string]string) (map[string]string, error) {
	err := validateRequiredFields(settings, []string{"id"})
	if nil != err {
		return nil, err
	}

	return mergeSettings(settings, p.Defaults()), nil
}
