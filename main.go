package main

import (
	"database/sql"
	"fmt"
	"gym_management_system/db"
	"log"
)

func main() {
	store, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(store.Db)

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
	member := db.Account{FirstName: "alex", LastName: "popovic", Email: "alexxx.popovic@gmail.com", EncryptedPassword: "asdf"}
	if _, err := store.CreateAccount(&member); err != nil {
		log.Fatal(err)
	}
	member.Email = "john@doe.com"
	member.FirstName = "john"
	member.LastName = "doe"
	if _, err := store.CreateAccount(&member); err != nil {
		log.Fatal(err)
	}
	accounts, err := store.GetAllAccounts()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(accounts)
	account, err := store.GetAccountById(3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(account)
}
