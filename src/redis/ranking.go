package redis

import (
	"context"
	"encoding/json"

	"github.com/ari1021/redis-ranking/src/db"
	goRedis "github.com/go-redis/redis/v8"
)

const (
	RedisRanking string = "RedisRanking"
)

type UserResponse struct {
	ID        string
	Name      string
	HighScore int
	Rank      int
}

// userIDとuserNameを持った構造体(json文字列にして扱う)
type Member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// AddRanking は，ランキングにユーザデータを追加します
func AddRanking(ctx context.Context, id string, name string, score int) error {
	conn := db.Conn.GetRedisConn()
	member := &Member{
		ID:   id,
		Name: name,
	}
	// memberをserializeする
	serializedMember, err := json.Marshal(member)
	if err != nil {
		return err
	}
	if err := conn.ZAdd(ctx, RedisRanking, &goRedis.Z{
		Score:  float64(score),
		Member: serializedMember,
	}).Err(); err != nil {
		return err
	}
	return nil
}

// GetRankings は，上位{limit}件のユーザデータを返します
func GetRankings(ctx context.Context, limit int) ([]*UserResponse, error) {
	// redisは0始まり
	// ex) 1~10 -> start:0, stop:9
	start := 0
	stop := start + limit - 1
	conn := db.Conn.GetRedisConn()
	serializedMembersWithScores, err := conn.ZRevRangeWithScores(ctx, RedisRanking, int64(start), int64(stop)).Result()
	if err != nil {
		return nil, err
	}
	res := make([]*UserResponse, 0, limit)
	member := &Member{}
	for i, serializedMemberWithScore := range serializedMembersWithScores {
		serializedMember := serializedMemberWithScore.Member
		score := serializedMemberWithScore.Score
		if err := json.Unmarshal([]byte(serializedMember.(string)), member); err != nil {
			return nil, err
		}
		u := &UserResponse{
			ID:        member.ID,
			Name:      member.Name,
			HighScore: int(score),
			Rank:      i + 1,
		}
		res = append(res, u)
	}
	return res, nil
}
