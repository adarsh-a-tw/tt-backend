package db

import (
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateMatch(match *Match) (int64, error)
	CreatePlayer(player *Player) error
	CreateTeam(team *Team) error
	CreateSet(set *Set) (int64, error)
	UpdateSet(set *Set) error
	AddTeamToMatch(mapping *TeamMatchMapping) error
	AddPlayerToMatch(mapping *PlayerMatchMapping) error
	UpdateMatchWinner(match *Match, isOppA bool) error
	ResetMatchWinner(match *Match) error
	GetAllMatches(matches *[]Match, statusFilter string) error
	GetMatchById(id int) (*Match, error)
	GetSetsByMatchId(id int) ([]Set, error)
	GetTeamInfoByMatchId(matchId int) ([]TeamInfoByMatchIdRow, error)
	GetPlayerInfoByMatchId(matchId int) ([]PlayerInfoByMatchIdRow, error)
	UpdateMatchStatus(matchId int, status string) error
	CreateSetLog(setLog *SetLog) error
	DeleteSetLog(id int) error
	GetSetLogsBySetId(setId int, limit *int) ([]SetLog, error)
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

func (r *repository) CreateSet(set *Set) (int64, error) {
	query := `
		INSERT INTO set (set_number, match_id, opp_a_score, opp_b_score, is_completed)
		VALUES (:set_number, :match_id, :opp_a_score, :opp_b_score, :is_completed)
		RETURNING id;
	`

	var id int64 // Declare a variable to store the returned ID

	// Use r.db.QueryRow to execute the query and return a single row
	rows, err := r.db.NamedQuery(query, set)
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

func (r *repository) UpdateSet(set *Set) error {
	query := `
		UPDATE set 
		SET set_number = :set_number, match_id = :match_id,
		opp_a_score = :opp_a_score, opp_b_score = :opp_b_score,
		is_completed = :is_completed
		WHERE id = :id;
	`

	_, err := r.db.NamedExec(query, set)

	return err
}

func (r *repository) GetMatchById(id int) (*Match, error) {
	query := `
		SELECT * FROM match WHERE id = $1;
	`
	var match Match
	err := r.db.Get(&match, query, id)
	if err != nil {
		return nil, err
	}

	return &match, nil
}

func (r *repository) GetSetsByMatchId(id int) ([]Set, error) {
	query := `
		SELECT * FROM set WHERE match_id = :id ORDER BY set_number ASC;
	`

	rows, err := r.db.NamedQuery(query, map[string]interface{}{"id": id})
	if err != nil {
		return nil, err
	}
	var sets []Set
	for rows.Next() {
		var set Set
		err = rows.StructScan(&set)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}

	return sets, nil
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

func (r *repository) UpdateMatchWinner(match *Match, isOppA bool) error {
	return r.writeMatchWinnerStatus(match, isOppA, true)
}

func (r *repository) ResetMatchWinner(match *Match) error {
	err := r.writeMatchWinnerStatus(match, true, false)
	if err != nil {
		return err
	}
	return r.writeMatchWinnerStatus(match, false, false)
}

func (r *repository) writeMatchWinnerStatus(match *Match, isOppA bool, isWinner bool) error {
	query := ``
	if match.Format == "SINGLES" {
		query = `
			UPDATE player_match_mapping SET is_winner = :isWinner
			WHERE match_id = :match_id AND is_opp_a = :is_opp_a
		`
	} else {
		query = `
			UPDATE team_match_mapping SET is_winner = :isWinner
			WHERE match_id = :match_id AND is_opp_a = :is_opp_a
		`
	}
	_, err := r.db.NamedExec(
		query,
		map[string]interface{}{"match_id": match.Id, "is_opp_a": isOppA, "isWinner": isWinner},
	)
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

func (r *repository) UpdateMatchStatus(matchId int, status string) error {
	query := `UPDATE match SET status = :status WHERE id = :matchId`

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(map[string]interface{}{"matchId": matchId, "status": status})

	return err
}

func (r *repository) CreateSetLog(setLog *SetLog) error {
	query := `
		INSERT INTO set_log (set_id, opp_a_score, opp_b_score, scored_by_a)
		VALUES (:set_id, :opp_a_score, :opp_b_score, :scored_by_a);
	`

	_, err := r.db.NamedExec(query, setLog)

	return err
}

func (r *repository) DeleteSetLog(id int) error {
	query := `
		DELETE FROM set_log WHERE id = :id;
	`

	_, err := r.db.NamedExec(query, map[string]interface{}{"id": id})

	return err
}

func (r *repository) GetSetLogsBySetId(setId int, limit *int) ([]SetLog, error) {
	var query string
	params := make(map[string]interface{}, 0)
	params["setId"] = setId
	if limit != nil {
		query = `
			SELECT * FROM set_log WHERE set_id = :setId ORDER BY id DESC LIMIT :limit;
		`
		params["limit"] = *limit
	} else {
		query = `
			SELECT * FROM set_log WHERE set_id = :setId ORDER BY id DESC;
		`
	}

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	setLogs := []SetLog{}

	if err := stmt.Select(&setLogs, params); err != nil {
		return nil, err
	}

	return setLogs, nil
}
