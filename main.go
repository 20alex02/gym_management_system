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

//member := db.Account{FirstName: "john", LastName: "doe", Email: "john.doe@mail.com", EncryptedPassword: "asdf"}
//if _, errors := store.CreateAccount(&member); errors != nil {
//	log.Fatal(errors)
//}
//member.Email = "john@doe.com"
//member.FirstName = "john"
//member.LastName = "doe"
//if _, errors := store.CreateAccount(&member); errors != nil {
//	log.Fatal(errors)
//}
//accounts, errors := store.GetAllAccounts()
//if errors != nil {
//	log.Fatal(errors)
//}
//fmt.Println(accounts)
//account, errors := store.GetAccountById(3)
//if errors != nil {
//	log.Fatal(errors)
//}
//fmt.Println(account)
