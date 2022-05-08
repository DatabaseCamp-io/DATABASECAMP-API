package cache

import "time"

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
}
