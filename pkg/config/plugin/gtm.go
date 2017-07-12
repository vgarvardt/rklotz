package plugin

// GoogleAnalytics is https://tagmanager.google.com Plugin implementation
type GoogleTagManager struct{}

func (p *GoogleTagManager) Defaults() map[string]string {
	return map[string]string{}
}

func (p *GoogleTagManager) Configure(settings map[string]string) (map[string]string, error) {
	err := validateRequiredFields(settings, []string{"id"})
	if nil != err {
		return nil, err
	}

	return mergeSettings(settings, p.Defaults()), nil
}
