package repository

import (
	"context"
	"strconv"

	"github.com/ponyo877/easy-rating/domain"
	"github.com/ponyo877/easy-rating/usecase"
	"github.com/redis/go-redis/v9"
)

type RatingRepository struct {
	redis *redis.Client
}

func NewRatingRepository(redis *redis.Client) usecase.Repository {
	return &RatingRepository{redis}
}

func (r *RatingRepository) GetRate(playerID string) (int, error) {
	val, err := r.redis.Get(context.Background(), playerID).Result()
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(val)
}

func (r *RatingRepository) StoreRate(user string, rate int) error {
	return r.redis.Set(context.Background(), user, rate, 0).Err()
}

func (r *RatingRepository) StoreResultCheck(matchID string, playerID string, check domain.Check) error {
	return r.redis.Set(context.Background(), matchID+":"+playerID, check.Val(), 0).Err()
}

func (r *RatingRepository) StoreResult(matchID string, playerID1, playerID2 string, result domain.Result) error {
	return r.redis.Set(context.Background(), matchID+":"+playerID1+"vs"+playerID2, result.EloScore(), 0).Err()
}

func (r *RatingRepository) GetResult(matchID string, playerID1, playerID2 string) (domain.Result, error) {
	val, err := r.redis.Get(context.Background(), matchID+":"+playerID1+"vs"+playerID2).Result()
	return domain.NewFromEloScore(val), err
}
