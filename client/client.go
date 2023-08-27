package client

import (
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.17.128:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	Rdb = rdb
}
