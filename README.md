# 🚦 Rate Limiter in Go (IP and Token)

This project implements a configurable rate limiter in Go, capable of controlling requests per second based on the IP address or an access token sent via the `API_KEY` header. The logic is based on middleware, using the Gin framework and storing data in Redis. If a token is present, its rate limiting configurations must override those of the IP. The project is prepared to run via Docker and features a decoupled storage strategy, allowing for the future use of other solutions besides Redis.

## 🧱 Project Structure

The project is organized as follows:

- `cmd/server/`: application entry point
- `internal/config/`: environment variable reading
- `internal/limiter/`: limiting logic and Redis abstraction
- `internal/middleware/`: rate limiting middleware for Gin
- `test/`: automated tests
- `Dockerfile` and `docker-compose.yml`: orchestration with Redis
- `.env`: limiting configurations
- `README.md`: documentation

## 🚀 How to run the application

🐳 With Docker installed, simply run the following command:

\`\`\`bash
docker-compose up --build
\`\`\`

The application will be available at `http://localhost:8080` and Redis will automatically start at `localhost:6379`.

## ⚙️ Configuration via .env

The application can be configured via environment variables in the `.env` file, as shown in the example below:

\`\`\`env
REDIS_HOST=localhost:6379
IP_RATE_LIMIT=5
IP_BLOCK_DURATION=300
TOKEN_RATE_LIMIT=10
TOKEN_BLOCK_DURATION=300
\`\`\`

These configurations define limits per IP (5 requests per second) and per token (10 requests per second). If the limit is exceeded, the IP or token will be blocked for 300 seconds (5 minutes). The token logic overrides the IP logic.

## 🧪 Automated Tests

The automated tests validate the limits per IP and per token. To run them, with Redis running, use:

\`\`\`bash
go test -v ./...
\`\`\`

Or for a specific package:

\`\`\`bash
go test -v ./internal/middleware
\`\`\`

## 🧪 How to test manually

You can use `curl` to test the API with or without a token. Examples:

Request with token:

\`\`\`bash
curl -H "API_KEY: abc123" http://localhost:8080
\`\`\`

Request without token (uses the IP as the limiting key):

\`\`\`bash
curl http://localhost:8080
\`\`\`

If the limit is exceeded, the response will be:

- **Status**: 429 Too Many Requests
- **Message**: `you have reached the maximum number of requests or actions allowed within a certain time frame`

## ♻️ Persistence Strategy

The rate limiting logic uses Redis by default. However, an interface named `RateLimiter` allows implementing new storage strategies without impacting the rest of the application, simply by replacing the current implementation (`RedisRateLimiter`) with another of your choice, such as a relational database, local memory, or external services.

## 📚 Technologies Used

- Go (1.21+)
- Gin Web Framework
- Redis (via go-redis v9)
- Docker and Docker Compose
- Tests with `testing` and `httptest` package

## 👨‍💻 Autor

Developed as a technical challenge by vitorlrrcamargo.
