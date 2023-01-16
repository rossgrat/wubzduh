package wubzduh

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"golang.org/x/oauth2/clientcredentials"
)

const (
	host = "localhost"
	port = 5432
)

//Utility functions
func ConnectToDB() (db *sql.DB) {
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

//Connec to spotify client
func ConnectToSpotify() (client *spotify.Client, ctx context.Context) {
	ctx = context.Background()
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}
	token, err := config.Token(ctx)
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client = spotify.New(httpClient)
	return client, ctx
}
