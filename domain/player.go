package domain

type Player struct {
	id   string
	rate int
}

func NewPlayer(id string, rate int) *Player {
	return &Player{
		id:   id,
		rate: rate,
	}
}

func (p *Player) ID() string {
	return p.id
}

func (p *Player) Rate() int {
	return p.rate
}
