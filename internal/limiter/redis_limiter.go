// internal/limiter/redis_limiter.go
package limiter

import (
	"context"
	"fmt"

	//"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	rdb *redis.ClusterClient
}

func NewRedisClusterLimiter(rdb *redis.ClusterClient) *RedisLimiter {
	return &RedisLimiter{rdb: rdb}
}

func (r *RedisLimiter) SetRate(key string, rate float64, burst int) {
	ctx := context.Background()
	rateKey := fmt.Sprintf("rate:{%s}", key)
	r.rdb.HSet(ctx, rateKey, map[string]interface{}{
		"rate":  rate,
		"burst": burst,
	})
	r.rdb.Expire(ctx, rateKey, 24*time.Hour) // Set TTL for cleanup
}

var luaScript = redis.NewScript(`
local rate_key = KEYS[1]
local token_key = KEYS[2]

local rate_data = redis.call('HGETALL', rate_key)
if not rate_data[1] then
  return 1  -- allow by default if no rate set
end

local rate = tonumber(rate_data[2])
local burst = tonumber(rate_data[4])
local now = tonumber(ARGV[1])
local ttl = tonumber(ARGV[2])

local last_time = tonumber(redis.call('HGET', token_key, 'ts') or 0)
local tokens = tonumber(redis.call('HGET', token_key, 'tokens') or burst)

local elapsed = now - last_time
local refill = elapsed * rate
local new_tokens = math.min(tokens + refill, burst)

if new_tokens < 1 then
  redis.call('HSET', token_key, 'tokens', new_tokens, 'ts', now)
  redis.call('EXPIRE', token_key, ttl)
  return 0  -- rate limited
else
  redis.call('HSET', token_key, 'tokens', new_tokens - 1, 'ts', now)
  redis.call('EXPIRE', token_key, ttl)
  return 1  -- allowed
end
`)

func (r *RedisLimiter) Allow(key string) bool {
	ctx := context.Background()
	rateKey := fmt.Sprintf("rate:{%s}", key)
	tokenKey := fmt.Sprintf("bucket:{%s}", key)
	now := time.Now().Unix()
	ttlSeconds := int64(24 * 60 * 60)

	result, err := luaScript.Run(ctx, r.rdb, []string{rateKey, tokenKey}, now, ttlSeconds).Int()
	if err != nil {
		fmt.Printf("Redis error: %v\n", err)
		return true // fail open
	}
	return result == 1
}
