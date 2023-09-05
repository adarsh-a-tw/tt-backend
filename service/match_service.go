package service

import (
	"github.com/adarsh-a-tw/tt-backend/db"
	"github.com/adarsh-a-tw/tt-backend/enums"
)

func (s *Service) CreateSinglesMatch(
	stage enums.MatchStage,
	playerAId int,
	playerBId int,
	maxSets int,
	gamePoint int,
) error {
	id, err := s.createMatch(enums.Singles, stage, maxSets, gamePoint)
	if err != nil {
		return err
	}

	err = s.repo.AddPlayerToMatch(
		&db.PlayerMatchMapping{MatchId: int(id), PlayerId: playerAId, IsOpponentA: true},
	)
	if err != nil {
		return err
	}

	err = s.repo.AddPlayerToMatch(
		&db.PlayerMatchMapping{MatchId: int(id), PlayerId: playerBId, IsOpponentA: false},
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *Service) CreateDoublesMatch(
	stage enums.MatchStage,
	teamAId int,
	teamBId int,
	maxSets int,
	gamePoint int,
) error {
	id, err := s.createMatch(enums.Doubles, stage, maxSets, gamePoint)
	if err != nil {
		return err
	}

	err = s.repo.AddTeamToMatch(
		&db.TeamMatchMapping{MatchId: int(id), TeamId: teamAId, IsOpponentA: true},
	)
	if err != nil {
		return err
	}

	err = s.repo.AddTeamToMatch(
		&db.TeamMatchMapping{MatchId: int(id), TeamId: teamBId, IsOpponentA: false},
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *Service) createMatch(
	format enums.MatchFormat,
	stage enums.MatchStage,
	maxSets int,
	gamePoint int,
) (int64, error) {
	match := &db.Match{
		Format:    string(format),
		Stage:     string(stage),
		SetCount:  maxSets,
		GamePoint: gamePoint,
		Status:    string(enums.Upcoming),
	}
	return s.repo.CreateMatch(match)
}
