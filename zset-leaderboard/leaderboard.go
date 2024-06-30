package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Leaderboard struct {
	client *redis.Client
	key    string
}

func NewLeaderboard(client *redis.Client, key string) *Leaderboard {
	return &Leaderboard{
		client: client,
		key:    key,
	}
}

func (lb *Leaderboard) AddUser(user string, score float64) error {
	return lb.client.ZAdd(ctx, lb.key, &redis.Z{
		Score:  score,
		Member: user,
	}).Err()
}

func (lb *Leaderboard) GetUserRank(user string) (int64, error) {
	rank, err := lb.client.ZRevRank(ctx, lb.key, user).Result()
	if err != nil {
		return 0, err
	}
	return rank + 1, nil
}

func (lb *Leaderboard) GetTopUsers(limit int64) ([]redis.Z, error) {
	return lb.client.ZRevRangeWithScores(ctx, lb.key, 0, limit-1).Result()
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	leaderboard := NewLeaderboard(client, "game_leaderboard")
	leaderboard.AddUser("Alice", 100)
	leaderboard.AddUser("Bob", 200)
	leaderboard.AddUser("Charlie", 150)

	var name rune
	for i := 1; i <= 10; i++ {
		name = rune('A' + i)
		score := i * 200
		leaderboard.AddUser(string(name), float64(score))
	}

	topUsers, _ := leaderboard.GetTopUsers(3)
	for _, user := range topUsers {
		fmt.Printf("User: %s, Score: %f\n", user.Member, user.Score)
	}
}
