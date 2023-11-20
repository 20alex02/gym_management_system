package main

import (
	"flag"
	"fmt"
	"gym_management_system/api"
	"gym_management_system/db"
	"log"
	"net/url"
	"time"
)

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	request := api.GetEventsRequest{
		From: time.Now().AddDate(0, 0, -1),
		To:   time.Now().AddDate(0, 0, 7),
	}

	// Convert the struct to a URL-encoded query string
	queryValues := url.Values{}
	queryValues.Set("from", request.From.Format(time.RFC3339))
	queryValues.Set("to", request.To.Format(time.RFC3339))

	// Append the query string to the base URL
	baseURL := "http://localhost:3000/api/events"
	fullURL := baseURL + "?" + queryValues.Encode()

	// Print the full URL with query parameters
	fmt.Println(fullURL)

	store, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close(store.Db)

	//if errors := store.Init(); errors != nil {
	//	log.Fatal(errors)
	//}

	if *seed {
		fmt.Println("seeding the database")
		if errors := store.Seed(); errors != nil {
			log.Fatal(errors)
		}
	}

	server := api.NewServer(":3000", store)
	server.Run()
}
