package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"

	goRedis "github.com/go-redis/redis/v8"

	"github.com/ari1021/redis-ranking/src/db"
	"github.com/ari1021/redis-ranking/src/mysql"
	"github.com/ari1021/redis-ranking/src/redis"
)

type User struct {
	ID        string
	Name      string
	HighScore int
}

// init/main は，コマンドライン引数から初期データ数を受け取り，初期データの挿入を行います
func main() {
	var dataNum int
	flag.IntVar(&dataNum, "dataNum", 1000, "the number of data")
	flag.Parse()
	db.InitializeConnection()
	data := make([]*User, 0, dataNum)
	for i := 0; i < dataNum; i++ {
		id := fmt.Sprintf("ID-%d", i)
		name := fmt.Sprintf("Name-%d", i)
		score := rand.Intn(dataNum)
		u := &User{
			ID:        id,
			Name:      name,
			HighScore: score,
		}
		data = append(data, u)
	}
	if err := InsertMySQLData(data); err != nil {
		log.Fatal(err)
	}
	log.Println("finished inserting data to MySQL")
	if err := InsertRedisData(data); err != nil {
		log.Fatal(err)
	}
	log.Println("finished inserting data to Redis")
}

// InsertMySQLData は，dataを受け取り，batch insertを行います
func InsertMySQLData(data []*User) error {
	users := make([]*mysql.UserMap, 0, len(data))
	for _, user := range data {
		userMap := &mysql.UserMap{
			ID:        user.ID,
			Name:      user.Name,
			HighScore: user.HighScore,
		}
		users = append(users, userMap)
	}
	conn := db.Conn.GetMySQLConn()
	if err := conn.Table("users").Create(&users).Error; err != nil {
		return err
	}
	return nil
}

// InsertRedisData は，dataを受け取り，pipelineを用いてinsertを行います
func InsertRedisData(data []*User) error {
	conn := db.Conn.GetRedisConn()
	pipe := conn.Pipeline()
	for i, user := range data {
		member := &redis.Member{
			ID:   user.ID,
			Name: user.Name,
		}
		serializedMember, err := json.Marshal(member)
		if err != nil {
			return err
		}
		if err := pipe.ZAdd(context.TODO(), redis.RedisRanking, &goRedis.Z{
			Score:  float64(user.HighScore),
			Member: serializedMember,
		}).Err(); err != nil {
			return err
		}
		if i%1000 == 0 {
			if _, err := pipe.Exec(context.TODO()); err != nil {
				return err
			}
		}
	}
	if _, err := pipe.Exec(context.TODO()); err != nil {
		return err
	}
	return nil
}
