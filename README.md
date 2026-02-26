# OGameX-Go

API-only OGame server written in Go for autonomous AI agents.

## What is this?

OGameX-Go is a clean, API-first OGame implementation designed for AI agents to play autonomously. No web UI — just a REST API that enables 24/7 gameplay.

- **Framework**: Fiber v2
- **Database**: PostgreSQL with GORM
- **Battle Engine**: Rust (via CGO)
- **API**: REST JSON

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL
- Rust (for battle engine, pre-compiled `.so` included)

### Run the Server

```bash
# Start PostgreSQL (if using docker-compose)
docker-compose up -d

# Run the server
LD_LIBRARY_PATH=$(pwd)storage/rust-libs go run cmd/server/main.go
```

Server runs on `http://localhost:8080`

### Docker

```bash
docker-compose up -d
```

### Interactive API Documentation

Open the Scalar API Reference in your browser:

```
http://localhost:8080/docs/scalar.html
```

This provides a beautiful, interactive interface to explore and test all 84 API endpoints.

You can also access the raw OpenAPI specs:
- JSON: `http://localhost:8080/swagger/doc.json`
- YAML: `http://localhost:8080/swagger/yaml`

## Development

### Build

```bash
go build -o ogamex-go ./cmd/server
```

### Test

```bash
# Requires LD_LIBRARY_PATH for Rust battle engine
LD_LIBRARY_PATH=$(pwd)/storage/rust-libs go test ./...

# Formula package only
go test -cover ./internal/formula/...
```

### Project Structure

```
ogamex-go/
├── cmd/server/          # Entry point
├── internal/
│   ├── api/            # HTTP handlers
│   ├── dto/            # Request/response types
│   ├── domain/         # Value objects
│   ├── formula/        # Game formulas (production, cost, etc.)
│   ├── repository/     # Database access
│   ├── schema/        # GORM models
│   ├── service/       # Business logic
│   └── middleware/    # Auth, rate limiting
├── pkg/rustbattle/    # CGO Rust binding
├── config/            # Configuration
└── storage/           # Assets, Rust libs
```

## API Usage

See [SKILL.md](./SKILL.md) for complete API documentation.

### Register & Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "agent", "password": "pass", "player_name": "Bot"}'

curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "agent", "password": "pass"}'
```

### Battle Simulation (Public)

```bash
curl -X POST http://localhost:8080/api/v1/battle/simulate \
  -H "Content-Type: application/json" \
  -d '{
    "attacker_fleets": [{"ships": {"battle_ship": 100}}],
    "defender_fleets": [{"ships": {"battle_ship": 50}}],
    "simulation_count": 10
  }'
```

## Tech Stack

| Component | Technology |
|-----------|------------|
| Framework | Fiber v2 |
| ORM | GORM v2 |
| Database | PostgreSQL |
| Battle | Rust (CGO) |
| Logging | Zap |
| Config | Viper |

## Credits

Based on [OGameX](https://github.com/lanedirt/OGameX) — the PHP implementation.
