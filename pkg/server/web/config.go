package web

// HTTPConfig is the configuration for web application
type HTTPConfig struct {
	Port          int    `env:"WEB_PORT,default=8080"`
	StaticPath    string `env:"WEB_STATIC_PATH,default=/etc/rklotz/static"`
	TemplatesPath string `env:"WEB_TEMPLATES_PATH,default=/etc/rklotz/templates"`
}

// SSLConfig is the configuration for TLS/SSL
type SSLConfig struct {
	Enabled  bool   `env:"SSL_ENABLED,default=false"`
	Port     int    `env:"SSL_PORT,default=8443"`
	Host     string `env:"SSL_HOST"`
	Email    string `env:"SSL_EMAIL,default=vgarvardt@gmail.com"`
	CacheDir string `env:"SSL_CACHE_DIR,default=/tmp"`
}
