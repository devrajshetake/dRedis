package dredis

import "time"

type CacheItem struct {
	CacheType string
	ExpiresAt time.Duration
	Value     string
}
