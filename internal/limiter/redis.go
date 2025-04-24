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

	// Verifica se está bloqueado
	blocked, err := r.client.TTL(r.ctx, blockKey).Result()
	if err != nil {
		return false, 0, err
	}
	if blocked > 0 {
		return false, blocked, nil
	}

	// Determina o limite
	limit := r.defaultLimit
	if tokenLimit, exists := r.tokenLimits[identifier]; exists {
		limit = tokenLimit
	}

	// Incrementa contagem
	count, err := r.client.Incr(r.ctx, limitKey).Result()
	if err != nil {
		return false, 0, err
	}
	if count == 1 {
		r.client.Expire(r.ctx, limitKey, time.Second)
	}

	if count > int64(limit) {
		// Bloqueia o acesso
		r.client.Set(r.ctx, blockKey, "1", r.blockDuration)
		return false, r.blockDuration, nil
	}

	return true, 0, nil
}
