package api

import (
	"autograder/internal/config"
	"autograder/internal/docker"
	"autograder/internal/grader"
	"context"
	"fmt"
	"net/http"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter"
	"github.com/ulule/limiter/drivers/middleware/stdlib"
)

// Server wraps the HTTP server with graceful shutdown.
type Server struct {
	http *http.Server
	wg   sync.WaitGroup
	log  *logrus.Logger
}

// Stop gracefully shuts down the server.
func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		if closeErr := s.http.Close(); closeErr != nil {
			s.log.WithError(closeErr).Error("force-close failed")
			return closeErr
		}
	}
	s.wg.Wait()
	s.log.Info("server stopped")
	return nil
}

// NewServer creates and starts the API server.
func NewServer(cfg *config.Config, log *logrus.Logger, dc docker.Client, g grader.Grader) (*Server, error) {
	configureLogger(log)

	h := &handler{
		cfg:    cfg,
		log:    log,
		docker: dc,
		grader: g,
	}

	store, rate, err := initializeRedis(log, &cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("redis init: %w", err)
	}
	mw := stdlib.NewMiddleware(limiter.New(*store, *rate))

	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.Use(mw.Handler)
	api.HandleFunc("/labs", h.listLabs).Methods("GET")
	api.HandleFunc("/submit", h.submit).Methods("POST")

	// Serve frontend (Vite build output)
	spa := &spaHandler{staticDir: "web/dist", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	srv := &Server{
		http: &http.Server{
			Addr:           addr,
			Handler:        router,
			ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		log: log,
	}

	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()
		log.WithField("addr", addr).Info("server started")
		if err := srv.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("server failed")
		}
	}()

	return srv, nil
}

func configureLogger(log *logrus.Logger) {
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02T15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			return s[len(s)-1], fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
	})
}

// spaHandler serves a single-page application, falling back to index.html.
type spaHandler struct {
	staticDir string
	indexPath string
}

func (s *spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := path.Join(s.staticDir, r.URL.Path)

	// If the file exists, serve it directly
	if info, err := http.Dir(s.staticDir).Open(r.URL.Path); err == nil {
		defer info.Close()
		http.FileServer(http.Dir(s.staticDir)).ServeHTTP(w, r)
		return
	}

	// Otherwise serve index.html for client-side routing
	http.ServeFile(w, r, path.Join(s.staticDir, s.indexPath))
	_ = p
}
