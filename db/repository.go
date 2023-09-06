package db

import (
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateMatch(match *Match) (int64, error)
	CreatePlayer(player *Player) error
	CreateTeam(team *Team) error
	AddTeamToMatch(mapping *TeamMatchMapping) error
	AddPlayerToMatch(mapping *PlayerMatchMapping) error
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateMatch(match *Match) (int64, error) {
	query := `
		INSERT INTO match (stage, format, game_point, set_count, status)
		VALUES (:stage, :format, :game_point, :set_count, :status)
		RETURNING id;
	`

	var id int64 // Declare a variable to store the returned ID

	// Use r.db.QueryRow to execute the query and return a single row
	rows, err := r.db.NamedQuery(query, match)
	if err != nil {
		return 0, err
	}
	rows.Next()
	err = rows.Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil

}

func (r *repository) CreatePlayer(player *Player) error {
	query := `
		INSERT INTO player (name)
		VALUES (:name)
		RETURNING id
	`

	_, err := r.db.NamedExec(query, player)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreateTeam(team *Team) error {
	query := `
		INSERT INTO team (player_a, player_b)
		VALUES (:player_a, :player_b)
		RETURNING id
	`
	_, err := r.db.NamedExec(query, team)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) AddTeamToMatch(mapping *TeamMatchMapping) error {
	query := `
		INSERT INTO team_match_mapping (match_id, team_id, is_opp_a, is_winner)
		VALUES (:match_id, :team_id, :is_opp_a, :is_winner)
	`
	_, err := r.db.NamedExec(query, mapping)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) AddPlayerToMatch(mapping *PlayerMatchMapping) error {
	query := `
		INSERT INTO player_match_mapping (match_id, player_id, is_opp_a, is_winner)
		VALUES (:match_id, :player_id, :is_opp_a, :is_winner)
	`
	_, err := r.db.NamedExec(query, mapping)
	if err != nil {
		return err
	}
	return nil
}
