package main

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func TestLeaderboard(t *testing.T) {
	// 启动 miniredis 服务器
	mrs, err := miniredis.Run()
	assert.NoError(t, err)
	defer mrs.Close()

	// 创建 Redis 客户端连接到 miniredis
	client := redis.NewClient(&redis.Options{
		Addr: mrs.Addr(),
	})

	// 创建排行榜实例
	leaderboard := NewLeaderboard(client, "test_leaderboard")

	// 添加用户和分数
	err = leaderboard.AddUser("Alice", 100)
	assert.NoError(t, err)
	err = leaderboard.AddUser("Bob", 200)
	assert.NoError(t, err)
	err = leaderboard.AddUser("Charlie", 150)
	assert.NoError(t, err)

	// 获取用户排名
	rank, err := leaderboard.GetUserRank("Bob")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rank)

	rank, err = leaderboard.GetUserRank("Charlie")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), rank)

	rank, err = leaderboard.GetUserRank("Alice")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), rank)

	// 获取前 3 名用户
	topUsers, err := leaderboard.GetTopUsers(3)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(topUsers))
	assert.Equal(t, "Bob", topUsers[0].Member)
	assert.Equal(t, "Charlie", topUsers[1].Member)
	assert.Equal(t, "Alice", topUsers[2].Member)
}
