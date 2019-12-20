package conn

import (
	"github.com/go-redis/redis/v7"
)

type RedisConfig struct {
	host     string
	port     string
	password string
}

var (
	RedisClient *redis.Client
)

func Redis() *redis.Client {
	return RedisClient
}

func NewRedisConfig(host, port, password string) *RedisConfig {
	return &RedisConfig{host: host, port: port, password: password}
}

func (r *RedisConfig) Conn() error {
	if RedisClient == nil {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}
	return nil
}

func (r *RedisConfig) Close() error {
	return RedisClient.Close()
}

