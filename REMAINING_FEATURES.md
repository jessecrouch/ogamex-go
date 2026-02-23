# OGameX-Go Remaining Features

**Last Updated:** Feb 2026  
**Purpose:** Document all features that still need implementation for full agent compatibility

---

## Overview

The OGameX-Go project has made significant progress (~90% complete), but there are several features needed for autonomous AI agents to play the full game. This document lists all remaining items in priority order.

---

## Priority 1: Agent Compatibility (Critical)

These features are required for agents to play the complete game loop.

### 1.1 Fleet Arrival Processing (Scheduler)

**Status:** Partial  
**Description:** Complete the scheduler to handle fleet arrivals, execute missions, and process battles.

**Required:**
- Process incoming fleets at arrival time
- Execute attack missions (battle simulation)
- Execute transport missions (resource transfer)
- Execute colonization missions
- Execute recycle missions (debris collection)
- Execute expedition missions
- Return fleets to origin

**Location:** `internal/scheduler/`

**Reference:** 
- OGameX `FleetService.php` - mission execution
- OGameX `QueueService.php` - timing logic

---

### 1.2 Moon Buildings - Jump Gate

**Status:** Not Started  
**Description:** Allow construction and use of Jump Gates on moons for fleet teleportation.

**Required:**
- Moon building: Jump Gate (building ID 43)
- Jump Gate activation/deactivation
- Fleet teleportation between moons
- Jump Gate cooldown timer (60 minutes)
- Jump Gate fuel consumption ( deuterium)

**Endpoints Needed:**
```
POST /planets/:id/jump-gate/activate
POST /planets/:id/jump-gate/target
POST /planets/:id/jump-gate/execute
```

---

### 1.3 Moon Buildings - Sensor Phalanx

**Status:** Not Started  
**Description:** Allow construction of Sensor Phalanx to detect incoming fleets.

**Required:**
- Moon building: Sensor Phalanx (building ID 42)
- Phalanx scan functionality
- Display incoming fleet missions

**Endpoints Needed:**
```
POST /planets/:id/phalanx/scan
```

---

## Priority 2: Multiplayer Features

### 2.1 Alliance System

**Status:** Not Started  
**Description:** Full alliance management system.

**Required:**
- Create alliance
- Apply to join alliance
- Accept/reject applications
- Alliance members list
- Alliance diplomacy ( NAP, war)
- Alliance chat/information

**Database:**
```sql
CREATE TABLE alliances (
    id PRIMARY KEY,
    name VARCHAR(50),
    tag VARCHAR(10),
    founder_id,
    description TEXT,
    created_at
);

CREATE TABLE alliance_members (
    alliance_id,
    user_id,
    rank_id,
    joined_at
);
```

**Endpoints:**
```
POST /alliances/create
GET /alliances/:id
POST /alliances/:id/apply
POST /alliances/:id/accept/:user_id
POST /alliances/:id/reject/:user_id
GET /alliances/:id/members
POST /alliances/:id/leave
```

---

### 2.2 Buddy System

**Status:** Not Started  
**Description:** Buddy request and management system.

**Required:**
- Send buddy request
- Accept/reject buddy requests
- List buddies

**Database:**
```sql
CREATE TABLE buddy_requests (
    id PRIMARY KEY,
    sender_id,
    receiver_id,
    message TEXT,
    status VARCHAR(20), -- pending, accepted, rejected
    created_at
);
```

**Endpoints:**
```
GET /buddy
POST /buddy/request
POST /buddy/:id/accept
POST /buddy/:id/reject
```

---

## Priority 3: Economy & Premium

### 3.1 Premium/Merchant Services

**Status:** Not Started  
**Description:** Dark Matter purchases and resource trading.

**Required:**
- Resource trading (exchange resources at merchant)
- Dark Matter packages
- Officer packages (admiral, engineer, geologist, etc.)

**Endpoints:**
```
POST /merchant/trade
GET /premium/packages
POST /premium/activate
```

---

### 3.2 Planet Move Service

**Status:** Not Started  
**Description:** Allow players to relocate their planet to a different position.

**Required:**
- Planet relocation (with dark matter cost)
- Position validation (not occupied)
- Fleet recall during move

**Endpoints:**
```
POST /planets/:id/relocate
```

---

## Priority 4: Galaxy Features

### 4.1 Debris Field Service

**Status:** Not Started  
**Description:** Manage debris fields from battles.

**Required:**
- Create debris on battle
- Position debris at coordinates
- Debris lifetime (12 hours in OGame)
- Debris recycling via fleet mission

**Database:**
```sql
CREATE TABLE debris_fields (
    id PRIMARY KEY,
    galaxy, system, position,
    metal INT64,
    crystal INT64,
    created_at,
    expires_at
);
```

---

### 4.2 Wreck Field Service

**Status:** Not Started  
**Description:** Handle wreck fields from large battles (OGame v8+).

**Required:**
- Create wreck fields on massive battles
- Wreck field collection
- Scrap ship technology bonus

---

## Priority 5: Additional Services

### 5.1 Highscore Enhancements

**Status:** Partial (basic implemented)  
**Description:** Full highscore with additional categories.

**Enhance:**
- Add research points calculation
- Add fleet points calculation
- Add total points (combined)
- Periodic ranking updates

---

### 5.2 NPC Fleet Generator

**Status:** Not Started  
**Description:** Generate NPC fleets for expeditions and missions.

**Required:**
- NPC expedition fleets
- NPC defense waves
- Pirate raids

---

### 5.3 NPC Planet Generation

**Status:** Not Started  
**Description:** Generate NPC planets for new players to attack.

**Required:**
- NPC planets at specific coordinates
- NPC defense (varies by level)
- NPC resources
- Loot calculation

---

## Priority 6: Optional Enhancements

### 6.1 Redis Caching

**Status:** Not Started  
**Description:** Cache frequently accessed data.

**Optional:**
- Cache /api/status endpoint (TTL 10s)
- Cache galaxy view (per user)
- Cache highscore (TTL 1 hour)

---

### 6.2 Rate Limiting Enhancements

**Status:** Partial  
**Description:** More granular rate limiting.

**Enhance:**
- Per-endpoint rate limits
- Fleet send rate limits
- Battle cooldown tracking

---

## Implementation Checklist

- [x] Fleet arrival processing (scheduler)
- [x] Jump Gate (moon teleport)
- [x] Sensor Phalanx (fleet detection)
- [x] Alliance system
- [x] Buddy system
- [x] Espionage reports
- [x] Vacation mode
- [x] Premium/Merchant
- [x] Planet move
- [x] Debris fields
- [x] Wreck fields
- [x] ACS (Fleet Union)
- [x] NPC fleets
- [x] NPC planets
- [ ] Redis caching (optional)

---

## Files Reference

**Schema (models to add):**
- `internal/schema/models.go` - Add Alliance, Buddy, Debris models

**Repositories:**
- `internal/repository/alliance_repo.go` - New
- `internal/repository/buddy_repo.go` - New
- `internal/repository/debris_repo.go` - New

**Services:**
- `internal/service/alliance_service.go` - New
- `internal/service/buddy_service.go` - New
- `internal/service/moon_service.go` - New (Jump Gate, Phalanx)
- `internal/service/debris_service.go` - New

**API Handlers:**
- Add to `internal/api/handlers.go`

---

## Notes

- Use formula package for all cost/time calculations
- Follow existing code patterns (services, repositories, interfaces)
- Add comprehensive tests for new features
- Ensure proper logging with trace_id
- Use context.Context for all database operations
