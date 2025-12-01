# ShortfyURL - Encurtador de Links Simples

Encurtador de URLs em Go com Cassandra, Redis e codificação Base62.

## Pré-requisitos

- Docker
- Docker Compose

## Executar com Docker

```bash
# Subir todos os serviços
docker-compose up -d

# Aguardar 30 segundos e criar o keyspace
bash init-cassandra.sh

# Ou no Windows PowerShell:
# Start-Sleep -Seconds 30
# docker exec -it shortfy-cassandra cqlsh -e "CREATE KEYSPACE IF NOT EXISTS shortfy WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};"
```

## Acessar a aplicação

- **Frontend**: http://localhost:8080
- **Admin**: http://localhost:8080/pages/admin.html

## Executar localmente (sem Docker)

### Pré-requisitos
- Go 1.21+
- Cassandra 3.x+
- Redis 6.x+

```bash
cd backend
go mod download
go run .
```

## Endpoints

### Criar URL curta
```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://exemplo.com/url-muito-longa"}'
```

### Redirecionar
```bash
curl http://localhost:8080/{shortCode}
```

### Estatísticas
```bash
curl http://localhost:8080/api/stats/{shortCode}
```

## Parar os serviços

```bash
docker-compose down

# Remover volumes (dados)
docker-compose down -v
```
