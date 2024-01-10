package sessions

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb *redis.Client
	ctx context.Context
}

func (rs *RedisStore) New(addr string) RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",
		DB: 0,
	})
	rs.rdb = rdb
	rs.ctx = context.Background()
}

func (rs *RedisStore) Get(sessionID string) (int, error) {
	val, err := rs.rdb.Get(rs.ctx, sessionID).Result()
	if err != nil {
		return 0, err
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (rs *RedisStore) Set(sessionID string, newUserID int) error {
	return rs.rdb.Set(rs.ctx, sessionID, newUserID, 1800).Err()
} 

func (rs *RedisStore) Delete(sessionID string) error {
	return rs.rdb.Del(rs.ctx, sessionID).Err()
}