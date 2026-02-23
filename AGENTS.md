**AGENTS.md**

```markdown
# AGENTS.md — Guide for AI Coding Agents: Porting OGameX to Clean Go API

**Version:** 1.0.0 (Feb 2026)  
**Project:** ogamex-go — A minimal, API-only, agent-native OGame server in Go.  
**Repo:** (your fork)  
**Original:** https://github.com/lanedirt/OGameX  
**Goal:** Build the cleanest, fastest, most maintainable OGame backend ever — explicitly designed so that autonomous AI agents (OpenClaw/Moltbot etc.) can play 24/7 with zero friction. No frontend, no Laravel magic, no PHP verbosity.

You are an elite Go coding agent. Follow this document **exactly**. Deviate only after explicit approval and a new section added here.

---

## 1. Mission & Non-Negotiables

- **API-first, agent-first**: Every endpoint must return clean, predictable JSON. Agents poll `/api/status` every 5-30 min.
- **Rust battle engine stays untouched** — call it via cgo. Do **not** reimplement combat.
- **100% TDD** from day 1.
- **Logging-driven development**: Every decision, every calculation, every side-effect must be logged at the right level.
- **Production-ready for agents**: rate limiting, structured logs (JSON), graceful shutdown, observability hooks.
- **Zero magic**: explicit, readable, testable code. No global state, no hidden callbacks.

---

## 2. Target Tech Stack (fixed)

| Layer          | Choice                  | Why (for agents)                          |
|----------------|-------------------------|-------------------------------------------|
| Framework      | Fiber v2 (`go-fiber`)   | Fastest, zero-allocation routes, perfect JSON |
| ORM / DB       | GORM v2 + PostgreSQL    | Simple migrations, but we also keep raw sqlc option |
| Battle         | cgo → Rust `.so`        | 200× faster, already battle-tested        |
| Scheduler      | `github.com/robfig/cron/v3` + in-memory queue + Redis option |
| Logging        | `go.uber.org/zap` (production JSON) + `slog` fallback |
| Config         | `viper` + env + `godotenv` |
| Testing        | `testing` + `testify` + `httptest` + `gofakeit` |
| Docker         | Alpine + multi-stage    | < 50 MB image                             |

---

## 3. Repository Structure (strictly follow)

```bash
ogamex-go/
├── cmd/server/main.go
├── internal/
│   ├── domain/          # pure value objects (Planet, Fleet, Resources, etc.)
│   ├── service/         # BuildingService, FleetService, ResearchService...
│   ├── engine/          # rustbattle.go (cgo wrapper)
│   ├── repository/      # gorm repos or sqlc queries
│   ├── scheduler/       # fleet arrival, production queues
│   ├── formula/         # pure functions for EVERY OGame formula
│   ├── api/             # Fiber handlers (thin)
│   ├── dto/             # request/response structs only
│   └── middleware/      # auth, rate-limit, logging
├── pkg/
│   └── rustbattle/      # cgo binding (single file)
├── migrations/          # goose or GORM auto
├── tests/               # integration + load tests
├── config/
├── docker-compose.yml
├── go.mod
└── AGENTS.md            # ← you are here
```

---

## 4. Go Best Practices (mandatory)

- **Package organization**: small, focused packages. One responsibility per file when possible.
- **Error handling**: `errors.Is`, `errors.As`, custom sentinel errors in `internal/errors`.
- **Context everywhere**: all service methods take `ctx context.Context`.
- **Immutability**: value types for Resources, Coordinates, etc.
- **Interfaces for testability**: `type PlanetRepository interface { ... }`
- **No `interface{}`**, no `any` except in logging.
- **Constants & iota** for enums (MissionType, UnitType, etc.).
- **Generics** where they actually reduce duplication (e.g., cost calculators).
- **Zero dependencies** unless justified in this doc.
- **Use go-related skills liberally**: golang-pro and golang-patterns. Also use api-design-principles for API development and postgresql-table-design for postgres.

---

## 5. Test-Driven Development (TDD) — NON-NEGOTIABLE

Rule: **Red → Green → Refactor** for every new function.

1. Write failing test first.
2. Make it pass with minimal code.
3. Refactor.
4. Commit with message: `feat: <thing> (TDD)`.

**Test coverage requirements:**
- 95%+ on `formula/`, `domain/`, `service/`
- Use table-driven tests for all formulas (feed OGameX test cases).
- Integration tests spin up real PostgreSQL + Rust battle via Docker (testcontainers-go).

Example test file name: `formula_production_test.go`

**Wiki Verification Tests:** See [WIKI_TEST_ROADMAP.md](./WIKI_TEST_ROADMAP.md) for comprehensive cross-checking of our implementation against the OGame wiki documentation. Run with `go test -v -run Wiki ./...`.

---

## 6. Logging-Driven Development

Use `zap` with structured fields. Every important action must log:

```go
logger.Info("building_upgrade_started",
    zap.String("player_id", playerID),
    zap.String("planet_id", planetID),
    zap.String("building", "metal_mine"),
    zap.Int("level", level),
    zap.Int64("metal_cost", cost.Metal),
    zap.Duration("time", duration),
)
```

**Levels:**
- Debug: calculation details (formulas)
- Info: player-visible actions (build, fleet send, battle start)
- Warn: soft failures (insufficient resources, fleet recall)
- Error: hard failures (DB, Rust panic, invalid state)
- Always include `trace_id` (from Fiber middleware)

Log every fleet arrival, every queue tick, every battle outcome.

---

## 7. Reference Original OGameX Files (study these first)

**Rust battle engine (DO NOT TOUCH, just bind):**
- https://github.com/lanedirt/OGameX/tree/main/rust/battle_engine_ffi
- Core lib: https://github.com/lanedirt/OGameX/blob/main/rust/battle_engine_ffi/src/lib.rs  (functions: `SimulateBattle`, `CalculateDebris`, etc.)
- Build script: https://github.com/lanedirt/OGameX/blob/main/rust/compile.sh

**PHP logic to port (reference for exact formulas & edge cases):**
- All services directory: https://github.com/lanedirt/OGameX/tree/main/app/Services
- Queue logic: look for `QueueService.php`, `UnitQueue`, dark-matter halving (recent changes)
- Resource production: search for `calculateProduction` or similar in Building/Planet models
- Fleet missions: https://github.com/lanedirt/OGameX/tree/main/app/Services (FleetService)
- Database schema (exact tables & fields): https://github.com/lanedirt/OGameX/tree/main/database/migrations
  - Pay special attention to `unit_queues`, `fleets`, `planets`, `research`, `dm_halved` flags

**Tests to replicate:**
- https://github.com/lanedirt/OGameX/tree/main/tests (all expedition, battle, queue tests are gold)

Study the PHP → extract **exact integer math** (rounding, flooring, caps). Reproduce 1:1 in Go `formula/` package.

---

## 8. Common Pitfalls & How to Avoid Them

| Pitfall                              | Why it hurts agents                  | Prevention |
|--------------------------------------|--------------------------------------|----------|
| Floating-point in formulas           | Resource drift over weeks            | Use `int64` + explicit flooring everywhere |
| Timezone / server time drift         | Fleets arrive at wrong ticks         | Always use UTC + `time.Time` with monotonic clock |
| Queue race conditions                | Double-spend resources               | Single goroutine per planet + DB transactions |
| Battle result non-determinism        | Agent distrust                       | Seed Rust engine deterministically + log full input |
| CGO thread safety                    | Crashes under load                   | One global battle pool + mutex |
| Missing edge cases (expeditions, ACS, moon chance) | Broken late-game | Every formula test case must include OGameX wiki + real server examples |
| Over-polling without caching         | High DB load                         | Redis cache layer for `/status` (TTL 10s) |
| Poor error messages in API           | Agents retry forever                 | `error_code` + `human_message` + `details` in every error response |

---

## 9. Implementation Roadmap (phased — follow order)

**Phase 0** — Setup & Rust binding (1 day)
**Phase 1** — Domain + Formula package (TDD every function) (2-3 days)
**Phase 2** — Repository + Services (Building, Research, Shipyard) (3 days)
**Phase 3** — Scheduler + Queue engine (4 days)
**Phase 4** — Fleet missions + Battle integration (3 days)
**Phase 5** — Full REST API + auth + rate limiting (2 days)
**Phase 6** — Agent skill compatibility + load testing (2 days)

After each phase: run full TDD suite + manual agent test (curl + simple Go agent).

---

## 10. API Design Rules (for agents)

- All endpoints under `/api/v1/`
- Auth: Bearer token (Sanctum-like)
- Consistent naming: `POST /buildings/upgrade`, `POST /fleets/send`
- Response shape: `{ "data": {...}, "meta": {...} }` or `{ "error": {...} }`
- Pagination, filtering, timestamps in RFC3339
- Rate limit headers: `X-RateLimit-Remaining`

---

## 11. Final Checklist Before PR

- [ ] 95%+ test coverage
- [ ] All logs structured + include trace_id
- [ ] Rust battle integration tested with 1M-unit fleets
- [ ] Matches OGameX behavior on 10+ known test cases (wiki + OGameX tests)
- [ ] Docker image < 60 MB
- [ ] Agents can play full game using only this API (no web UI needed)

---

## 12. Development Progress (UPDATE REGULARLY)

**Last Updated:** Feb 2026

### Overall Progress: ~60% Complete

### Completed Components

| Component | Status | Notes |
|-----------|--------|-------|
| Tech Stack | 95% | Fiber v2, GORM, Zap, Viper, Docker |
| Domain Types | 80% | Building, Unit, Coordinates, Resources, Mission, Research types |
| Formula Package | 85% | Production, cost, distance, ship stats with wiki verification |
| Wiki Tests | 100% | All 10 phases complete (WIKI_TEST_ROADMAP.md) |
| Auth Service | 80% | Register, login, token validation |
| Building Service | 70% | Start/cancel building, queue management |
| Production Service | 75% | Resource production calculation |
| Research Service | 70% | Start/cancel research, tech tree |
| Unit Service | 70% | Ship/defense building, queues |
| Fleet Service | 75% | Most missions implemented, battle uses Go not Rust |
| Repositories | 70% | User, planet, fleet, queues, tech |
| Scheduler | 80% | Building, research, unit, fleet, production processing |
| API Endpoints | 60% | ~25 endpoints, auth, rate limiting |
| Rust Battle Engine | 60% | CGO binding, SimulateBattleWithRust added - needs call from processAttack |

### Known Issues & Gaps

| Issue | Severity | Status |
|-------|----------|--------|
| Test failures in `cost_test.go` | HIGH | FIXED |
| Missing DTOs | MEDIUM | FIXED |
| Missing custom middleware | MEDIUM | FIXED |
| Fleet missions | MEDIUM | MOSTLY DONE - attack/transport/colonize/recycle/expedition implemented |
| Rust battle integration | HIGH | DONE - added SimulateBattleWithRust method, needs to be called from processAttack |
| Test coverage | MEDIUM | 76.5% on formula/ - need 95%+ |
| Agent compatibility | LOW | Phase 6 not started |

### Repository Structure Status

```
ogamex-go/
├── cmd/server/main.go           ✅ Complete
├── internal/
│   ├── domain/                 ✅ 6 files - building, unit, coordinates, resources, mission, research types
│   ├── service/                ✅ 6 services - auth, building, fleet, production, research, unit
│   ├── engine/                 ❌ Empty - rustbattle is in pkg/
│   ├── repository/             ✅ 6 repos - user, planet, fleet, queues, tech, interfaces
│   ├── scheduler/              ⚠️ Partial - basic cron, needs fleet processing
│   ├── formula/                ✅ 16 files - production, cost, distance, ship_stats + wiki tests
│   ├── api/                    ✅ handlers.go + errors.go
│   ├── dto/                    ✅ dto.go - request/response structs
│   ├── middleware/             ✅ auth.go - custom middleware
│   ├── logger/                 ✅ logger.go
│   └── schema/                 ✅ models.go (GORM)
├── pkg/rustbattle/             ✅ battle.go (CGO binding)
├── tests/wiki/                 ✅ 8 JSON test data files
├── config/                    ✅ config.yaml
├── docker-compose.yml          ✅ PostgreSQL + app
├── Dockerfile                 ✅ Multi-stage build
└── AGENTS.md                  📝 This file
```

### Test Status

```bash
# Run wiki verification tests
go test -v -run Wiki ./...

# Run all tests (requires LD_LIBRARY_PATH for rustbattle)
LD_LIBRARY_PATH=/home/bolt/Documents/ogamex-go/storage/rust-libs go test ./...

# Check coverage
go test -cover ./internal/formula/...
```

### Next Steps (Priority Order)

1. **FIX:** Test failures in `cost_test.go` (TestCalculatePositionBonus)
2. **ADD:** DTOs for request/response structs
3. **ADD:** Custom middleware for auth/rate-limit
4. **COMPLETE:** Fleet mission types (attack, transport, colonize, etc.)
5. **INTEGRATE:** Rust battle engine with fleet service
6. **IMPROVE:** Test coverage to 95%+
7. **ADD:** Redis caching for /api/status endpoint

---

*NOTE: Update this section as development progresses. Run `go test ./...` before each commit to ensure no regressions.*

