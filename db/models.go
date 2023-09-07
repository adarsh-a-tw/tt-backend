package db

type Match struct {
	Id        int    `db:"id"`
	Stage     string `db:"stage"`
	Format    string `db:"format"`
	GamePoint int    `db:"game_point"`
	SetCount  int    `db:"set_count"`
	Status    string `db:"status"`
}

type Set struct {
	Id             int  `db:"id"`
	SetNumber      int  `db:"set_number"`
	MatchId        int  `db:"match_id"`
	OpponentAScore int  `db:"opp_a_score"`
	OpponentBScore int  `db:"opp_b_score"`
	IsCompleted    bool `db:"is_completed"`
}

type Team struct {
	Id      int    `db:"id"`
	PlayerA string `db:"player_a"`
	PlayerB string `db:"player_b"`
}

type Player struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type TeamMatchMapping struct {
	MatchId     int  `db:"match_id"`
	TeamId      int  `db:"team_id"`
	IsOpponentA bool `db:"is_opp_a"`
	IsWinner    bool `db:"is_winner"`
}

type PlayerMatchMapping struct {
	MatchId     int  `db:"match_id"`
	PlayerId    int  `db:"player_id"`
	IsOpponentA bool `db:"is_opp_a"`
	IsWinner    bool `db:"is_winner"`
}

type SetLog struct {
	Id        int  `db:"id"`
	SetId     int  `db:"set_id"`
	OppAScore int  `db:"opp_a_score"`
	OppBScore int  `db:"opp_b_score"`
	ScoredByA bool `db:"scored_by_a"`
}
