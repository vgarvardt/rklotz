package plugin

// GoogleAnalytics is http://www.google.com/analytics/ Plugin implementation
type GoogleAnalytics struct{}

// Defaults returns maps of default plugin configurations
func (p *GoogleAnalytics) Defaults() map[string]string {
	return map[string]string{}
}

// SetUp applies settings map to a plugin
func (p *GoogleAnalytics) SetUp(settings map[string]string) (map[string]string, error) {
	err := validateRequiredFields(settings, []string{"tracking_id"})
	if nil != err {
		return nil, err
	}

	return mergeSettings(settings, p.Defaults()), nil
}
