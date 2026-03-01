package cmd

import (
	"fmt"

	"github.com/rossgrat/wubzduh/internal/config"
	"github.com/rossgrat/wubzduh/internal/db"
	"github.com/rossgrat/wubzduh/internal/service"
	"github.com/rossgrat/wubzduh/plugins/logger"
	spotifyplugin "github.com/rossgrat/wubzduh/plugins/spotify"
	"github.com/spf13/cobra"
)

var fetchConfigPath string

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Manually fetch new releases",
	RunE:  runFetch,
}

func init() {
	fetchCmd.Flags().StringVar(&fetchConfigPath, "config", "./config.yaml", "path to config file")
	fetchCmd.Flags().Bool("no-date-check", false, "fetch all latest releases regardless of release date")
	rootCmd.AddCommand(fetchCmd)
}

func runFetch(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(fetchConfigPath)
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

	noDateCheck, _ := cmd.Flags().GetBool("no-date-check")
	fetchSvc := service.NewFetchService(store, spotifyClient, service.FetchWithLogger(log))

	log.Info("running manual fetch", "check_release_date", !noDateCheck)
	return fetchSvc.Fetch(!noDateCheck)
}
