package dao

import (
	"github.com/go-redis/redis/v8"
	"tmp/conf"
)

var Rdb *redis.Client

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: conf.RedisAddr,
		DB:   0, // use default DB
	})
}
