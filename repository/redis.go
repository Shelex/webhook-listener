package repository

import (
	"context"
	"encoding/json"
	"fmt"
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
	return r.client.LPush(ctx, hook.Channel, hook).Err()
}

func (r *Redis) Get(channel string, pagination Pagination) ([]entities.Hook, int64, error) {
	endIndex := pagination.Offset + pagination.Limit - 1

	keys, err := r.client.LRange(ctx, channel, pagination.Offset, endIndex).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read items for channel %s, %s", channel, err)

	}

	hooks := make([]entities.Hook, len(keys))

	count, err := r.client.LLen(ctx, channel).Result()
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
