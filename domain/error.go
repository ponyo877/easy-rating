package domain

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidPlayer error = errors.New("invalid player")
	ErrInvalidResult error = errors.New("invalid result")
	ErrNotFound      error = redis.Nil
)
