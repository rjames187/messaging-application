package sessions

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb *redis.Client
	ctx context.Context
	exp time.Duration
}

func NewRedisStore(addr string, expiration string) RedisStore {
	res := RedisStore{}
	res.rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	res.ctx = context.Background()
	res.exp, _ = time.ParseDuration(expiration)
	return res
}

func (rs *RedisStore) Get(sessionID string) (int, error) {
	pipe := rs.rdb.Pipeline()

	getResult := pipe.Get(rs.ctx, sessionID)
	pipe.Expire(rs.ctx, sessionID, rs.exp).Err()

	_, err := pipe.Exec(rs.ctx)
	if err != nil && err.Error() != "redis: nil" {
		return 0, err
	}

	val := getResult.Val()
	if val == "" {
		return 0, nil
	}

	res, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (rs *RedisStore) Set(sessionID string, userID int) error {
	return rs.rdb.Set(rs.ctx, sessionID, userID, rs.exp).Err()
}

func (rs *RedisStore) Delete(sessionID string) error {
	return rs.rdb.Del(rs.ctx, sessionID).Err()
}
