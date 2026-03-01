package worker

import (
	"context"
	"log/slog"
	"time"
)

type TaskFunc func() error

type Worker struct {
	name     string
	task     TaskFunc
	interval time.Duration
	logger   *slog.Logger
}

type Option func(*Worker)

func WithLogger(log *slog.Logger) Option {
	return func(w *Worker) { w.logger = log }
}

func New(name string, task TaskFunc, interval time.Duration, opts ...Option) *Worker {
	w := &Worker{
		name:     name,
		task:     task,
		interval: interval,
		logger:   slog.Default(),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *Worker) Run(ctx context.Context) {
	w.logger.Info("worker started", "name", w.name, "interval", w.interval.String())

	if err := w.task(); err != nil {
		w.logger.Error("worker task failed", "name", w.name, "error", err)
	}

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker stopping", "name", w.name)
			return
		case <-ticker.C:
			if err := w.task(); err != nil {
				w.logger.Error("worker task failed", "name", w.name, "error", err)
			}
		}
	}
}
