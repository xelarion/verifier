package verifier

import (
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
