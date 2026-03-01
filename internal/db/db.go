package db

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq"
)

type Store struct {
	db     *sql.DB
	logger *slog.Logger
}

type Option func(*Store)

func WithLogger(log *slog.Logger) Option {
	return func(s *Store) { s.logger = log }
}

func NewStore(connStr string, opts ...Option) (*Store, error) {
	s := &Store{logger: slog.Default()}
	for _, opt := range opts {
		opt(s)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	s.db = db

	if err := s.RunMigrations(connStr); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	s.logger.Info("database connected")
	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
