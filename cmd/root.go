package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wubzduh",
	Short: "New music release feed from hand-curated artists",
	Long:  "A service that monitors Spotify for new music releases from configured artists and serves a web feed.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
