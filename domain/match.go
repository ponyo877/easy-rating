package domain

import (
	elogo "github.com/kortemy/elo-go"
)

var elo *elogo.Elo

func init() {
	elo = elogo.NewElo()
}

type Match struct {
	p1     *Player
	p2     *Player
	result Result
}

func NewMatch(p1, p2 *Player, r Result) *Match {
	return &Match{p1, p2, r}
}

func (m *Match) CalcRate() (*Player, *Player) {
	o1, o2 := elo.Outcome(m.p1.rate, m.p2.rate, m.result.EloScore())
	return NewPlayer(m.p1.id, o1.Rating), NewPlayer(m.p2.id, o2.Rating)
}
