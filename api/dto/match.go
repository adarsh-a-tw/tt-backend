package api

type Match struct {
	Id        int    `json:"id"`
	Stage     string `json:"stage"`
	Format    string `json:"format"`
	GamePoint int    `json:"game_point"`
	SetCount  int    `json:"set_count"`
	Status    string `json:"status"`
	Sets      []Set  `json:"sets"`
	OpponentA string `json:"opponent_A"`
}

type Set struct{}
