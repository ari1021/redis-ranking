package db

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connection struct {
	redisConn *redis.Client
	mysqlConn *gorm.DB
}

// GetRedisConn は，構造体Connectionが保持している *redis.Clientを返します
func (c *Connection) GetRedisConn() *redis.Client {
	return c.redisConn
}

// GetMySQLConn は，構造体Connectionが保持している *gorm.DBを返します
func (c *Connection) GetMySQLConn() *gorm.DB {
	return c.mysqlConn
}

var Conn = &Connection{}

// InitializeConnection は，Connectionを確立して，Connに保存します
func InitializeConnection() {
	redisConn := NewRedisConnection()
	mysqlConn, _ := NewMySQLConnection()
	Conn.redisConn = redisConn
	Conn.mysqlConn = mysqlConn
}

// NewRedisConnection は，*redis.Clientを生成します
func NewRedisConnection() *redis.Client {
	addr := fmt.Sprintf(
		"%s:%s",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"),
	)
	conn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return conn
}

// NewMySQLConnection は，migrationを行い，*gorm.DBを生成します
func NewMySQLConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	) + "?parseTime=true&collation=utf8mb4_bin"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL: %w", err)
	}
	return conn, nil
}
