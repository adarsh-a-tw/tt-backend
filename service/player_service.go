package service

import (
	"github.com/adarsh-a-tw/tt-backend/db"
)

func (s *service) CreatePlayer(name string) error {
	return s.repo.CreatePlayer(&db.Player{Name: name})
}
