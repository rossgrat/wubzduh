package main

// This module allows a user to add a module to the database via the command line
// TODO: Abstract database connect, abstract spotify connection

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	wubzduh "github.com/rossgrat/wubzduh/src"

	_ "github.com/lib/pq"

	"github.com/zmb3/spotify/v2"
)

const (
	insertArtist           = `INSERT INTO Artists (ArtistName, SpotifyID) values ($1, $2)`
	getAllArtistsNameAndID = `SELECT ArtistName, SpotifyID FROM Artists ORDER BY ArtistName ASC;`
	getAllArtistsWithName  = `SELECT ArtistName, SpotifyID FROM Artists WHERE ArtistName=$1`
)

func addArtist(db *sql.DB, client *spotify.Client, ctx context.Context) {
	var artist string
	for artist != "exit" {
		fmt.Printf("\nArtist Addition \nUse '_' instead of ' '\nUse 'exit' to return to menu\n\nEnter Artist: ")
		fmt.Scan(&artist)
		if artist == "exit" {
			fmt.Printf("Exiting CLI")
			return
		}
		fmt.Printf("Entered Artist %s\n\n", artist)

		results, err := client.Search(ctx, artist, spotify.SearchTypeArtist, spotify.Limit(5))
		if err != nil {
			log.Fatal(err)
		}
		if results.Artists != nil {
			fmt.Printf("Found %d Artists.\n", len(results.Artists.Artists))
		}

		var searchMap [10]spotify.FullArtist
		var counter int = 0
		// handle album results
		if results.Artists != nil {
			fmt.Printf("---Results---\n")
			for _, item := range results.Artists.Artists {
				searchMap[counter] = item
				fmt.Printf("%d: %s", counter, item.Name)
				fmt.Printf("\n\t %d ", item.Popularity)
				for _, genre := range item.Genres {
					fmt.Printf("%s ", genre)
				}
				fmt.Printf("\n")
				counter++
			}
		}

		var selection int
		fmt.Printf("Select entry: ")
		fmt.Scanf("%d", selection)
		fmt.Printf("Selected entry %d: %s\n", selection, searchMap[selection].Name)

		_, err = db.Exec(insertArtist, searchMap[selection].Name, searchMap[selection].SimpleArtist.ID)
		if err != nil {
			log.Fatalf("couldn't insert artist: %v", err)
		}
	}
}

func showArtists(db *sql.DB) {
	rows, err := db.Query(getAllArtistsNameAndID)
	if err != nil {
		log.Fatalf("Failed to query artists: %v", err)
	}
	var name, spotify_id string
	fmt.Printf("\n\n")
	fmt.Printf("---Results---\n")
	for rows.Next() {
		err = rows.Scan(&name, &spotify_id)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("%s %s\n", name, spotify_id)
	}
}

func searchForArtist(db *sql.DB) {
	var artist string
	for artist != "exit" {
		fmt.Printf("\nArtist Search \nUse '_' instead of ' '\nUse 'exit' to return to menu\n\nEnter Artist: ")
		fmt.Scan(&artist)
		if artist == "exit" {
			fmt.Printf("Exiting CLI")
			return
		}
		fmt.Printf("Entered Artist %s\n\n", artist)

		rows, err := db.Query(getAllArtistsWithName, artist)
		if err != nil {
			log.Fatalf("Failed to query artists: %v", err)
		}
		var name, spotify_id string
		fmt.Printf("---Results---\n")
		for rows.Next() {
			err = rows.Scan(&name, &spotify_id)
			if err != nil {
				log.Fatalf("Failed to scan row: %v", err)
			}
			fmt.Printf("%s %s\n", name, spotify_id)
		}
	}
}

func optionMenu() (option string) {
	fmt.Printf("\n\nOptions:\n")
	fmt.Printf("\t'exit': Exit CLI\n")
	fmt.Printf("\t1: Add Artists to DB\n")
	fmt.Printf("\t2: Show All Artists in DB\n")
	fmt.Printf("\t3: Search For Artist in DB\n")
	fmt.Printf("\t4: Run Fetch\n")
	fmt.Printf("\t5: Run Fetch without Current Date Matching\n")
	fmt.Printf("\t6: Run Purge\n")

	fmt.Printf("\nSelect Option: ")
	fmt.Scan(&option)
	return option
}

func main() {
	db := wubzduh.ConnectToDB()
	client, ctx := wubzduh.ConnectToSpotify()

	for {
		option := optionMenu()
		switch option {
		case "exit":
			fmt.Printf("Exiting.\n")
			return
		case "1":
			addArtist(db, client, ctx)
		case "2":
			showArtists(db)
		case "3":
			searchForArtist(db)
		case "4":
			fmt.Print("Executing Fetch...")
			wubzduh.Fetch(db, client, ctx)
		case "5":
			fmt.Print("Executing Fetch...")
			wubzduh.FetchAllLatest(db, client, ctx)
		case "6":
			fmt.Print("Executing Purge..")
			wubzduh.Purge(db)
		}
	}
}
