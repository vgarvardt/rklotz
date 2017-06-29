package plugin

// GoogleAnalytics is http://www.google.com/analytics/ Plugin implementation
type GoogleAnalytics struct{}

func (p *GoogleAnalytics) Defaults() map[string]string {
	return map[string]string{}
}

func (p *GoogleAnalytics) Configure(settings map[string]string) (map[string]string, error) {
	err := validateRequiredFields(settings, []string{"tracking_id"})
	if nil != err {
		return nil, err
	}

	return mergeSettings(settings, p.Defaults()), nil
}
