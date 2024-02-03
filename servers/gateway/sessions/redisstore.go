package sessions

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb *redis.Client
	ctx context.Context
	exp string
}

func NewRedisStore(addr string, expiration string) RedisStore {
	res := RedisStore{}
	res.rdb = redis.NewClient(&redis.Options{
		Addr: addr,
		Password: "",
		DB: 0,
	})
	res.ctx = context.Background()
	res.exp = expiration
	return res
}

func (rs *RedisStore) Get(sessionID string) (int, error) {
	val, err := rs.rdb.Get(rs.ctx, sessionID).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, nil
		}
		return 0, err
	}
	duration, _ := time.ParseDuration(rs.exp)
	err = rs.rdb.Expire(rs.ctx, sessionID, duration).Err()
	if err != nil {
		log.Print("error resetting Redis expiration timer")
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (rs *RedisStore) Set(sessionID string, userID int) error {
	duration, _ := time.ParseDuration(rs.exp)
	return rs.rdb.Set(rs.ctx, sessionID, userID, duration).Err()
} 

func (rs *RedisStore) Delete(sessionID string) error {
	return rs.rdb.Del(rs.ctx, sessionID).Err()
}