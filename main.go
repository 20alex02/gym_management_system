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

	now := time.Now()
	start := now.AddDate(0, 0, 1-int(now.Weekday()))
	from := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	queryValues := url.Values{}
	queryValues.Set("from", from.Format(time.RFC3339))
	queryValues.Set("to", to.Format(time.RFC3339))
	fmt.Println(queryValues.Encode())

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
		if err = store.Seed(); err != nil {
			log.Fatal(err)
		}
	}

	validator, err := api.NewCustomValidator()
	if err != nil {
		log.Fatal(err)
	}
	//v := validator.New()
	//if err := v.RegisterValidation("password", api.PasswordValidation); err != nil {
	//	log.Fatal(err)
	//}
	//if err := v.RegisterValidation("gteCurrentDay", api.GreaterThanOrEqualCurrentDayValidation); err != nil {
	//	log.Fatal(err)
	//}
	//if err := v.RegisterValidation("gtNow", api.GreaterThanNowValidation); err != nil {
	//	log.Fatal(err)
	//}
	server := api.NewServer(":3000", store, validator)
	server.Run()
}
