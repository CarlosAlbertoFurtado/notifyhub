# NotifyHub

![Go](https://img.shields.io/badge/go-1.23-00ADD8)
![CI](https://github.com/CarlosAlbertoFurtado/notifyhub/actions/workflows/ci.yml/badge.svg)

Microsserviço de notificações. Recebe HTTP, persiste no Postgres, dispara email/SMS/webhook via goroutine. Resposta volta na hora com status `queued`.

Extrai a lógica de notificação que eu tava copiando entre o [FinTrack](https://github.com/CarlosAlbertoFurtado/fintrack-api) e o [SmartBooking](https://github.com/CarlosAlbertoFurtado/smart-booking). Agora é um serviço separado que qualquer um chama via HTTP.

## Endpoints

```
POST  /api/notifications/send     envia (email, sms ou webhook)
GET   /api/notifications           lista com filtro (?channel=email&status=sent&page=1)
GET   /api/notifications/stats     contagem por status
GET   /api/notifications/:id       detalhe
GET   /health
```

### Exemplo

```bash
curl -X POST http://localhost:8080/api/notifications/send \
  -H "Content-Type: application/json" \
  -d '{"channel":"email","recipient":"user@example.com","subject":"Oi","body":"teste"}'
```

## Como roda

```bash
cp .env.example .env
docker-compose up -d      # ou: go run ./cmd/api (precisa PG rodando)
```

## Estrutura

```
cmd/api/           entrypoint
internal/
├── domain/        entidade Notification + interfaces
├── application/   use case de envio (goroutine)
├── config/        env vars
├── handler/       HTTP handlers (Gin)
└── infra/         repo postgres, senders (email, webhook)
```

Go 1.23, Gin, pgx (driver nativo, sem ORM), goroutines pra dispatch, Docker multi-stage (~15MB final), GitHub Actions com `-race`.

## Testes

`go test ./... -v -race` — 8 testes, cobertura ~71%.

## Pendente

- [ ] Template com variáveis (tipo `{{nome}}`)
- [ ] Retry com backoff

MIT
