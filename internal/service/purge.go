package service

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/rossgrat/wubzduh/internal/db"
)

type PurgeService struct {
	store  *db.Store
	logger *slog.Logger
}

type PurgeOption func(*PurgeService)

func PurgeWithLogger(log *slog.Logger) PurgeOption {
	return func(p *PurgeService) { p.logger = log }
}

func NewPurgeService(store *db.Store, opts ...PurgeOption) *PurgeService {
	p := &PurgeService{
		store:  store,
		logger: slog.Default(),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *PurgeService) Purge() error {
	cutoff := time.Now().AddDate(0, 0, -14)

	albums, err := p.store.GetAllAlbumReleaseDates()
	if err != nil {
		return fmt.Errorf("getting album release dates: %w", err)
	}

	purged := 0
	for _, album := range albums {
		if album.ReleaseDate.Before(cutoff) {
			if err := p.store.DeleteTracksForAlbum(album.ID); err != nil {
				p.logger.Error("failed to delete tracks for album", "album_id", album.ID, "error", err)
				continue
			}
			if err := p.store.DeleteAlbum(album.ID); err != nil {
				p.logger.Error("failed to delete album", "album_id", album.ID, "error", err)
				continue
			}
			purged++
		}
	}

	if purged > 0 {
		p.logger.Info("purged old albums", "count", purged)
	}
	return nil
}
