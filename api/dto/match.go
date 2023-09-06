package dto

type OpponentResponse struct {
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
