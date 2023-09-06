package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adarsh-a-tw/tt-backend/cli"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Cannot connect to database: %s", err.Error())
	}
	defer db.Close()

	app := cli.New(db)

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
