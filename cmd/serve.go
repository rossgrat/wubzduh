package cmd

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/rossgrat/wubzduh/internal/config"
	"github.com/rossgrat/wubzduh/internal/db"
	"github.com/rossgrat/wubzduh/internal/service"
	"github.com/rossgrat/wubzduh/internal/worker"
	"github.com/rossgrat/wubzduh/plugins/logger"
	spotifyplugin "github.com/rossgrat/wubzduh/plugins/spotify"
	"github.com/rossgrat/wubzduh/web"
	"github.com/spf13/cobra"
)

var serveConfigPath string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	RunE:  runServe,
}

func init() {
	serveCmd.Flags().StringVar(&serveConfigPath, "config", "./config.yaml", "path to config file")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(serveConfigPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log := logger.New(cfg.Log)

	store, err := db.NewStore(cfg.DB.ConnStr(), db.WithLogger(log))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer store.Close()

	spotifyClient, err := spotifyplugin.NewClient(cfg.Spotify.ClientID, cfg.Spotify.ClientSecret)
	if err != nil {
		return fmt.Errorf("connecting to spotify: %w", err)
	}

	fetchSvc := service.NewFetchService(store, spotifyClient, service.FetchWithLogger(log))
	purgeSvc := service.NewPurgeService(store, service.PurgeWithLogger(log))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fetchWorker := worker.New("fetch", func() error {
		return fetchSvc.Fetch(true)
	}, 12*time.Hour, worker.WithLogger(log))

	purgeWorker := worker.New("purge", func() error {
		return purgeSvc.Purge()
	}, 24*time.Hour, worker.WithLogger(log))

	go fetchWorker.Run(ctx)
	go purgeWorker.Run(ctx)

	srv := web.New(store, web.WithPort(cfg.Server.Port), web.WithLogger(log))
	return srv.Run(ctx)
}
