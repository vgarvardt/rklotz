package web

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"

	"github.com/cappuccinotm/slogx"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pkgz/routegroup"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"

	"github.com/vgarvardt/rklotz/pkg/server/handler"
	m "github.com/vgarvardt/rklotz/pkg/server/middleware"
)

const staticPath = "/static"

// NewRouter initialises and builds new HTTP router
func NewRouter(ph *handler.Posts, fh *handler.Feed, cfgHTTP HTTPConfig, theme string, logger *slog.Logger) http.Handler {
	mux := routegroup.New(http.NewServeMux())

	mux.Use(
		middleware.RequestID,
		middleware.RealIP,
		m.NewLogger(logger).Handler,
		m.NewRequestLogger().Handler,
		m.Recovery,
	)

	mux.HandleFunc("GET /{$}", ph.Front)
	mux.HandleFunc("GET /tag/{tag}", ph.Tag)
	mux.HandleFunc("/", ph.Post)

	mux.Mount("/feed").Route(func(b *routegroup.Bundle) {
		b.HandleFunc("GET /atom", fh.Atom)
		b.HandleFunc("GET /rss", fh.Rss)
	})

	staticRoot := http.Dir(cfgHTTP.StaticPath)
	staticHandler := http.StripPrefix(staticPath, http.FileServer(staticRoot))
	faviconPath := filepath.Join(cfgHTTP.StaticPath, theme, "favicon.ico")

	mux.HandleFunc("GET "+staticPath+"/", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})

	return mux
}

// ListenAndServe launches web server that listens to HTTP(S) requests
func ListenAndServe(ctx context.Context, router http.Handler, cfgSSL SSLConfig, cfgHTTP HTTPConfig, logger *slog.Logger) error {
	if !cfgSSL.Enabled {
		server := &http.Server{
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      10 * time.Second,
			Addr:              fmt.Sprintf(":%d", cfgHTTP.Port),
			Handler:           router,
		}

		logger.Info("Running HTTP server...", slog.String("address", server.Addr))
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
		Handler:           router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				cert, err := tlsCertManager.GetCertificate(info)
				if err != nil {
					logger.Error(
						"TLS cert manager could not get certificate",
						slogx.Error(err),
						slog.String("server-name", info.ServerName),
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
		logger.Info("Running HTTPS server...", slog.String("address", httpsServer.Addr))
		if err := httpsServer.ListenAndServeTLS("", ""); err != nil {
			return fmt.Errorf("failed to run HTTPS server: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		logger.Info("Running HTTP to HTTPS redirect server...", slog.String("address", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil {
			return fmt.Errorf("failed to run HTTPS redirect server: %w", err)
		}
		return nil
	})

	<-ctx.Done()
	logger.Info("One of the servers stopped, stopping all of them")
	logger.Info("Stopping HTTPS server", slogx.Error(httpsServer.Shutdown(context.Background())))
	logger.Info("Stopping HTTP to HTTPS redirect server", slogx.Error(httpServer.Shutdown(context.Background())))

	return g.Wait()
}
