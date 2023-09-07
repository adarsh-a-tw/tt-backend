package service

import (
	"errors"
	"math"

	"github.com/adarsh-a-tw/tt-backend/db"
	"github.com/adarsh-a-tw/tt-backend/enums"
)

var ErrGameOverOrSetCountExceeded = errors.New("game over or set count exceeded")
var ErrPreviousSetNotCompleted = errors.New("previous set not completed")
var ErrSetNotFound = errors.New("set not found")
var ErrSetAlreadyCompleted = errors.New("set already completed")

func (s *service) CreateSet(matchId int) error {

	match, err := s.repo.GetMatchById(matchId)
	if err != nil {
		return err
	}
	existing_sets, err := s.repo.GetSetsByMatchId(matchId)
	if err != nil {
		return err
	}
	if match.Status == string(enums.Past) || len(existing_sets) == match.SetCount {
		return ErrGameOverOrSetCountExceeded
	}

	for _, set := range existing_sets {
		if !set.IsCompleted {
			return ErrPreviousSetNotCompleted
		}
	}

	set := db.Set{
		SetNumber: len(existing_sets) + 1,
		MatchId:   matchId,
	}

	_, err = s.repo.CreateSet(&set)

	if match.Status == string(enums.Upcoming) {
		err = s.repo.UpdateMatchStatus(matchId, string(enums.Ongoing))
	}

	return err
}

func (s *service) UpdateScore(
	matchId int,
	setId int,
	scoredByA bool,
) error {
	match, err := s.repo.GetMatchById(matchId)
	if err != nil {
		return err
	}
	existing_sets, err := s.repo.GetSetsByMatchId(matchId)
	if err != nil {
		return err
	}
	if match.Status == string(enums.Past) {
		return ErrGameOverOrSetCountExceeded
	}
	var set *db.Set
	for _, s := range existing_sets {
		if s.Id == setId {
			set = &s
		}
	}
	if set == nil {
		return ErrSetNotFound
	}

	if set.IsCompleted {
		return ErrSetAlreadyCompleted
	}

	if scoredByA {
		set.OpponentAScore += 1
	} else {
		set.OpponentBScore += 1
	}

	setLog := &db.SetLog{
		SetId:     set.Id,
		OppAScore: set.OpponentAScore,
		OppBScore: set.OpponentBScore,
		ScoredByA: scoredByA,
	}

	err = s.repo.CreateSetLog(setLog)

	if err != nil {
		return err
	}

	set.IsCompleted = isSetComplete(*set, match.GamePoint)

	err = s.repo.UpdateSet(set)

	if err != nil {
		return err
	}

	if set.IsCompleted {
		s.handleMatchCompletion(match)
	}

	return err
}

func (s *service) handleMatchCompletion(match *db.Match) error {
	existing_sets, err := s.repo.GetSetsByMatchId(match.Id)
	if err != nil {
		return err
	}

	wonByA := 0
	wonByB := 0
	for _, set := range existing_sets {
		if set.OpponentAScore > set.OpponentBScore {
			wonByA += 1
		} else {
			wonByB += 1
		}
	}

	if hasMajorityWins(wonByA, match.SetCount) {
		err = s.repo.UpdateMatchWinner(match, true)
		if err != nil {
			return err
		}
		return s.repo.UpdateMatchStatus(match.Id, string(enums.Past))
	}

	if hasMajorityWins(wonByB, match.SetCount) {
		err = s.repo.UpdateMatchWinner(match, false)
		if err != nil {
			return err
		}
		return s.repo.UpdateMatchStatus(match.Id, string(enums.Past))
	}

	return nil
}

func hasMajorityWins(wins int, setCount int) bool {
	return (wins*100)/setCount > 50
}

func isSetComplete(set db.Set, gamePoint int) bool {
	scoreDiff := int(math.Abs(float64(set.OpponentAScore - set.OpponentBScore)))
	maxScore := int(math.Max(float64(set.OpponentAScore), float64(set.OpponentBScore)))

	return maxScore >= gamePoint && scoreDiff > 1
}
