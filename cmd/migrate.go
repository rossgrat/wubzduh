package cmd

import (
	"fmt"

	"github.com/rossgrat/wubzduh/internal/config"
	"github.com/rossgrat/wubzduh/internal/db"
	"github.com/spf13/cobra"
)

var migrateConfigPath string

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	RunE:  runMigrateUp,
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back all migrations",
	RunE:  runMigrateDown,
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	RunE:  runMigrateVersion,
}

func init() {
	migrateCmd.PersistentFlags().StringVar(&migrateConfigPath, "config", "./config.yaml", "path to config file")
	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd, migrateVersionCmd)
	rootCmd.AddCommand(migrateCmd)
}

func runMigrateUp(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(migrateConfigPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if err := db.RunMigrateUp(cfg.DB.ConnStr()); err != nil {
		return fmt.Errorf("running migrations up: %w", err)
	}

	fmt.Println("Migrations applied successfully.")
	return nil
}

func runMigrateDown(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(migrateConfigPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if err := db.RunMigrateDown(cfg.DB.ConnStr()); err != nil {
		return fmt.Errorf("running migrations down: %w", err)
	}

	fmt.Println("Migrations rolled back successfully.")
	return nil
}

func runMigrateVersion(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(migrateConfigPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	version, dirty, err := db.GetMigrateVersion(cfg.DB.ConnStr())
	if err != nil {
		return fmt.Errorf("getting migration version: %w", err)
	}

	fmt.Printf("Version: %d, Dirty: %v\n", version, dirty)
	return nil
}
