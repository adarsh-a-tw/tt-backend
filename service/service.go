package service

import (
	"github.com/adarsh-a-tw/tt-backend/db"
	"github.com/adarsh-a-tw/tt-backend/enums"
)

type Service interface {
	CreateDoublesMatch(stage enums.MatchStage, teamAId int, teamBId int, maxSets int, gamePoint int) error
	CreatePlayer(name string) error
	CreateSinglesMatch(stage enums.MatchStage, playerAId int, playerBId int, maxSets int, gamePoint int) error
	CreateTeam(playerAName string, playerBName string) error
	GetMatchInfoList(status string) ([]matchInfo, error)
}

type service struct {
	repo db.Repository
}

func NewService(repo db.Repository) Service {
	return &service{repo: repo}
}
