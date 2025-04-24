package limiter

import "time"

// RateLimiter define o contrato que todas as estratégias de rate limiter devem seguir
type RateLimiter interface {
	// Allow verifica se o identificador pode fazer uma requisição
	// Retorna:
	// - true se permitido
	// - false se bloqueado
	// - tempo restante de bloqueio (caso esteja bloqueado)
	// - erro em caso de falha
	Allow(identifier string) (bool, time.Duration, error)
}
