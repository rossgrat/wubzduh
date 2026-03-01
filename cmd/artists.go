package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rossgrat/wubzduh/internal/config"
	"github.com/rossgrat/wubzduh/internal/db"
	"github.com/rossgrat/wubzduh/plugins/logger"
	spotifyplugin "github.com/rossgrat/wubzduh/plugins/spotify"
	"github.com/spf13/cobra"
)

var artistsConfigPath string

var artistsCmd = &cobra.Command{
	Use:   "artists",
	Short: "Manage tracked artists",
}

var artistsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tracked artists",
	RunE:  runArtistsList,
}

var artistsAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Search Spotify and add an artist",
	Args:  cobra.ExactArgs(1),
	RunE:  runArtistsAdd,
}

var artistsSearchCmd = &cobra.Command{
	Use:   "search [name]",
	Short: "Search for an artist in the database",
	Args:  cobra.ExactArgs(1),
	RunE:  runArtistsSearch,
}

func init() {
	artistsCmd.PersistentFlags().StringVar(&artistsConfigPath, "config", "./config.yaml", "path to config file")
	artistsCmd.AddCommand(artistsListCmd, artistsAddCmd, artistsSearchCmd)
	rootCmd.AddCommand(artistsCmd)
}

func runArtistsList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(artistsConfigPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log := logger.New(cfg.Log)
	store, err := db.NewStore(cfg.DB.ConnStr(), db.WithLogger(log))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer store.Close()

	artists, err := store.GetAllArtists()
	if err != nil {
		return fmt.Errorf("getting artists: %w", err)
	}

	if len(artists) == 0 {
		fmt.Println("No artists found.")
		return nil
	}

	for _, a := range artists {
		fmt.Printf("  %s (spotify: %s)\n", a.Name, a.SpotifyID)
	}
	return nil
}

func runArtistsAdd(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(artistsConfigPath)
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

	name := args[0]
	results, err := spotifyClient.SearchArtists(name, 5)
	if err != nil {
		return fmt.Errorf("searching spotify: %w", err)
	}

	if results.Artists == nil || len(results.Artists.Artists) == 0 {
		fmt.Println("No artists found on Spotify.")
		return nil
	}

	fmt.Println("Search results:")
	for i, a := range results.Artists.Artists {
		genres := "none"
		if len(a.Genres) > 0 {
			genres = strings.Join(a.Genres, ", ")
		}
		fmt.Printf("  [%d] %s (popularity: %d, genres: %s)\n", i, a.Name, a.Popularity, genres)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nSelect artist number (or 'q' to cancel): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" {
		return nil
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 0 || idx >= len(results.Artists.Artists) {
		return fmt.Errorf("invalid selection: %s", input)
	}

	selected := results.Artists.Artists[idx]

	existing, err := store.GetArtistsByName(selected.Name)
	if err != nil {
		return fmt.Errorf("checking existing artists: %w", err)
	}
	if len(existing) > 0 {
		fmt.Printf("Artist %q already exists in database.\n", selected.Name)
		return nil
	}

	artist := db.Artist{
		Name:      selected.Name,
		SpotifyID: selected.ID.String(),
	}
	if err := store.InsertArtist(artist); err != nil {
		return fmt.Errorf("inserting artist: %w", err)
	}

	fmt.Printf("Added %q to tracked artists.\n", selected.Name)
	return nil
}

func runArtistsSearch(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(artistsConfigPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	log := logger.New(cfg.Log)
	store, err := db.NewStore(cfg.DB.ConnStr(), db.WithLogger(log))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer store.Close()

	artists, err := store.GetArtistsByName(args[0])
	if err != nil {
		return fmt.Errorf("searching artists: %w", err)
	}

	if len(artists) == 0 {
		fmt.Printf("No artists found matching %q.\n", args[0])
		return nil
	}

	for _, a := range artists {
		fmt.Printf("  %s (spotify: %s)\n", a.Name, a.SpotifyID)
	}
	return nil
}
