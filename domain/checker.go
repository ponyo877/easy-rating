package domain

type Check int

const (
	CheckNotYet Check = iota
	CheckOK
)

func (c Check) Val() int {
	return int(c)
}
