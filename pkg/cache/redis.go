package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

func Init(addr string) {
	rdb = redis.NewClient(&redis.Options{Addr: addr})
}

func Get(key string, value interface{}) error {
	data, err := rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func Set(key string, value interface{}, ttl int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, data, time.Duration(ttl)*time.Second).Err()
}
