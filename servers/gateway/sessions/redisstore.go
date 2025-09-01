package sessions

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb *redis.Client
	ctx context.Context
	exp time.Duration
}

func NewRedisStore(client *redis.Client, expiration string) RedisStore {
	res := RedisStore{}
	res.rdb = client
	res.ctx = context.Background()
	res.exp, _ = time.ParseDuration(expiration)
	return res
}

func (rs *RedisStore) Get(key string) (string, error) {
	pipe := rs.rdb.Pipeline()

	getResult := pipe.Get(rs.ctx, key)
	pipe.Expire(rs.ctx, key, rs.exp).Err()

	_, err := pipe.Exec(rs.ctx)
	if err != nil && err.Error() != "redis: nil" {
		return "", err
	}

	val := getResult.Val()

	return val, nil
}

func (rs *RedisStore) Set(key string, value string) error {
	return rs.rdb.Set(rs.ctx, key, value, rs.exp).Err()
}

func (rs *RedisStore) Delete(key string) error {
	return rs.rdb.Del(rs.ctx, key).Err()
}
