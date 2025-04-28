package limiter

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	client        *redis.Client
	defaultLimit  int
	blockDuration time.Duration
	tokenLimits   map[string]int
	ctx           context.Context
}

func NewRedisRateLimiter() *RedisRateLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // Adapte se necessário
		DB:       0,
	})

	defaultLimit, _ := strconv.Atoi(os.Getenv("DEFAULT_LIMIT"))
	blockDuration, _ := strconv.Atoi(os.Getenv("BLOCK_DURATION"))

	// Carrega limites específicos de token do .env
	tokenLimits := make(map[string]int)
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "TOKEN_LIMIT_") {
			parts := strings.SplitN(env, "=", 2)
			key := strings.TrimPrefix(parts[0], "TOKEN_LIMIT_")
			val, _ := strconv.Atoi(parts[1])
			tokenLimits[key] = val
		}
	}

	return &RedisRateLimiter{
		client:        rdb,
		defaultLimit:  defaultLimit,
		blockDuration: time.Duration(blockDuration) * time.Second,
		tokenLimits:   tokenLimits,
		ctx:           context.Background(),
	}
}

func (r *RedisRateLimiter) Allow(identifier string) (bool, time.Duration, error) {
	blockKey := fmt.Sprintf("block:%s", identifier)
	limitKey := fmt.Sprintf("rate:%s", identifier)

	script := redis.NewScript(`
		local blockKey = KEYS[1]
		local rateKey = KEYS[2]
		local limit = tonumber(ARGV[1])
		local blockDuration = tonumber(ARGV[2])

		if redis.call("TTL", blockKey) > 0 then
			return {0, redis.call("TTL", blockKey)}
		end

		local current = redis.call("INCR", rateKey)
		if current == 1 then
			redis.call("EXPIRE", rateKey, blockDuration)
		end

		if current > limit then
			redis.call("SET", blockKey, "1", "EX", blockDuration)
			return {0, blockDuration}
		end

		return {1, 0}
	`)

	limit := r.defaultLimit
	if val, ok := r.tokenLimits[identifier]; ok {
		limit = val
	}

	res, err := script.Run(r.ctx, r.client, []string{blockKey, limitKey},
		limit,
		int(r.blockDuration.Seconds()),
	).Result()

	if err != nil {
		fmt.Println("❌ Redis script error:", err)
		return false, 0, err
	}

	data := res.([]interface{})
	allowed := data[0].(int64) == 1
	retryAfter := time.Duration(data[1].(int64)) * time.Second

	if allowed {
		fmt.Printf("🔁 Incrementando %s (limite: %d)\n", limitKey, limit)
	} else {
		fmt.Printf("⛔ BLOQUEADO %s (aguarde %v)\n", identifier, retryAfter)
	}

	return allowed, retryAfter, nil
}
