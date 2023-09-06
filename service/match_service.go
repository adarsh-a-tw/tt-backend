package service

import (
	"fmt"

	"github.com/adarsh-a-tw/tt-backend/db"
	"github.com/adarsh-a-tw/tt-backend/enums"
)

type matchInfo struct {
	Id        int
	Format    enums.MatchFormat
	Stage     enums.MatchStage
	Status    enums.MatchStatus
	Opponents []struct {
		Name     string
		IsWinner bool
	}
}

func (s *service) GetMatchInfoList(status string) ([]matchInfo, error) {
	matches := []db.Match{}
	err := s.repo.GetAllMatches(&matches, status)
	if err != nil {
		return nil, err
	}

	matchInfoList := make([]matchInfo, 0)
	for _, match := range matches {
		opponents := make([]struct {
			Name     string
			IsWinner bool
		}, 2)
		if match.Format == string(enums.Doubles) {
			rows, err := s.repo.GetTeamInfoByMatchId(match.Id)
			if err != nil {
				return nil, err
			}
			for i, row := range rows {
				opponents[i] = struct {
					Name     string
					IsWinner bool
				}{fmt.Sprintf("%s & %s", row.PlayerA, row.PlayerB), row.IsWinner}
			}
		} else {
			rows, err := s.repo.GetPlayerInfoByMatchId(match.Id)
			if err != nil {
				return nil, err
			}
			for i, row := range rows {
				opponents[i] = struct {
					Name     string
					IsWinner bool
				}{row.PlayerName, row.IsWinner}
			}
		}
		matchInfoList = append(matchInfoList, matchInfo{
			Id:        match.Id,
			Format:    enums.MatchFormat(match.Format),
			Stage:     enums.MatchStage(match.Stage),
			Status:    enums.MatchStatus(match.Status),
			Opponents: opponents,
		})
	}

	return matchInfoList, nil
}

func (s *service) CreateSinglesMatch(
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

func (s *service) CreateDoublesMatch(
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

func (s *service) createMatch(
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
