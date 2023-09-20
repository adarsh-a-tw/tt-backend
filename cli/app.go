package cli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/adarsh-a-tw/tt-backend/api"
	database "github.com/adarsh-a-tw/tt-backend/db"
	"github.com/adarsh-a-tw/tt-backend/enums"
	"github.com/adarsh-a-tw/tt-backend/service"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli"
)

var importChoices = []string{"player", "team", "match"}

func New(db *sqlx.DB, rdb *redis.Client) *cli.App {
	app := cli.NewApp()
	app.Name = "TT Backend"
	app.Usage = "Cli to run import fixture commands or run the backend server"
	app.Version = "1.0.0"

	registerCommands(app, db, rdb)

	return app
}

func registerCommands(app *cli.App, db *sqlx.DB, rdb *redis.Client) {
	app.Commands = []cli.Command{
		{
			Name:        "serve",
			ShortName:   "s",
			Description: "Starts the server",
			Action: func(c *cli.Context) error {
				return runServer(db, rdb)
			},
		},
		{
			Name:        "import",
			ShortName:   "i",
			Description: "Import resources to db",
			Action: func(c *cli.Context) error {
				csvFileName := c.String("csv")
				if csvFileName == "" {
					return fmt.Errorf("CSV file not specified")
				}

				importType := c.String("import-type")
				if importType == "" {
					return fmt.Errorf("import type not specified")
				}

				csvFile, err := os.Open(csvFileName)
				if err != nil {
					return err
				}
				defer csvFile.Close()

				reader := csv.NewReader(csvFile)
				svc := service.NewService(database.NewRepository(db))

				switch importType {
				case "player":
					return createPlayers(reader, svc)
				case "team":
					return createTeams(reader, svc)
				case "match":
					return createMatches(reader, svc)
				default:
					log.Println("Unknown resource type")
				}

				return nil
			},
			Flags: []cli.Flag{
				NewChoiceFlag(cli.StringFlag{
					Name:  "import-type",
					Usage: "Type of resource to import",
				}, importChoices,
				),
				cli.StringFlag{
					Name:  "csv",
					Usage: "CSV file containing import data",
				},
			},
		},
	}
}

func runServer(db *sqlx.DB, rdb *redis.Client) error {
	var port = 8080
	addr := fmt.Sprintf(":%d", port)

	svc := service.NewService(database.NewRepository(db))
	api := api.New(svc, rdb)
	return api.Serve(addr)
}

func createPlayers(reader *csv.Reader, svc service.Service) error {
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	keys := map[string]int{}
	for i, record := range records {
		if i == 0 {
			for j, value := range record {
				keys[value] = j
			}
			continue
		}
		nameIndex, ok := keys["name"]
		if !ok {
			return errors.New("field not found in csv: name")
		}

		err = svc.CreatePlayer(record[nameIndex])
		if err != nil {
			return err
		}
	}
	log.Println("Data imported successfully.")
	return nil
}

func createTeams(reader *csv.Reader, svc service.Service) error {
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	keys := map[string]int{}
	for i, record := range records {
		if i == 0 {
			for j, value := range record {
				keys[value] = j
			}
			continue
		}
		playerAIndex, ok := keys["player_a"]
		if !ok {
			return errors.New("field not found in csv: player_a")
		}
		playerBIndex, ok := keys["player_b"]
		if !ok {
			return errors.New("field not found in csv: player_b")
		}

		err = svc.CreateTeam(record[playerAIndex], record[playerBIndex])
		if err != nil {
			return err
		}
	}
	log.Println("Data imported successfully.")
	return nil
}

func createMatches(reader *csv.Reader, svc service.Service) error {
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	keys := map[string]int{}
	for i, record := range records {
		if i == 0 {
			for j, value := range record {
				keys[value] = j
			}

			keyNames := []string{"format", "stage", "opp_a_id", "opp_b_id", "max_sets", "game_point"}

			for _, key := range keyNames {
				_, ok := keys[key]
				if !ok {
					return fmt.Errorf("field not found in csv: %s", key)
				}
			}
			continue
		}

		format := enums.MatchFormat(record[keys["format"]])
		stage := enums.MatchStage(record[keys["stage"]])
		opp_a_id, err := strconv.Atoi(record[keys["opp_a_id"]])
		if err != nil {
			return err
		}
		opp_b_id, err := strconv.Atoi(record[keys["opp_b_id"]])
		if err != nil {
			return err
		}
		max_sets, err := strconv.Atoi(record[keys["max_sets"]])
		if err != nil {
			return err
		}
		game_point, err := strconv.Atoi(record[keys["game_point"]])
		if err != nil {
			return err
		}

		switch format {
		case enums.Singles:
			err = svc.CreateSinglesMatch(stage, opp_a_id, opp_b_id, max_sets, game_point)
			if err != nil {
				return err
			}
		case enums.Doubles:
			err = svc.CreateDoublesMatch(stage, opp_a_id, opp_b_id, max_sets, game_point)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid match format: %s", format)
		}
	}
	log.Println("Data imported successfully.")
	return nil
}
