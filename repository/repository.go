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

func (r *RatingRepository) SetRate(playerID string, rate int) error {
	return r.redis.Set(context.Background(), playerID, rate, 0).Err()
}

func (r *RatingRepository) SetSortedRate(playerID string, rate int) error {
	return r.redis.ZAdd(context.Background(), "player-ranking", redis.Z{
		Score:  float64(rate),
		Member: playerID,
	}).Err()
}

func (r *RatingRepository) UpdateRate(playerID string, delta int) error {
	_, err := r.redis.ZIncrBy(context.Background(), "player-ranking", float64(delta), playerID).Result()
	return err
}

func (r *RatingRepository) GetRanking(offset int) ([]*domain.Player, error) {
	vals, err := r.redis.ZRevRangeWithScores(context.Background(), "player-ranking", 0, int64(offset)).Result()
	if err != nil {
		return nil, err
	}
	players := make([]*domain.Player, len(vals))
	for i, v := range vals {
		players[i] = domain.NewPlayer(v.Member.(string), int(v.Score))
	}
	return players, nil
}

func (r *RatingRepository) StoreRate(playerID string, rate int) error {
	return r.redis.Set(context.Background(), playerID, rate, 0).Err()
}

func (r *RatingRepository) GetResultByMatch(matchID string) (domain.Result, error) {
	val, err := r.redis.Get(context.Background(), "result:"+matchID).Result()
	return domain.NewFromEloScore(val), err
}

func (r *RatingRepository) StoreResultByMatch(matchID string, result domain.Result) error {
	return r.redis.Set(context.Background(), "result:"+matchID, result.EloScore(), 0).Err()
}

func (r *RatingRepository) GetPlayerByMatch(matchID string) (string, error) {
	return r.redis.Get(context.Background(), "player:"+matchID).Result()
}

func (r *RatingRepository) StorePlayerByMatch(matchID string, pID string) error {
	return r.redis.Set(context.Background(), "player:"+matchID, pID, 0).Err()
}
