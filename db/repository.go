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
	GetAllMatches(matches *[]Match, statusFilter string) error
	GetTeamInfoByMatchId(matchId int) ([]TeamInfoByMatchIdRow, error)
	GetPlayerInfoByMatchId(matchId int) ([]PlayerInfoByMatchIdRow, error)
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

func (r *repository) GetAllMatches(matches *[]Match, statusFilter string) error {
	query := `SELECT * FROM match`

	if statusFilter != "" {
		query += ` WHERE status = :filter`
	}

	query += ` ORDER BY id ASC`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err := stmt.Select(matches, map[string]interface{}{"filter": statusFilter}); err != nil {
		return err
	}

	return nil
}

type TeamInfoByMatchIdRow struct {
	MatchId     int    `db:"match_id"`
	Id          int    `db:"id"`
	TeamId      int    `db:"team_id"`
	PlayerA     string `db:"player_a"`
	PlayerB     string `db:"player_b"`
	IsOpponentA bool   `db:"is_opp_a"`
	IsWinner    bool   `db:"is_winner"`
}

func (r *repository) GetTeamInfoByMatchId(matchId int) ([]TeamInfoByMatchIdRow, error) {
	query := `SELECT * FROM team_match_mapping JOIN team ON team_match_mapping.team_id = team.id`
	query += ` WHERE match_id = :matchId`
	query += ` ORDER BY is_opp_a DESC`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows := []TeamInfoByMatchIdRow{}

	if err := stmt.Select(&rows, map[string]interface{}{"matchId": matchId}); err != nil {
		return nil, err
	}

	return rows, nil
}

type PlayerInfoByMatchIdRow struct {
	MatchId     int    `db:"match_id"`
	Id          int    `db:"id"`
	PlayerId    int    `db:"player_id"`
	PlayerName  string `db:"name"`
	IsOpponentA bool   `db:"is_opp_a"`
	IsWinner    bool   `db:"is_winner"`
}

func (r *repository) GetPlayerInfoByMatchId(matchId int) ([]PlayerInfoByMatchIdRow, error) {
	query := `SELECT * FROM player_match_mapping JOIN player ON player_match_mapping.player_id = player.id`
	query += ` WHERE match_id = :matchId`
	query += ` ORDER BY is_opp_a DESC`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows := []PlayerInfoByMatchIdRow{}

	if err := stmt.Select(&rows, map[string]interface{}{"matchId": matchId}); err != nil {
		return nil, err
	}

	return rows, nil
}
