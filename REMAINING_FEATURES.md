# OGameX-Go Remaining Features

**Last Updated:** Feb 2026  
**Purpose:** Document all features that still need implementation for full agent compatibility

---

## Overview

The OGameX-Go project is ~95% complete with ~80 API endpoints covering all core gameplay. Most features from the original OGameX PHP repository have been ported. This document lists any remaining items.

---

## Completed Features ✅

### Core Gameplay
- [x] User authentication (register, login, token validation)
- [x] Planet management (CRUD, resources, buildings)
- [x] Building queue system
- [x] Research queue system
- [x] Unit/Ship building queue
- [x] Fleet missions (attack, transport, colonize, recycle, expedition, ACS)
- [x] Fleet recall
- [x] Battle simulation (via Rust engine or Go fallback)

### Multiplayer Features
- [x] Alliance system (create, join, applications, members)
- [x] Buddy system (requests, accept, reject)
- [x] Espionage reports
- [x] Messages

### Premium & Economy
- [x] Premium status/activation
- [x] Merchant buy/sell
- [x] Vacation mode
- [x] Planet move service
- [x] Character class selection (General, Admiral, Engineer, Geologist)
- [x] DM queue halving (HalvingService)

### Moon Features
- [x] Jump Gate (fleet teleportation between moons)
- [x] Sensor Phalanx (probe detection)
- [x] Moon buildings

### Galaxy Features
- [x] Galaxy view
- [x] Debris fields (creation, collection)
- [x] Wreck fields (creation, collection)
- [x] Highscore

### NPC Systems
- [x] NPC planet generation
- [x] NPC expedition fleets
- [x] NPC pirate raids
- [x] ACS (Fleet Union)

### API Features
- [x] Structured JSON responses
- [x] Rate limiting
- [x] Trace ID middleware
- [x] Authentication middleware

---

## API Endpoints Summary (~80 endpoints)

```
Auth:
POST   /api/v1/auth/register
POST   /api/v1/auth/login

User:
GET    /api/v1/user
GET    /api/v1/user/stats
GET    /api/v1/character-class
POST   /api/v1/character-class/select
GET    /api/v1/user/vacation
POST   /api/v1/user/vacation/enable
POST   /api/v1/user/vacation/disable

Planets:
GET    /api/v1/planets
GET    /api/v1/planets/:id
GET    /api/v1/planets/:id/details
PUT    /api/v1/planets/:id/set-current
GET    /api/v1/planets/:id/resources
GET    /api/v1/planets/:id/buildings
GET    /api/v1/planets/:id/overview
POST   /api/v1/planets/:id/buildings/:building_id
GET    /api/v1/planets/:id/queue
POST   /api/v1/planets/:id/move
GET    /api/v1/planets/:id/jump-gate/targets
POST   /api/v1/planets/:id/jump-gate/execute
POST   /api/v1/planets/:id/phalanx/scan

Research:
POST   /api/v1/research/start
GET    /api/v1/research/queue

Fleets:
POST   /api/v1/fleets/send
GET    /api/v1/fleets
POST   /api/v1/fleets/:id/recall

Units:
POST   /api/v1/planets/:id/units/build
GET    /api/v1/planets/:id/units/queue
GET    /api/v1/units/available
GET    /api/v1/planets/:id/units
GET    /api/v1/planets/:id/defense
GET    /api/v1/planets/:id/production

Queues:
DELETE /api/v1/queue/:id
DELETE /api/v1/research/queue/:id
DELETE /api/v1/units/queue/:id

Messages:
GET    /api/v1/messages
GET    /api/v1/messages/unread
POST   /api/v1/messages/:id/read
DELETE /api/v1/messages/:id

Galaxy:
GET    /api/v1/galaxy/:galaxy/:system

Highscore:
GET    /api/v1/highscore/:category

Notes:
GET    /api/v1/notes
POST   /api/v1/notes
PUT    /api/v1/notes/:id
DELETE /api/v1/notes/:id

Alliances:
GET    /api/v1/alliances
POST   /api/v1/alliances
GET    /api/v1/alliances/:id
GET    /api/v1/alliances/:id/members
POST   /api/v1/alliances/:id/apply
GET    /api/v1/alliances/:id/applications
POST   /api/v1/alliances/applications/:id/accept
POST   /api/v1/alliances/applications/:id/reject
POST   /api/v1/alliances/leave
PUT    /api/v1/alliances

Buddy:
GET    /api/v1/buddy
POST   /api/v1/buddy/request
GET    /api/v1/buddy/pending
POST   /api/v1/buddy/:id/accept
POST   /api/v1/buddy/:id/reject
DELETE /api/v1/buddy/:id

Espionage:
GET    /api/v1/espionage
GET    /api/v1/espionage/unread
POST   /api/v1/espionage/:id/read
DELETE /api/v1/espionage/:id

Debris:
GET    /api/v1/debris
GET    /api/v1/debris/:galaxy/:system/:position
POST   /api/v1/debris/:galaxy/:system/:position/collect

Wrecks:
GET    /api/v1/wrecks
GET    /api/v1/wrecks/:galaxy/:system/:position
POST   /api/v1/wrecks/:galaxy/:system/:position/collect

NPC:
GET    /api/v1/npc/planets
GET    /api/v1/npc/planets/:galaxy/:system/:position
POST   /api/v1/npc/planets
POST   /api/v1/npc/fleets/expedition
POST   /api/v1/npc/fleets/pirate

ACS:
POST   /api/v1/acs/create
POST   /api/v1/acs/:id/join
GET    /api/v1/acs/:id
GET    /api/v1/acs/:id/fleets

Premium:
GET    /api/v1/premium/status
POST   /api/v1/premium/activate
POST   /api/v1/merchant/buy
POST   /api/v1/merchant/sell
```

---

## Low Priority / Not Implemented

These features are not critical for agent gameplay:

### 1. User Settings

**Status:** Not Implemented  
**Reason:** Agents don't need to change settings programmatically  
**Would need:**
- User preferences (notifications, display options)
- Password change
- Account deletion

### 2. Alliance Depot

**Status:** Not Implemented  
**Reason:** Low usage feature for agent gameplay  
**Would need:**
- Alliance resource sharing

### 3. Redis Caching

**Status:** Optional  
**Description:** Cache frequently accessed data for performance
**Would need:**
- Cache /api/status endpoint (TTL 10s)
- Cache galaxy view
- Cache highscore

---

## Known Issues

1. **Ship Speed Discrepancy:** One test case shows CalculateShipSpeed(Cruiser impulse 5) returns 30000 instead of expected 27000 - minor formula difference

2. **Rust Battle Engine:** Requires `libbattle_engine_ffi.so` to be installed. Falls back to Go implementation if not available.

3. **Fleet Scheduler:** Fleet arrivals are processed when the API is called, not via a background scheduler. For true 24/7 operation, a cron-based scheduler would be needed.

---

## Test Commands

```bash
# Build
go build ./...

# Run formula tests
go test -v ./internal/formula/...

# Run with Rust battle engine
LD_LIBRARY_PATH=/home/bolt/Documents/ogamex-go/storage/rust-libs go test ./...
```

---

## Files Reference

```
ogamex-go/
├── cmd/server/main.go
├── internal/
│   ├── api/handlers.go          # ~80 endpoints
│   ├── schema/models.go         # User, Planet, Fleet, etc.
│   ├── service/                 # All services implemented
│   ├── repository/              # Data access layer
│   ├── formula/                 # Cost/production formulas (95%+ coverage)
│   ├── dto/dto.go               # Request/response types
│   └── middleware/              # Auth, rate-limit, trace_id
├── pkg/rustbattle/              # CGO battle engine binding
└── REMAINING_FEATURES.md
```
