package domain

type Result int

const (
	ResultNotYet Result = iota
	ResultOneWin
	ResultTwoWin
	ResultDraw
)

func (r Result) EloScore() float64 {
	switch r {
	case ResultOneWin:
		return 1
	case ResultTwoWin:
		return 0
	case ResultDraw:
		return 0.5
	default:
		return -1
	}
}

func NewFromEloScore(score string) Result {
	switch score {
	case "1":
		return ResultOneWin
	case "0":
		return ResultTwoWin
	case "0.5":
		return ResultDraw
	default:
		return ResultNotYet
	}
}

func (r1 Result) IsEquel(r2 Result) bool {
	return r1 == r2
}
