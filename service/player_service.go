package service

import (
	"github.com/adarsh-a-tw/tt-backend/db"
)

func (s *Service) CreatePlayer(name string) error {
	return s.repo.CreatePlayer(&db.Player{Name: name})
}
