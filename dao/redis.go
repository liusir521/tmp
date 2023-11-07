package dao

import "github.com/go-redis/redis/v8"

var Rdb *redis.Client

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "10.67.68.161:6379",
		DB:   0, // use default DB
	})
}
