package models

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RateLimiter struct {
	Client   *redis.Client
	Rate     int
	Interval time.Duration
}

func NewRateLimiter(client *redis.Client, rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		Client:   client,
		Rate:     rate,
		Interval: interval,
	}
}

func (rl *RateLimiter) Allow(key string) (bool, error) {
	now := time.Now().UnixNano() / int64(time.Millisecond)

	script := `
	local key = KEYS[1]
	local rate = tonumber(ARGV[1])
	local interval = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local expireTime = interval / 1000

	redis.call('ZREMRANGEBYSCORE', key, 0, now - interval)
	local currentCount = redis.call('ZCARD', key)

	if currentCount < rate then
		redis.call('ZADD', key, now, now)
		redis.call('EXPIRE', key, expireTime)
		return 1
	else
		return 0
	end
	`

	result, err := rl.Client.Eval(ctx, script, []string{key}, rl.Rate, rl.Interval.Milliseconds(), now).Result()
	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}
