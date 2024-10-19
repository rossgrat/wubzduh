package main

// This module allows a user to add a module to the database via the command line

import (
	"context"
	"fmt"
	"log"

	"github.com/rossgrat/wubzduh/src/lib/db"
	"github.com/rossgrat/wubzduh/src/lib/threads"
	"github.com/rossgrat/wubzduh/src/lib/util"

	_ "github.com/lib/pq"

	"github.com/zmb3/spotify/v2"
)

func addArtist(client *spotify.Client, ctx context.Context) {
	for {
		var artistName string
		fmt.Printf("\nArtist Addition \nUse '_' instead of ' '\nUse 'exit' to return to menu\n\nEnter Artist: ")
		fmt.Scan(&artistName)
		if artistName == "exit" {
			fmt.Printf("Exiting CLI")
			return
		}
		fmt.Printf("Entered Artist %s\n\n", artistName)

		results, err := client.Search(ctx, artistName, spotify.SearchTypeArtist, spotify.Limit(5))
		if err != nil {
			log.Fatal(err)
		}
		if results.Artists == nil {
			fmt.Printf("Did not find any artists.")
			continue
		}
		fmt.Printf("Found %d Artists.\n", len(results.Artists.Artists))

		var searchMap [10]spotify.FullArtist
		var counter int = 0
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

		var selection int
		fmt.Printf("Select entry: ")
		fmt.Scan(&selection)
		fmt.Printf("Selected entry %d: %s\n", selection, searchMap[selection].Name)

		artists, err := db.GetAllArtistsWithName(searchMap[selection].Name)
		if err != nil {
			log.Fatal(err.Error())
		}
		if len(artists) > 0 {
			fmt.Printf("Error - Artist exists in database.")
			return
		}

		artist := db.Artist{
			Name:      searchMap[selection].Name,
			SpotifyID: searchMap[selection].SimpleArtist.ID.String(),
		}
		if err := db.InsertArtist(artist); err != nil {
			log.Fatalf("couldn't insert artist: %v", err)
		}
	}
}

func showArtists() {
	artists, err := db.GetAllArtists()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("\n\n")
	fmt.Printf("---Results---\n")
	for _, artist := range artists {
		fmt.Printf("%s %s\n", artist.Name, artist.SpotifyID)
	}
}

func searchForArtist() {
	for {
		var artistName string
		fmt.Printf("\nArtist Search \nUse '_' instead of ' '\nUse 'exit' to return to menu\n\nEnter Artist: ")
		fmt.Scan(&artistName)
		if artistName == "exit" {
			fmt.Printf("Exiting CLI")
			return
		}
		fmt.Printf("Entered Artist %s\n\n", artistName)

		artists, err := db.GetAllArtistsWithName(artistName)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("\n\n")
		fmt.Printf("---Results---\n")
		for _, artist := range artists {
			fmt.Printf("%s %s\n", artist.Name, artist.SpotifyID)
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
	db.Connect()
	client, ctx := util.ConnectToSpotify()

	for {
		option := optionMenu()
		switch option {
		case "exit":
			fmt.Printf("Exiting.\n")
			return
		case "1":
			addArtist(client, ctx)
		case "2":
			showArtists()
		case "3":
			searchForArtist()
		case "4":
			fmt.Print("Executing Fetch...")
			threads.Fetch(true)
		case "5":
			fmt.Print("Executing Fetch, no date matching...")
			threads.Fetch(false)
		case "6":
			fmt.Print("Executing Purge..")
			threads.Purge()
		}
	}
}
