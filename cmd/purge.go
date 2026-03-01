package cmd

import (
	"fmt"

	"github.com/rossgrat/wubzduh/internal/config"
	"github.com/rossgrat/wubzduh/internal/db"
	"github.com/rossgrat/wubzduh/internal/service"
	"github.com/rossgrat/wubzduh/plugins/logger"
	"github.com/spf13/cobra"
)

var purgeConfigPath string

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Manually purge releases older than 14 days",
	RunE:  runPurge,
}

func init() {
	purgeCmd.Flags().StringVar(&purgeConfigPath, "config", "./config.yaml", "path to config file")
	rootCmd.AddCommand(purgeCmd)
}

func runPurge(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(purgeConfigPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log := logger.New(cfg.Log)

	store, err := db.NewStore(cfg.DB.ConnStr(), db.WithLogger(log))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer store.Close()

	purgeSvc := service.NewPurgeService(store, service.PurgeWithLogger(log))

	log.Info("running manual purge")
	return purgeSvc.Purge()
}
