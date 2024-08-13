package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ponyo877/easy-rating/repository"
	"github.com/ponyo877/easy-rating/usecase"
	"github.com/redis/go-redis/v9"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type Request struct {
	MatchID  string `json: "match_id"`
	PlayerID string `json: "player_id"`
	Result   int    `json: "result"`
}

func main() {
	flag.Parse()
	redisURL := os.Getenv("REDIS_URL")
	redisToken := os.Getenv("REDIS_TOKEN")

	opt, _ := redis.ParseURL(fmt.Sprintf("rediss://default:%s@%s:6379", redisToken, redisURL))
	repository := repository.NewRatingRepository(redis.NewClient(opt))
	service := usecase.NewRatingService(repository)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		json.NewDecoder(r.Body).Decode(&req)
		service.UpdateRate("1", "2", 1)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
