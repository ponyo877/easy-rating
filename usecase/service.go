package usecase

import (
	"errors"

	"github.com/ponyo877/easy-rating/domain"
)

type RatingService struct {
	repo Repository
}

func NewRatingService(repo Repository) Usecase {
	return &RatingService{repo}
}

func (u *RatingService) FetchPlayerRate(pID string) (int, error) {
	rate, err := u.repo.GetRate(pID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return 0, err
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return rate, nil
	}
	initRate := 1500
	if err := u.repo.StoreRate(pID, initRate); err != nil {
		return 0, err
	}
	if err := u.repo.SetSortedRate(pID, initRate); err != nil {
		return 0, err
	}
	return initRate, nil
}

func (u *RatingService) GetRanking(offset int) ([]*domain.Player, error) {
	return u.repo.GetRanking(offset)
}

func (u *RatingService) IsExistRepost(matchID string) (bool, error) {
	_, err := u.repo.GetPlayerByMatch(matchID)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return false, err
	}
	return !errors.Is(err, domain.ErrNotFound), nil
}

func (u *RatingService) SaveReport(matchID, pID string, result domain.Result) error {
	if err := u.repo.StorePlayerByMatch(matchID, pID); err != nil {
		return err
	}
	if err := u.repo.StoreResultByMatch(matchID, result); err != nil {
		return err
	}
	return nil
}

func (u *RatingService) CheckReportWithPID(matchID, pID string, result domain.Result) (string, bool, error) {
	setPID, err := u.repo.GetPlayerByMatch(matchID)
	if err != nil {
		return "", false, err
	}
	if setPID == pID {
		return "", false, domain.ErrInvalidPlayer
	}
	setResult, err := u.repo.GetResultByMatch(matchID)
	if err != nil {
		return "", false, err
	}
	if !result.IsEquel(setResult) {
		return "", false, domain.ErrInvalidResult
	}
	return setPID, true, nil
}

func (u *RatingService) DisableMatch(matchID string) error {
	return u.repo.StoreResultByMatch(matchID, domain.ResultNotYet)
}

func (u *RatingService) UpdateRate(p1ID, p2ID string, result domain.Result) error {
	p1Rate, err := u.repo.GetRate(p1ID)
	if err != nil {
		return err
	}
	p1 := domain.NewPlayer(p1ID, p1Rate)
	p2Rate, err := u.repo.GetRate(p2ID)
	if err != nil {
		return err
	}
	p2 := domain.NewPlayer(p2ID, p2Rate)
	match := domain.NewMatch(p1, p2, result)
	newP1, newP2 := match.LatestPlayer()
	if err := u.repo.StoreRate(newP1.ID(), newP1.Rate()); err != nil {
		return err
	}
	if err := u.repo.StoreRate(newP2.ID(), newP2.Rate()); err != nil {
		return err
	}
	delta1, delta2 := match.Delta()
	if err := u.repo.UpdateRate(newP1.ID(), delta1); err != nil {
		return err
	}
	if err := u.repo.UpdateRate(newP2.ID(), delta2); err != nil {
		return err
	}
	return nil
}
