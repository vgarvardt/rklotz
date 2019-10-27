package web

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"

	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/handler"
	m "github.com/vgarvardt/rklotz/pkg/middleware"
)

// NewRouter initialises and builds new HTTP router
func NewRouter(pH *handler.Posts, fH *handler.Feed, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(m.NewLogger(logger).Handler)
	r.Use(middleware.RequestLogger(m.NewRequestLogger()))
	r.Use(middleware.Recoverer)

	r.Get("/", pH.Front)
	r.Get("/tag/{tag}", pH.Tag)
	r.NotFound(pH.Post)

	r.Route("/feed", func(r chi.Router) {
		r.Get("/atom", fH.Atom)
		r.Get("/rss", fH.Rss)
	})

	return r
}

// ServeStatic registers static handler for the router
func ServeStatic(r chi.Router, cfgWeb config.Web, theme string) {
	staticRoot := http.Dir(cfgWeb.StaticPath)
	staticPath := "/static"
	staticHandler := http.StripPrefix(staticPath, http.FileServer(staticRoot))

	faviconPath := filepath.Join(cfgWeb.StaticPath, theme, "favicon.ico")

	r.Get(staticPath+"/*", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})
}

// ListenAndServe launches web server that listens to HTTP(S) requests
func ListenAndServe(handler chi.Router, cfgSSL config.SSL, cfgWeb config.Web, logger *zap.Logger) error {
	if !cfgSSL.Enabled {
		address := fmt.Sprintf(":%d", cfgWeb.Port)
		logger.Info("Running HTTP server...", zap.String("address", address))

		return http.ListenAndServe(address, handler)
	}

	logger.Info("SSL is enabled, starting HTTPS server")

	tlsCertManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfgSSL.Host),
		Cache:      autocert.DirCache(cfgSSL.CacheDir),
		Email:      cfgSSL.Email,
	}

	httpsAddress := fmt.Sprintf(":%d", cfgSSL.Port)
	server := &http.Server{
		Addr:    httpsAddress,
		Handler: handler,
		TLSConfig: &tls.Config{
			GetCertificate: tlsCertManager.GetCertificate,
		},
	}

	go func() {
		logger.Info("Running HTTPS server...", zap.String("address", httpsAddress))
		logger.Fatal("Failed to run HTTPS server", zap.Error(server.ListenAndServeTLS("", "")))
	}()

	httpAddress := fmt.Sprintf(":%d", cfgWeb.Port)

	logger.Info("Running HTTP to HTTPS redirect server...", zap.String("address", httpAddress))

	return http.ListenAndServe(httpAddress, tlsCertManager.HTTPHandler(nil))
}
