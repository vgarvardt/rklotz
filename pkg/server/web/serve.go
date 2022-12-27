package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

	"github.com/vgarvardt/rklotz/pkg/server/handler"
	m "github.com/vgarvardt/rklotz/pkg/server/middleware"
)

// NewRouter initialises and builds new HTTP router
func NewRouter(pH *handler.Posts, fH *handler.Feed, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(m.NewLogger(logger).Handler)
	r.Use(m.NewRequestLogger().Handler)
	r.Use(m.Recovery)

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
func ServeStatic(r chi.Router, cfgHTTP HTTPConfig, theme string) {
	staticRoot := http.Dir(cfgHTTP.StaticPath)
	staticPath := "/static"
	staticHandler := http.StripPrefix(staticPath, http.FileServer(staticRoot))

	faviconPath := filepath.Join(cfgHTTP.StaticPath, theme, "favicon.ico")

	r.Get(staticPath+"/*", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})
}

// ListenAndServe launches web server that listens to HTTP(S) requests
func ListenAndServe(ctx context.Context, handler chi.Router, cfgSSL SSLConfig, cfgHTTP HTTPConfig, logger *zap.Logger) error {
	if !cfgSSL.Enabled {
		server := &http.Server{
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      10 * time.Second,
			Addr:              fmt.Sprintf(":%d", cfgHTTP.Port),
			Handler:           handler,
		}

		logger.Info("Running HTTP server...", zap.String("address", server.Addr))
		return server.ListenAndServe()
	}

	logger.Info("SSLConfig is enabled, starting HTTPS server")

	tlsCertManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfgSSL.Host),
		Cache:      autocert.DirCache(cfgSSL.CacheDir),
		Email:      cfgSSL.Email,
	}

	httpsServer := &http.Server{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		Addr:              fmt.Sprintf(":%d", cfgSSL.Port),
		Handler:           handler,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				cert, err := tlsCertManager.GetCertificate(info)
				if err != nil {
					logger.Error(
						"TLS cert manager could not get certificate",
						zap.Error(err),
						zap.String("server-name", info.ServerName),
					)
				}

				return cert, err
			},
		},
	}

	httpServer := &http.Server{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		Addr:              fmt.Sprintf(":%d", cfgHTTP.Port),
		Handler:           tlsCertManager.HTTPHandler(nil),
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		logger.Info("Running HTTPS server...", zap.String("address", httpsServer.Addr))
		if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
			return fmt.Errorf("failed to run HTTPS server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		logger.Info("Running HTTP to HTTPS redirect server...", zap.String("address", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil {
			return fmt.Errorf("failed to run HTTPS redirect server: %w", err)
		}
		return nil
	})

	<-ctx.Done()
	logger.Info("One of the servers stopped, stopping all of them")
	logger.Info("Stopping HTTPS server", zap.Error(httpsServer.Shutdown(context.Background())))
	logger.Info("Stopping HTTP to HTTPS redirect server", zap.Error(httpServer.Shutdown(context.Background())))

	return g.Wait()
}
