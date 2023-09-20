package main

import (
	"fmt"
	"log"
	"os"

	"github.com/adarsh-a-tw/tt-backend/cli"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Cannot connect to database: %s", err.Error())
	}
	defer db.Close()

	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	app := cli.New(db, rdb)

	err = app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
