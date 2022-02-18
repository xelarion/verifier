package verifier

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// TokenStorage token 存储器
type TokenStorage interface {
	Get(key string) (string, error) // return error when key does not exist
	Set(key string, value interface{}, expiration time.Duration) error
	SetNX(key string, value interface{}, expiration time.Duration) bool
	Del(key string) error
	DelByKeyPrefix(keyPrefix string) error
	Exists(key string) bool
}

type RedisStorage struct {
	Client *redis.Client
}

// a non-nil, empty Context
var backgroundContext = context.Background()

func (r *RedisStorage) Get(key string) (string, error) {
	return r.Client.Get(backgroundContext, key).Result()
}

func (r *RedisStorage) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(backgroundContext, key, value, expiration).Err()
}

func (r *RedisStorage) SetNX(key string, value interface{}, expiration time.Duration) bool {
	return r.Client.SetNX(backgroundContext, key, value, expiration).Val()
}

func (r *RedisStorage) Del(key string) error {
	return r.Client.Del(backgroundContext, key).Err()
}

func (r *RedisStorage) DelByKeyPrefix(keyPrefix string) error {
	var cursor uint64
	var keys []string
	var err error

	for {
		keys, cursor, err = r.Client.Scan(backgroundContext, cursor, keyPrefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err = r.Client.Del(backgroundContext, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *RedisStorage) Exists(key string) bool {
	return r.Client.Exists(backgroundContext, key).Val() == 1
}
