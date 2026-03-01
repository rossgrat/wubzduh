package web

import (
	"context"
	"embed"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/rossgrat/wubzduh/internal/db"
)

//go:embed templates/*.html
var embeddedTemplates embed.FS

//go:embed static/*
var embeddedStatic embed.FS

type Server struct {
	store     *db.Store
	templates *template.Template
	logger    *slog.Logger
	port      string
}

type Option func(*Server)

func WithPort(port string) Option {
	return func(s *Server) { s.port = port }
}

func WithLogger(log *slog.Logger) Option {
	return func(s *Server) { s.logger = log }
}

func New(store *db.Store, opts ...Option) *Server {
	s := &Server{
		store:  store,
		logger: slog.Default(),
		port:   "8080",
	}
	for _, opt := range opts {
		opt(s)
	}

	s.templates = template.Must(template.ParseFS(embeddedTemplates, "templates/*.html"))

	return s
}

func (s *Server) Run(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/artists/", s.artistsHandler)
	mux.HandleFunc("/feed/", s.feedHandler)
	mux.HandleFunc("/favicon.ico", s.faviconHandler)
	mux.HandleFunc("/", s.feedRedirect)
	mux.Handle("/static/", http.FileServerFS(embeddedStatic))

	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: s.requestLogger(mux),
	}

	go func() {
		<-ctx.Done()
		s.logger.Info("shutting down web server")
		srv.Shutdown(context.Background())
	}()

	s.logger.Info("starting web server", "port", s.port)
	return srv.ListenAndServe()
}
