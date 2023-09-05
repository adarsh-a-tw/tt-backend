package service

import (
	"github.com/adarsh-a-tw/tt-backend/db"
)

func (s *Service) CreateTeam(playerAName, playerBName string) error {
	return s.repo.CreateTeam(&db.Team{PlayerA: playerAName, PlayerB: playerBName})
}
