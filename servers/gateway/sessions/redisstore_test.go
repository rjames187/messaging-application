//go:build !no_db

package sessions

import (
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func getClient() *redis.Client {
	addr := os.Getenv("REDISADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return client
}

var redisClient = getClient()

func TestSetGet(t *testing.T) {
	client := NewRedisStore(redisClient, "3s")
	err := client.Set("key", "3")
	if err != nil {
		t.Error("error setting key to 3")
	}
	err = client.Set("key2", "5")
	if err != nil {
		t.Error("error setting key2 to 5")
	}
	val, err := client.Get("key")
	if err != nil {
		t.Error("error fetching key")
	}
	if val != "3" {
		t.Errorf("fetched key should have equaled 3 and not %s", val)
	}
}

func TestSetOverride(t *testing.T) {
	client := NewRedisStore(redisClient, "3s")
	err := client.Set("key", "3")
	if err != nil {
		t.Error("error setting key to 3")
	}
	err = client.Set("key", "5")
	if err != nil {
		t.Error("error setting key to 5")
	}
	val, err := client.Get("key")
	if err != nil {
		t.Error("error fetching key")
	}
	if val != "5" {
		t.Errorf("fetched key should have equaled 5 and not %s", val)
	}
}

func TestExpiration(t *testing.T) {
	client := NewRedisStore(redisClient, "2s")
	err := client.Set("k", "3")
	if err != nil {
		t.Error("error setting k to 3")
	}
	duration, _ := time.ParseDuration("3s")
	time.Sleep(duration)
	_, err = client.Get("k")
	if err != nil {
		t.Errorf("error fetching key k: %s", err)
	}
}

func TestReset(t *testing.T) {
	client := NewRedisStore(redisClient, "10s")
	err := client.Set("i", "3")
	if err != nil {
		t.Error("error setting i to 3")
	}
	duration, _ := time.ParseDuration("5s")
	time.Sleep(duration)
	client.Get("i")
	time.Sleep(duration)
	val, err := client.Get("i")
	if err != nil {
		t.Errorf("error fetching key i: %s", err)
	}
	if val == "0" {
		t.Error("key k should have been present")
	}
	if val != "3" {
		t.Error("key k should have equaled 3")
	}
}
