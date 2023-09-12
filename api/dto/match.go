package dto

type OpponentResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsWinner bool   `json:"is_winner"`
}

type MatchInfoResponse struct {
	Id        int                `json:"id"`
	Format    string             `json:"format"`
	Stage     string             `json:"stage"`
	Status    string             `json:"status"`
	Opponents []OpponentResponse `json:"opponents"`
}

type MatchSubscribeRequest struct {
	MatchId int `json:"match_id"`
}

type SetResponse struct {
	Id             int              `json:"id"`
	SetNumber      int              `json:"set_number"`
	OpponentAScore int              `json:"opp_a_score"`
	OpponentBScore int              `json:"opp_b_score"`
	IsCompleted    bool             `json:"is_completed"`
	Logs           []SetLogResponse `json:"logs"`
}

type SetLogResponse struct {
	Id        int  `json:"id"`
	OppAScore int  `json:"opp_a_score"`
	OppBScore int  `json:"opp_b_score"`
	ScoredByA bool `json:"scored_by_a"`
}

type MatchDetail struct {
	Id        int                `json:"id"`
	Format    string             `json:"format"`
	Stage     string             `json:"stage"`
	Status    string             `json:"status"`
	Opponents []OpponentResponse `json:"opponents"`
	Sets      []SetResponse      `json:"sets"`
}

type MatchDetailResponse struct {
	Data  MatchDetail `json:"data"`
	Error string          `json:"error"`
}
