package service

import (
	"fmt"

	"github.com/adarsh-a-tw/tt-backend/db"
	"github.com/adarsh-a-tw/tt-backend/enums"
)

type opponent struct {
	Id       int
	Name     string
	IsWinner bool
}

type matchInfo struct {
	Id        int
	Format    enums.MatchFormat
	Stage     enums.MatchStage
	Status    enums.MatchStatus
	Opponents []opponent
}

func (s *service) GetMatchInfoList(status string) ([]matchInfo, error) {
	matches := []db.Match{}
	err := s.repo.GetAllMatches(&matches, status)
	if err != nil {
		return nil, err
	}

	matchInfoList := make([]matchInfo, 0)
	for _, match := range matches {
		opponents, err := s.opponentsFromMatch(match)
		if err != nil {
			return nil, err
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

func (s *service) opponentsFromMatch(match db.Match) ([]opponent, error) {
	opponents := make([]opponent, 2)
	if match.Format == string(enums.Doubles) {
		rows, err := s.repo.GetTeamInfoByMatchId(match.Id)
		if err != nil {
			return nil, err
		}
		for i, row := range rows {
			opponents[i] = opponent{row.TeamId, fmt.Sprintf("%s & %s", row.PlayerA, row.PlayerB), row.IsWinner}
		}
	} else {
		rows, err := s.repo.GetPlayerInfoByMatchId(match.Id)
		if err != nil {
			return nil, err
		}
		for i, row := range rows {
			opponents[i] = opponent{row.PlayerId, row.PlayerName, row.IsWinner}
		}
	}
	return opponents, nil
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

type set struct {
	Id             int
	SetNumber      int
	OpponentAScore int
	OpponentBScore int
	IsCompleted    bool
	Logs           []setLog
}

type setLog struct {
	Id        int
	OppAScore int
	OppBScore int
	ScoredByA bool
}

type MatchDetail struct {
	Id        int
	Format    string
	Stage     string
	Status    string
	Opponents []opponent
	Sets      []set
}

func (svc *service) GetMatchDetails(matchId int) (*MatchDetail, error) {
	match, err := svc.repo.GetMatchById(matchId)
	if err != nil {
		return nil, err
	}

	opponents, err := svc.opponentsFromMatch(*match)
	if err != nil {
		return nil, err
	}

	setsFromDb, err := svc.repo.GetSetsByMatchId(matchId)
	if err != nil {
		return nil, err
	}

	sets := make([]set, 0)
	for _, s := range setsFromDb {
		setLogsFromDb, err := svc.repo.GetSetLogsBySetId(s.Id, nil)
		if err != nil {
			return nil, err
		}

		setLogs := make([]setLog, 0)
		for _, sl := range setLogsFromDb {
			setLogs = append(setLogs, setLog{
				Id:        sl.Id,
				OppAScore: sl.OppAScore,
				OppBScore: sl.OppBScore,
				ScoredByA: sl.ScoredByA,
			})
		}

		sets = append(sets, set{
			Id:             s.Id,
			SetNumber:      s.SetNumber,
			OpponentAScore: s.OpponentAScore,
			OpponentBScore: s.OpponentBScore,
			IsCompleted:    s.IsCompleted,
			Logs:           setLogs,
		})
	}

	return &MatchDetail{
		Id:        matchId,
		Format:    match.Format,
		Stage:     match.Stage,
		Status:    match.Status,
		Opponents: opponents,
		Sets:      sets,
	}, nil
}
