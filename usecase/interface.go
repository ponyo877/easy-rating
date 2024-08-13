package usecase

import "github.com/ponyo877/easy-rating/domain"

type Repository interface {
	GetRate(playerID string) (int, error)
	StoreRate(playerID string, rate int) error
	GetResult(matchID string, playerID1 string, playerID2 string) (domain.Result, error)
	StoreResult(matchID string, playerID1 string, playerID2 string, result domain.Result) error
}

type Usecase interface {
	CheckReport(playerID1 string, playerID2 string, result domain.Result) error
	UpdateRate(playerID1 string, playerID2 string, result domain.Result) error
}
