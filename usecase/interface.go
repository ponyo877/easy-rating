package usecase

import "github.com/ponyo877/easy-rating/domain"

type Repository interface {
	GetRate(playerID string) (int, error)
	SetRate(playerID string, rate int) error
	SetSortedRate(playerID string, rate int) error
	UpdateRate(playerID string, delta int) error
	GetRanking(offset int) ([]*domain.Player, error)
	StoreRate(playerID string, rate int) error
	GetResultByMatch(matchID string) (domain.Result, error)
	StoreResultByMatch(matchID string, result domain.Result) error
	GetPlayerByMatch(matchID string) (string, error)
	StorePlayerByMatch(matchID string, pID string) error
}

type Usecase interface {
	GetPlayersRate(p1ID, p2ID string) (int, int, error)
	IsExistRepost(matchID string) (bool, error)
	SaveReport(matchID, pID string, result domain.Result) error
	CheckReportWithPID(matchID, pID string, result domain.Result) (string, bool, error)
	DisableMatch(matchID string) error
	UpdateRate(p1ID, p2ID string, result domain.Result) error
}
