package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Shelex/webhook-listener/entities"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Redis struct {
	client *redis.Client
}

func NewRedis() (Storage, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
		PoolSize: 100,
	})

	DB = &Redis{
		client: rdb,
	}

	return DB, nil
}

func (r *Redis) Add(hook entities.Hook) error {
	hook.Created_at = time.Now().UTC().Unix()
	return r.client.ZAddNX(ctx, hook.Channel, &redis.Z{
		Score:  float64(hook.Created_at),
		Member: hook,
	}).Err()
}

func (r *Redis) Get(channel string, pagination Pagination) ([]entities.Hook, int64, error) {
	minTime := strconv.Itoa(int(GetExpiryDate()))

	keys, err := r.client.ZRevRangeByScore(ctx, channel, &redis.ZRangeBy{
		Offset: pagination.Offset,
		Count:  pagination.Limit,
		Min:    minTime,
		Max:    "+inf",
	}).Result()

	if err != nil {
		return nil, 0, fmt.Errorf("failed to read items for channel %s, %s", channel, err)

	}

	hooks := make([]entities.Hook, len(keys))

	count, err := r.client.ZCount(ctx, channel, minTime, "+inf").Result()

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count for channel %s, %s", channel, err)
	}

	for i, key := range keys {
		var hook entities.Hook
		json.Unmarshal([]byte(key), &hook)
		hooks[i] = hook
	}

	return hooks, count, nil
}

func (r *Redis) Delete(channel string) error {
	return r.client.Unlink(ctx, channel).Err()
}

func (r *Redis) ClearExpired() error {
	log.Println("[redis-clear] clearing expired items")
	expiryDate := strconv.Itoa(int(GetExpiryDate()))
	iter := r.client.ScanType(ctx, 0, "*", 1000, "zset").Iterator()
	keysToRemove := make([]string, 0)

	for iter.Next(ctx) {
		channel := iter.Val()
		log.Printf("[redis-clear] checking channel: %s", channel)
		err := r.client.ZRemRangeByScore(ctx, channel, "-inf", expiryDate).Err()
		if err != nil {
			log.Printf("experienced error while clearing channel %s: %s", channel, err)
		}
		count, err := r.client.ZCount(ctx, channel, "-inf", "+inf").Result()
		if err != nil {
			log.Printf("experienced error while checking count for channel %s: %s", channel, err)
		}
		log.Printf("[redis-clear] channel %s has %d non-expired records", channel, count)
		if count == 0 {
			log.Printf("[redis-clear] channel %s will be removed", channel)
			keysToRemove = append(keysToRemove, channel)
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	// clear empty keys
	if len(keysToRemove) > 0 {
		log.Printf("[redis-clear] found %d empty channels to remove: %v", len(keysToRemove), keysToRemove)
		return r.client.Unlink(ctx, keysToRemove...).Err()
	}

	return nil
}
