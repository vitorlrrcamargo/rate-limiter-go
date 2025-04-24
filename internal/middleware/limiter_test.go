package middleware_test

import (
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"rate-limiter-go/internal/limiter"
	"rate-limiter-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	_ = godotenv.Load("../../.env")
}

func setupRouter(l limiter.RateLimiter) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.RateLimitMiddleware(l))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})
	return r
}

func TestIPRateLimiting(t *testing.T) {
	os.Setenv("DEFAULT_LIMIT", "3") // 3 req/s
	os.Setenv("BLOCK_DURATION", "5")

	lim := limiter.NewRedisRateLimiter()
	router := setupRouter(lim)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	}

	// 4ª requisição deve falhar
	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 429, resp.Code)
	assert.Contains(t, resp.Body.String(), "you have reached the maximum")
}

func TestTokenRateLimiting(t *testing.T) {
	os.Setenv("DEFAULT_LIMIT", "1")
	os.Setenv("TOKEN_LIMIT_testtoken", "5")
	os.Setenv("BLOCK_DURATION", "5")

	lim := limiter.NewRedisRateLimiter()
	router := setupRouter(lim)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "testtoken")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, 200, resp.Code)
	}

	// 6ª requisição deve falhar
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "testtoken")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 429, resp.Code)
}

func TestUnblockAfterExpiration(t *testing.T) {
	os.Setenv("DEFAULT_LIMIT", "1")
	os.Setenv("BLOCK_DURATION", "2")

	lim := limiter.NewRedisRateLimiter()
	router := setupRouter(lim)

	// Primeira: passa
	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)

	// Segunda: bloqueia
	req = httptest.NewRequest("GET", "/", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 429, resp.Code)

	// Espera expirar
	time.Sleep(3 * time.Second)

	// Deve passar de novo
	req = httptest.NewRequest("GET", "/", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, 200, resp.Code)
}
