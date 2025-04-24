package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"rate-limiter-go/internal/http"
	"rate-limiter-go/internal/limiter"
	"rate-limiter-go/internal/middleware"
)

func main() {
	// Carrega variáveis de ambiente do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file not found, using system environment variables")
	}

	// Inicializa o rate limiter
	rateLimiter := limiter.NewRedisRateLimiter()

	// Inicializa o servidor web
	router := gin.Default()
	router.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Rota de teste
	router.GET("/", http.RootHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor rodando na porta %s...", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
