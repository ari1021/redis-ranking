package mysql

import (
	"context"

	"github.com/arkuchy/redis-ranking/src/db"
)

type UserResponse struct {
	ID        string
	Name      string
	HighScore int
	Rank      int
}

type UserMap struct {
	ID        string `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	HighScore int    `gorm:"not null;index:,sort:desc"`
}

// AddRanking は，usersテーブルにユーザデータを追加します
func AddRanking(ctx context.Context, id string, name string, score int) error {
	conn := db.Conn.GetMySQLConn()
	u := &UserMap{
		ID:        id,
		Name:      name,
		HighScore: score,
	}
	if err := conn.WithContext(ctx).Table("users").Create(u).Error; err != nil {
		return err
	}
	return nil
}

// GetRankings は，上位{limit}件のユーザデータを返します
func GetRankings(ctx context.Context, limit int) ([]*UserResponse, error) {
	conn := db.Conn.GetMySQLConn()
	us := make([]*UserMap, 0, limit)
	if err := conn.WithContext(ctx).Table("users").Order("high_score desc").Find(&us).Limit(limit).Error; err != nil {
		return nil, err
	}
	res := make([]*UserResponse, 0, limit)
	for i, u := range us {
		ur := &UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			HighScore: u.HighScore,
			Rank:      i,
		}
		res = append(res, ur)
	}
	return res, nil
}
