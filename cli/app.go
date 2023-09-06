package cli

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	database "github.com/adarsh-a-tw/tt-backend/db"
	"github.com/adarsh-a-tw/tt-backend/enums"
	"github.com/adarsh-a-tw/tt-backend/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli"
)

var importChoices = []string{"player", "team", "match"}

func New(db *sqlx.DB) *cli.App {
	app := cli.NewApp()
	app.Name = "TT Backend"
	app.Usage = "Cli to run import fixture commands or run the backend server"
	app.Version = "1.0.0"

	registerCommands(app, db)

	// registerFlags(app)
	// registerCliActions(app, db)

	return app
}

func registerCommands(app *cli.App, db *sqlx.DB) {
	app.Commands = []cli.Command{
		{
			Name:        "serve",
			ShortName:   "s",
			Description: "Starts the server",
			Action: func(c *cli.Context) error {
				runServer(db)
				return nil
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
				fmt.Println(importType)
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

func runServer(db *sqlx.DB) {
	var port = 8080
	addr := fmt.Sprintf(":%d", port)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(addr)
}

func createPlayers(reader *csv.Reader, svc *service.Service) error {
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

func createTeams(reader *csv.Reader, svc *service.Service) error {
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

func createMatches(reader *csv.Reader, svc *service.Service) error {
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
