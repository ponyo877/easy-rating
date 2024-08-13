package usecase

import "github.com/ponyo877/easy-rating/domain"

type RatingService struct {
	repo Repository
}

func NewRatingService(repo Repository) Usecase {
	return &RatingService{repo}
}

func (u *RatingService) CheckReport(playerID1 string, playerID2 string, result domain.Result) error {
	// TODO not implemented
	return nil
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
	newP1, newP2 := match.CalcRate()
	if err := u.repo.StoreRate(newP1.ID(), newP1.Rate()); err != nil {
		return err
	}
	if err := u.repo.StoreRate(newP2.ID(), newP2.Rate()); err != nil {
		return err
	}
	return nil
}
