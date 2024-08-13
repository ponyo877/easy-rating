package domain

import (
	"crypto/sha256"
	"fmt"
)

type Hash string

func NewHash(s string) Hash {
	return Hash(s)
}

func (h Hash) IsValid(pID string, solt string) bool {
	hash := sha256.Sum256([]byte(solt + pID + solt))
	return fmt.Sprintf("%x", hash) == string(h)
}
