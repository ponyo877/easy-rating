package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ponyo877/easy-rating/domain"
	"github.com/ponyo877/easy-rating/repository"
	"github.com/ponyo877/easy-rating/usecase"
	"github.com/redis/go-redis/v9"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type Request struct {
	MatchID  string `json:"match_id"`
	PlayerID string `json:"player_id"`
	Number   int    `json:"number"`
	Hash     string `json:"hash"`
	Result   string `json:"result"`
}

type Player struct {
	ID   string `json:"id"`
	Rate int    `json:"rate"`
}

type Response struct {
	Player1 Player `json:"player1"`
	Player2 Player `json:"player2"`
}

func main() {
	flag.Parse()
	redisURL := os.Getenv("REDIS_URL")
	redisToken := os.Getenv("REDIS_TOKEN")
	solt := os.Getenv("SOLT")

	opt, _ := redis.ParseURL(fmt.Sprintf("rediss://default:%s@%s:6379", redisToken, redisURL))
	repository := repository.NewRatingRepository(redis.NewClient(opt))
	service := usecase.NewRatingService(repository)
	http.HandleFunc("/start", startHandler(service, solt))
	http.HandleFunc("/finish", finishHandler(service, solt))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func startHandler(service usecase.Usecase, solt string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		v := r.URL.Query()
		p1ID, p2ID := v.Get("p1"), v.Get("p2")
		p1Rate, p2Rate, err := service.GetPlayersRate(p1ID, p2ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res := Response{
			Player1: Player{
				ID:   p1ID,
				Rate: p1Rate,
			},
			Player2: Player{
				ID:   p2ID,
				Rate: p2Rate,
			},
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
		}
	}
}

func finishHandler(service usecase.Usecase, solt string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !domain.NewHash(req.Hash).IsValid(req.PlayerID, solt) {
			http.Error(w, "Invalid Hash", http.StatusUnauthorized)
			return
		}
		result := domain.NewFromEloScore(req.Result)
		isExist, err := service.IsExistRepost(req.MatchID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !isExist {
			if err := service.SaveReport(req.MatchID, req.PlayerID, result); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		pID, ok, err := service.CheckReportWithPID(req.MatchID, req.PlayerID, result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "Invalid Request", http.StatusBadRequest)
			return
		}
		if err := service.DisableMatch(req.MatchID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var p1ID, p2ID string
		if req.Number == 1 {
			p1ID, p2ID = pID, req.PlayerID
		} else {
			p1ID, p2ID = req.PlayerID, pID
		}
		if err := service.UpdateRate(p1ID, p2ID, result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
