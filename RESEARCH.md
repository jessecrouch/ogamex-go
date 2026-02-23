# OGameX Research - Porting Notes

## Source: https://github.com/lanedirt/OGameX

## Location: /tmp/OGameX/

---

## Architecture Overview

### Tech Stack (Original)
- **Backend**: Laravel (PHP 8.x)
- **Database**: MySQL/PostgreSQL
- **Battle Engine**: Rust (via cgo FFI) - 200x faster than PHP
- **Queue**: Redis + Laravel Horizon

### Tech Stack (Target - per AGENTS.md)
- **Framework**: Fiber v2 (Go)
- **ORM**: GORM v2 + PostgreSQL
- **Battle Engine**: cgo → Rust `.so` (unchanged)
- **Scheduler**: cron + in-memory queue + Redis option

---

## Rust Battle Engine (DO NOT REIMPLEMENT - Use cgo)

**Location**: `rust/battle_engine_ffi/`

### FFI Interface
Single function signature:
```c
pub extern "C" fn fight_battle_rounds(input_json: *const c_char) -> *mut c_char
```

**Input**: JSON with attacker_fleets and defender_fleets  
**Output**: JSON battle results with rounds, losses, debris

### Compilation
```bash
cd rust
cargo build --release
cp target/release/libbattle_engine_ffi.so storage/rust-libs/
```

### Performance
- ~175ms vs ~54,000ms (PHP) for large fleets
- Memory: ~55MB vs ~226MB (PHP)

---

## Database Schema (Key Tables)

### Migrations Location: `database/migrations/`

**Core tables:**
- `users` - Player accounts
- `planets` - Planet data with resources
- `users_tech` - Research levels
- `building_queue` - Construction queue
- `research_queue` - Research queue
- `unit_queue` - Ship/defense construction
- `fleet_missions` - Active fleet missions
- `debris_fields` - Debris fields
- `espionage_reports` - Spy reports
- `battle_reports` - Combat reports
- `messages` - Player messages
- `alliances` - Alliance data

### Key Planet Fields
```php
metal, crystal, deuterium
metal_production, crystal_production, deuterium_production
energy_available, energy_max, energy_used
metal_mine, crystal_mine, deuterium_synthesizer
solar_plant, fusion_plant, solar_satellite
metal_mine_percent, crystal_mine_percent, deuterium_synthesizer_percent
```

---

## Services to Port (32 files)

**Location**: `app/Services/`

| Service | Purpose |
|---------|---------|
| BuildingQueueService | Building construction queue |
| ResearchQueueService | Research queue |
| UnitQueueService | Ship/defense queue |
| FleetMissionService | Fleet missions (attack, transport, etc.) |
| FleetUnionService | ACS (Alliance Combat System) |
| PlanetService | Resource production, planet management |
| PlayerService | Player management |
| DarkMatterService | Premium currency |
| DebrisFieldService | Debris collection |
| WreckFieldService | Wreck field processing |
| CoordinateDistanceCalculator | Flight time calculations |
| ObjectService | Game object definitions |
| CharacterClassService | Collector, General, Discoverer bonuses |
| HalvingService | Dark matter cost halving |
| JumpGateService | Moon teleport gates |
| PhalanxService | Moon defense scanner |
| MerchantService | Resource trading |
| AllianceService | Alliance management |
| MessageService | Player messages |
| HighscoreService | Rankings |

---

## Key Formulas (from tests)

### Resource Production
- **Metal Mine**: `20 * level * 1.1^level * energy_factor`
- **Crystal Mine**: `20 * level * 1.1^level * energy_factor`
- **Deuterium Synthesizer**: `20 * level * 1.1^level * energy_factor`
- **Solar Plant**: `20 * level * 1.1^level`
- **Fusion Plant**: `30 * level * 1.1^level` (consumes deuterium)
- **Solar Satellite**: `TempFactor * 0.1 * solar_sat_count`

### Energy Balance
- Positive energy = full production
- Zero/negative energy = production reduced to 0-50%
- Fusion plant requires deuterium in storage to produce energy

### Position Bonuses (Planet Position 1-15)
```
1-3:  100% base
4:    110% crystal
5-7:  100% base
8:    120% crystal
9:    130% crystal
10:   140% crystal
11:   160% crystal
12:   170% crystal
13-15: 100% base
```

### Fleet Speed
- Base speeds per ship type
- Speed bonus from research
- Universe speed multiplier

### Distance Calculation
```
distance = |x1 - x2| * 20000 + |y1 - y2| * 5 + |z1 - z2| * 20000
```

### Flight Time
```
time = 35000 * distance / (speed * 10)
```

---

## Unit Tests to Replicate

**Location**: `tests/Unit/`

| Test File | Coverage |
|-----------|----------|
| ResourceProductionTest.php | Mine production, energy, position bonuses |
| BuildingQueueServiceTest.php | Queue logic, costs |
| HalvingServiceTest.php | Dark matter halving |
| CoordinateDistanceCalculatorTest.php | Distance & flight time |
| RustFfiTest.php | Battle engine integration |
| FleetCheckTest.php | Fleet validation |
| UnitCollectionTest.php | Ship collections |

---

## Game Constants

**Location**: `app/GameConstants/`

Need to extract:
- Building costs (metal, crystal, deuterium)
- Building production rates
- Unit stats (attack, shield, cargo, speed)
- Unit costs
- Research costs and effects

---

## Edge Cases to Handle

1. **Energy deficit** - Reduce production when energy < 0
2. **Fusion plant** - Requires deuterium to produce energy
3. **Moon destruction** - 20% chance when > 5000 units destroyed
4. **Expedition** - Random rewards (none, resources, ships, nothing)
5. **ACS** - Multiple attackers, combined fleet combat
6. **Debris/Wreck** - 50% of destroyed resources become debris
7. **Dark matter halving** - Prices halved based on player DM purchases

---

## API Endpoints (for reference)

**Location**: `routes/`

Key endpoints to implement:
- `/api/v1/buildings/upgrade`
- `/api/v1/fleets/send`
- `/api/v1/planets/resources`
- `/api/v1/research/start`
- `/api/v1/shipyard/build`
- `/api/v1/status` (for agent polling)

---

## Notes

- All production and costs use **integer math** (int64 in Go)
- Use `ceil()` for positive production, floor not specified
- Server uses UTC timestamps
- Race conditions handled via database transactions
- Fleets arrive based on server time, not client time
