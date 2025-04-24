# 🚦 Rate Limiter em Go (IP e Token)

Este projeto implementa um rate limiter configurável em Go, capaz de controlar requisições por segundo com base em endereço IP ou token de acesso enviado via header \`API_KEY\`. A lógica é baseada em middleware, utilizando o framework Gin e armazenando dados no Redis. Caso um token esteja presente, suas configurações de limitação devem sobrepor as do IP. O projeto está preparado para funcionar via Docker e possui uma estratégia desacoplada de armazenamento, podendo futuramente utilizar outra solução além do Redis.

## 🧱 Estrutura do Projeto

O projeto está organizado da seguinte forma:

- \`cmd/server/\`: ponto de entrada da aplicação
- \`internal/config/\`: leitura de variáveis de ambiente
- \`internal/limiter/\`: lógica de limitação e abstração com Redis
- \`internal/middleware/\`: middleware de rate limiting para o Gin
- \`test/\`: testes automatizados
- \`Dockerfile\` e \`docker-compose.yml\`: orquestração com Redis
- \`.env\`: configurações de limitação
- \`README.md\`: documentação

## 🚀 Como executar a aplicação

🐳 Com Docker instalado, basta executar o seguinte comando:

\`\`\`bash
docker-compose up --build
\`\`\`

A aplicação ficará disponível em \`http://localhost:8080\` e o Redis será iniciado automaticamente em \`localhost:6379\`.

## ⚙️ Configuração via .env

A aplicação pode ser configurada por variáveis de ambiente no arquivo \`.env\`, conforme exemplo abaixo:

\`\`\`env
REDIS_HOST=localhost:6379
IP_RATE_LIMIT=5
IP_BLOCK_DURATION=300
TOKEN_RATE_LIMIT=10
TOKEN_BLOCK_DURATION=300
\`\`\`

Essas configurações definem limites por IP (5 requisições por segundo) e por token (10 requisições por segundo). Caso o limite seja excedido, o IP ou token será bloqueado por 300 segundos (5 minutos). A lógica do token se sobrepõe à do IP.

## 🧪 Testes Automatizados

Os testes automatizados validam os limites por IP e por token. Para executá-los, com Redis rodando, utilize:

\`\`\`bash
go test -v ./...
\`\`\`

Ou para um pacote específico:

\`\`\`bash
go test -v ./internal/middleware
\`\`\`

## 🧪 Como testar manualmente

Você pode usar \`curl\` para testar a API com ou sem token. Exemplos:

Requisição com token:

\`\`\`bash
curl -H "API_KEY: abc123" http://localhost:8080
\`\`\`

Requisição sem token (usa o IP como chave de limitação):

\`\`\`bash
curl http://localhost:8080
\`\`\`

Caso o limite seja excedido, a resposta será:

- **Status**: 429 Too Many Requests
- **Mensagem**: \`you have reached the maximum number of requests or actions allowed within a certain time frame\`

## ♻️ Estratégia de Persistência

A lógica de rate limiting utiliza Redis por padrão. No entanto, uma interface chamada \`RateLimiter\` permite implementar novas estratégias de armazenamento sem impactar o restante da aplicação, bastando substituir a implementação atual (\`RedisRateLimiter\`) por outra de sua escolha, como banco relacional, memória local ou serviços externos.

## 📚 Tecnologias Utilizadas

- Go (1.21+)
- Gin Web Framework
- Redis (via go-redis v9)
- Docker e Docker Compose
- Testes com pacote \`testing\` e \`httptest\`

## 👨‍💻 Autor

Desenvolvido como desafio técnico por vitorlrrcamargo.