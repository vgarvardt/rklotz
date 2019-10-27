package web

// HTTPConfig is the configuration for web application
type HTTPConfig struct {
	Port          int    `envconfig:"WEB_PORT" default:"8080"`
	StaticPath    string `envconfig:"WEB_STATIC_PATH" default:"/etc/rklotz/static"`
	TemplatesPath string `envconfig:"WEB_TEMPLATES_PATH" default:"/etc/rklotz/templates"`
}

// SSLConfig is the configuration for TLS/SSL
type SSLConfig struct {
	Enabled  bool   `envconfig:"SSL_ENABLED" default:"false"`
	Port     int    `envconfig:"SSL_PORT" default:"8443"`
	Host     string `envconfig:"SSL_HOST"`
	Email    string `envconfig:"SSL_EMAIL" default:"vgarvardt@gmail.com"`
	CacheDir string `envconfig:"SSL_CACHE_DIR" default:"/tmp"`
}
