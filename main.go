package main

import (
	"flag"
	"fmt"
	"gym_management_system/api"
	"gym_management_system/db"
	"log"
)

func main() {
	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

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
