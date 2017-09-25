package main

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vgarvardt/rklotz/pkg/config"
	"github.com/vgarvardt/rklotz/pkg/handler"
	"github.com/vgarvardt/rklotz/pkg/loader"
	m "github.com/vgarvardt/rklotz/pkg/middleware"
	"github.com/vgarvardt/rklotz/pkg/renderer"
	"github.com/vgarvardt/rklotz/pkg/storage"
	"golang.org/x/crypto/acme/autocert"
)

// RunServer initializes and runs web-server instance
func RunServer(cmd *cobra.Command, args []string) {
	appConfig, err := config.Load()
	failOnError(err, "Failed to load config")

	logLevel, err := log.ParseLevel(appConfig.LogLevel)
	failOnError(err, "Failed to parse log level")
	log.SetLevel(logLevel)

	hasher := md5.New()
	hasher.Write([]byte(time.Now().Format(time.RFC3339Nano)))
	instanceID := hex.EncodeToString(hasher.Sum(nil))[:5]

	log.WithFields(log.Fields{"version": version, "instance": instanceID}).Info("Starting rKlotz...")

	storageInstance, err := storage.NewStorage(appConfig.StorageDSN, appConfig.PostsPerPage)
	failOnError(err, "Failed to get storageInstance instance")
	defer storageInstance.Close()

	loaderInstance, err := loader.NewLoader(appConfig.PostsDSN)
	failOnError(err, "Failed to get loaderInstance instance")

	err = loaderInstance.Load(storageInstance)
	failOnError(err, "Failed to load posts")

	htmlRenderer, err := renderer.NewHTMLRenderer(appConfig.Web.TemplatesPath, instanceID, appConfig.UI, appConfig.Plugins, appConfig.RootURL)
	failOnError(err, "Failed to init HTML Renderer")
	xmlRenderer := renderer.NewXMLRenderer()

	postsHandler := handler.NewPostsHandler(storageInstance, htmlRenderer)
	feedHandler := handler.NewFeedHandler(storageInstance, xmlRenderer, appConfig.UI, appConfig.RootURL)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestLogger(&m.LoggerRequest{}))
	r.Use(middleware.Recoverer)

	r.Get("/", postsHandler.Front)
	r.Get("/tag/{tag}", postsHandler.Tag)
	r.NotFound(postsHandler.Post)

	r.Route("/feed", func(r chi.Router) {
		r.Get("/atom", feedHandler.Atom)
		r.Get("/rss", feedHandler.Rss)
	})

	serveStatic(r, appConfig)

	listenAndServe(r, appConfig)
}

func serveStatic(r chi.Router, appConfig *config.AppConfig) {
	staticRoot := http.Dir(appConfig.Web.StaticPath)
	staticPath := "/static"
	staticHandler := http.StripPrefix(staticPath, http.FileServer(staticRoot))

	faviconPath := filepath.Join(appConfig.Web.StaticPath, appConfig.UI.Theme, "favicon.ico")

	r.Get(staticPath+"/*", func(w http.ResponseWriter, r *http.Request) {
		staticHandler.ServeHTTP(w, r)
	})

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, faviconPath)
	})
}

func listenAndServe(r chi.Router, appConfig *config.AppConfig) {
	if appConfig.SSL.Enabled {
		log.Info("SSL is enabled, starting HTTPS server")

		hostPolicy := func(ctx context.Context, host string) error {
			if host == appConfig.SSL.Host {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", appConfig.SSL.Host)
		}

		tlsCertManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(appConfig.SSL.CacheDir),
			Email:      appConfig.UI.Email,
		}

		httpsAddress := fmt.Sprintf(":%d", appConfig.SSL.Port)
		server := &http.Server{
			Addr:      httpsAddress,
			Handler:   r,
			TLSConfig: &tls.Config{GetCertificate: tlsCertManager.GetCertificate},
		}

		go func() {
			log.WithField("address", httpsAddress).Info("Running HTTPS server...")
			log.Fatal(server.ListenAndServeTLS("", ""))
		}()
	}

	if appConfig.SSL.Enabled && appConfig.SSL.RedirectHTTP {
		mux := &http.ServeMux{}
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.WithField("url", r.URL.String()).Debug("Redirecting HTTP request to HTTPS")

			newURI := "https://" + r.Host + r.URL.String()
			http.Redirect(w, r, newURI, http.StatusFound)
		})

		httpAddress := fmt.Sprintf(":%d", appConfig.Web.Port)
		redirectServer := &http.Server{Handler: mux, Addr: httpAddress}

		log.WithField("address", httpAddress).Info("Running HTTP to HTTPS redirect server...")
		log.Fatal(redirectServer.ListenAndServe())
	} else {
		address := fmt.Sprintf(":%d", appConfig.Web.Port)
		log.WithField("address", address).Info("Running HTTP server...")

		log.Fatal(http.ListenAndServe(address, r))
	}
}

func failOnError(err error, msg string) {
	if nil != err {
		log.WithError(err).Panic(msg)
	}
}
