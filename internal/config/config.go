package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DefaultLimit   int
	BlockDuration  time.Duration
	RedisHost      string
	Port           string
	TokenOverrides map[string]int
}

func LoadConfig() Config {
	defaultLimit, _ := strconv.Atoi(os.Getenv("DEFAULT_LIMIT"))
	blockSecs, _ := strconv.Atoi(os.Getenv("BLOCK_DURATION"))

	// Extra tokens
	tokenOverrides := make(map[string]int)
	for _, env := range os.Environ() {
		if len(env) > 12 && env[:12] == "TOKEN_LIMIT_" {
			key := env[12:]
			val := os.Getenv("TOKEN_LIMIT_" + key)
			limit, _ := strconv.Atoi(val)
			tokenOverrides[key] = limit
		}
	}

	return Config{
		DefaultLimit:   defaultLimit,
		BlockDuration:  time.Duration(blockSecs) * time.Second,
		RedisHost:      os.Getenv("REDIS_HOST"),
		Port:           os.Getenv("PORT"),
		TokenOverrides: tokenOverrides,
	}
}
