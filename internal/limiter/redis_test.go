package limiter

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis(t *testing.T) *RedisRateLimiter {
	// Usa Redis local ou do docker-compose
	os.Setenv("REDIS_HOST", "localhost:6379")
	os.Setenv("DEFAULT_LIMIT", "3")
	os.Setenv("BLOCK_DURATION", "10")
	os.Setenv("TOKEN_LIMIT_abc123", "5")

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9, // banco separado só para testes
	})

	// Limpa banco de testes antes de rodar
	err := rdb.FlushDB(context.Background()).Err()
	if err != nil {
		t.Fatalf("Erro ao limpar Redis de teste: %v", err)
	}

	return NewRedisRateLimiter()
}

func TestAllow_DefaultLimit(t *testing.T) {
	limiter := setupTestRedis(t)
	identifier := "127.0.0.1"

	// Permite 3 requisições
	for i := 0; i < 3; i++ {
		allowed, _, err := limiter.Allow(identifier)
		assert.NoError(t, err)
		assert.True(t, allowed, "Requisição %d deveria ser permitida", i+1)
	}

	// 4ª requisição deve ser bloqueada
	allowed, retryAfter, err := limiter.Allow(identifier)
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.GreaterOrEqual(t, retryAfter.Seconds(), 1.0)
}

func TestAllow_TokenOverride(t *testing.T) {
	limiter := setupTestRedis(t)
	token := "abc123"

	// Permite até 5 requisições com token
	for i := 0; i < 5; i++ {
		allowed, _, err := limiter.Allow(token)
		assert.NoError(t, err)
		assert.True(t, allowed, "Token request %d deveria ser permitida", i+1)
	}

	// 6ª requisição deve ser bloqueada
	allowed, retryAfter, err := limiter.Allow(token)
	assert.NoError(t, err)
	assert.False(t, allowed)
	assert.GreaterOrEqual(t, retryAfter.Seconds(), 1.0)
}

func TestAllow_UnblockAfterDuration(t *testing.T) {
	limiter := setupTestRedis(t)
	identifier := "testuser"

	// Estoura limite
	for i := 0; i < 4; i++ {
		limiter.Allow(identifier)
	}

	// Espera até fim do bloqueio
	time.Sleep(11 * time.Second)

	allowed, _, err := limiter.Allow(identifier)
	assert.NoError(t, err)
	assert.True(t, allowed, "Deveria estar desbloqueado após duração do bloqueio")
}
