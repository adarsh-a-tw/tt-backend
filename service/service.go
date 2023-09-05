package service

import "github.com/adarsh-a-tw/tt-backend/db"

type Service struct {
	repo db.Repository
}

func NewService(repo db.Repository) *Service {
	return &Service{repo: repo}
}
