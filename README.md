# easy-rating
Simple rating management server implemented in Go for P2P games.

## Demo

![demo is WIP]()

The following is being used as the backend service. Please try it out.   
https://hit-and-blow.pages.dev/go

## Usage
Server start (default listen port: 8001, change with -flag option)
```bash
go run ./main.go
```

## API Endpoints
- Start Rating: Initialize rating for players
```bash
curl "http://localhost:8001/rating/start?p1=test-user1&p2=test-user2"
```

- Finish Rating: Submit match
```bash
curl -X POST -H "Content-Type: application/json" -d '{"match_id":"match1","player_id":"test-user1","number":1,"hash":"abc123","result":"win"}' "http://localhost:8001/rating/finish"
```

- Get Ranking: Retrieve player rankings
```bash
curl "http://localhost:8001/rating/ranking"
```

## Environment Variables
| KEY | Description |
| --- | --- |
|REDIS_URL | URL for the Redis instance|
|REDIS_TOKEN | Token for Redis authentication|
|SOLT | Secret key for hashing|
```	