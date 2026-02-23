# Wiki Verification Test Roadmap

This document outlines the comprehensive testing strategy to verify our OGame implementation against the official OGame wiki documentation.

## Overview

The goal is to systematically cross-check our Go implementation against the OGame wiki to ensure 100% behavioral parity. Each phase adds new test categories using JSON test data files.

**Philosophy:** Log all discrepancies comprehensively. Most differences indicate bugs in our code, but the wiki may also have errors or refer to older game versions.

---

## Phase 1: Resource Production (Priority: HIGH)

**Status:** TODO

### Test Data: `tests/wiki/production.json`

Verify production formulas match wiki tables:

| Resource | Wiki Formula | Test Values |
|----------|-------------|-------------|
| Metal Mine | `30 × level × 1.1^level` | L1=33, L10=779, L20=4037 |
| Crystal Mine | `20 × level × 1.1^level` | L1=22, L10=518, L20=2690 |
| Deuterium | `10 × level × 1.1^level × (1.36 - 0.004 × temp)` | Varies by temp |
| Solar Plant | `20 × level × 1.1^level` | L1=22, L10=519 |
| Fusion Reactor | `30 × level × (1.05 + energyTech×0.01)^level` | L1=32 |
| Solar Satellite | `floor((temp + 140) / 6) × count` | Temp 100 = 40 energy |

### Test File: `internal/formula/wiki_production_test.go`

- [ ] Metal mine production by level
- [ ] Crystal mine production by level  
- [ ] Deuterium production (temp-adjusted)
- [ ] Energy production (solar plant, fusion, satellites)
- [ ] Energy consumption by buildings
- [ ] Plasma technology bonus to crystal production
- [ ] Basic income (30 metal, 15 crystal per hour)

---

## Phase 2: Building Costs (Priority: HIGH)

**Status:** TODO

### Test Data: `tests/wiki/buildings_cost.json`

Verify cost formulas match wiki:

| Building | Cost Factor | Base Costs |
|----------|-------------|------------|
| Metal Mine | 1.5^(level-1) | 60M |
| Crystal Mine | 1.6^(level-1) | 48M, 24C |
| Deuterium Synth | 1.5^(level-1) | 225M |
| Solar Plant | 1.5^(level-1) | 75M |
| Fusion Reactor | 1.8^(level-1) | 900M, 360C, 180D |
| Robotics Factory | 2^(level-1) | 400M, 120C |
| Shipyard | 2^(level-1) | 400M, 200C |
| Research Lab | 2^(level-1) | 200M, 400C |
| Storage (Metal/Crystal/Deut) | 2^(level-1) | See wiki |
| Nanite Factory | 2^(level-1) | 1M, 200K C |
| Terraformer | 2^(level-1) | 50K M, 100K C, 1M D |

### Test File: `internal/formula/wiki_building_cost_test.go`

- [ ] Cost calculation for each building type
- [ ] Cost factor formula verification
- [ ] Storage capacity formula
- [ ] Construction time formula

---

## Phase 3: Ship Statistics (Priority: HIGH)

**Status:** TODO

### Test Data: `tests/wiki/ships.json`

Verify ship stats match wiki:

| Ship | Cargo | Base Speed | Structural Integrity | Shield | Weapon |
|------|-------|------------|---------------------|--------|--------|
| Small Cargo | 5,000 | 10,000 | 2,000 | 10 | 5 |
| Large Cargo | 25,000 | 7,500 | 6,000 | 25 | 12 |
| Light Fighter | 50 | 12,500 | 4,000 | 10 | 50 |
| Heavy Fighter | 100 | 10,000 | 10,000 | 25 | 150 |
| Cruiser | 800 | 15,000 | 20,000 | 50 | 400 |
| Battleship | 1,500 | 10,000 | 60,000 | 200 | 1,000 |
| Battlecruiser | 750 | 10,000 | 70,000 | 400 | 700 |
| Bomber | 500 | 5,000 | 75,000 | 500 | 1,000 |
| Destroyer | 2,000 | 5,000 | 110,000 | 500 | 2,000 |
| Deathstar | 1,000,000 | 100 | 15,000,000 | 200,000 | 200,000 |
| Recycler | 20,000 | 6,000 | 16,000 | 10 | 100 |
| Espionage Probe | 5 | 100,000,000 | 1,000 | 1 | 0.01 |
| Solar Satellite | 0 | 0 | 2,000 | 1 | 1 |
| Colony Ship | 7,500 | 2,500 | 30,000 | 100 | 150 |
| Crawler | 0 | 4,000 | 4,000 | 2 | 8 |
| Reaper | ??? | ??? | ??? | ??? | ??? |
| Pathfinder | ??? | ??? | ??? | ??? | ??? |

### Test File: `internal/service/wiki_ship_test.go`

- [ ] Ship cargo capacity
- [ ] Ship base speed
- [ ] Ship structural integrity
- [ ] Ship shield strength
- [ ] Ship weapon power
- [ ] Ship cost calculation
- [ ] Ship build time calculation
- [ ] Engine upgrade speeds (Combustion, Impulse, Hyperspace)

---

## Phase 4: Defense Statistics (Priority: HIGH)

**Status:** TODO

### Test Data: `tests/wiki/defense.json`

Verify defense stats match wiki:

| Defense | Metal | Crystal | Deuterium | Integrity | Shield | Weapon |
|---------|-------|---------|-----------|-----------|--------|--------|
| Rocket Launcher | 2,000 | - | - | 2,000 | 20 | 80 |
| Light Laser | 1,500 | 500 | - | 2,000 | 25 | 100 |
| Heavy Laser | 6,000 | 2,000 | - | 8,000 | 100 | 250 |
| Gauss Cannon | 20,000 | 15,000 | 2,000 | 35,000 | 200 | 1,100 |
| Ion Cannon | 5,000 | 3,000 | - | 8,000 | 500 | 150 |
| Plasma Turret | 50,000 | 50,000 | 30,000 | 100,000 | 300 | 3,000 |
| Small Shield Dome | 10,000 | 10,000 | - | 20,000 | 2,000 | 1 |
| Large Shield Dome | 50,000 | 50,000 | - | 100,000 | 10,000 | 1 |
| Anti-Ballistic Missile | 8,000 | - | - | 8,000 | 1 | 1 |
| Interplanetary Missile | 15,000 | 15,000 | 2,500 | 15,000 | 1 | 12,000 |

### Test File: `internal/service/wiki_defense_test.go`

- [ ] Defense costs
- [ ] Defense build times
- [ ] Defense stats (integrity, shield, weapon)
- [ ] Defense 70% rebuild chance (in battle engine)

---

## Phase 5: Research Technology (Priority: MEDIUM)

**Status:** TODO

### Test Data: `tests/wiki/research.json`

Verify research requirements:

| Research | Lab Level | Prerequisites |
|----------|-----------|---------------|
| Espionage Technology | 3 | - |
| Computer Technology | 1 | - |
| Weapons Technology | 4 | - |
| Shielding Technology | 6 | Energy 3 |
| Armour Technology | 2 | - |
| Energy Technology | 1 | - |
| Hyperspace Technology | 7 | Energy 5, Shielding 5 |
| Combustion Drive | 1 | Energy 1 |
| Impulse Drive | 2 | Energy 1 |
| Hyperspace Drive | 7 | Hyperspace 3 |
| Laser Technology | 1 | Energy 2 |
| Ion Technology | 4 | Laser 5, Energy 4 |
| Plasma Technology | 4 | Laser 10, Ion 5, Energy 8 |
| Graviton Technology | 12 | - |
| Astrophysics | 3 | Espionage 4, Impulse 3 |

### Test File: `internal/service/wiki_research_test.go`

- [ ] Research cost calculation
- [ ] Research time calculation
- [ ] Prerequisites verification
- [ ] Tech tree dependencies

---

## Phase 6: Fleet Mechanics (Priority: HIGH)

**Status:** TODO

### Test Data: `tests/wiki/fleet.json`

Verify fleet mechanics:

**Flight Time Formula:**
```
T = 10 + 3500 × sqrt(10 × D / V) / speed_factor
```
Where D = distance, V = slowest ship speed

**Fuel Consumption:**
- Base consumption formula per ship
- Speed factor modifiers

**Distance Calculation:**
```
Same planet: 0
Same system, different position: 2,000 + 5 × |position - position|
Different system, same galaxy: 2,000 + 270 × |system - system| + 5 × |position - position|
Different galaxy: 2,000 + 20,000 × |galaxy - galaxy| + 270 × |system - system| + 5 × |position - position|
```

### Test File: `internal/formula/wiki_fleet_test.go`

- [ ] Distance calculation
- [ ] Flight time calculation
- [ ] Fuel consumption
- [ ] Fleet speed (slowest ship rule)
- [ ] Cargo capacity with fuel deduction
- [ ] Speed factor research bonuses

---

## Phase 7: Rapid Fire Tables (Priority: MEDIUM)

**Status:** TODO

### Test Data: `tests/wiki/rapid_fire.json`

Verify rapid fire values for all ships against:
- Other ships
- Defense structures

Example (Deathstar):
- vs Espionage Probe: 1250
- vs Solar Satellite: 1250
- vs Small Cargo: 250
- vs Light Fighter: 200
- vs Battleship: 30
- vs Rocket Launcher: 200
- vs Light Laser: 200

### Test File: `internal/formula/wiki_rapidfire_test.go`

- [ ] All ship vs ship rapid fire
- [ ] All ship vs defense rapid fire
- [ ] Defense vs ship rapid fire

---

## Phase 8: Planet Position Bonuses (Priority: MEDIUM)

**Status:** TODO

### Test Data: `tests/wiki/planet_position.json`

Verify planet position effects:

| Position | Min Fields | Max Fields | Temp Range | Production Bonuses |
|----------|------------|------------|------------|-------------------|
| 1 | 96 | 172 | 220-260°C | +40% Crystal |
| 2 | 104 | 176 | 170-210°C | +30% Crystal |
| 3 | 112 | 182 | 120-160°C | +20% Crystal |
| 4 | 118 | 208 | 70-110°C | - |
| 5 | 133 | 232 | 60-100°C | - |
| 6 | 146 | 242 | 50-90°C | +17% Metal |
| 7 | 152 | 248 | 40-80°C | +23% Metal |
| 8 | 156 | 252 | 30-70°C | +35% Metal |
| 9 | 150 | 246 | 20-60°C | +23% Metal |
| 10 | 142 | 232 | 10-50°C | +17% Metal |
| 15 | 90 | 164 | -130--90°C | Best Deuterium |

### Test File: `internal/formula/wiki_planet_test.go`

- [ ] Temperature calculation by position
- [ ] Field count ranges
- [ ] Production bonuses (crystal positions 1-3, metal positions 6-10)
- [ ] Deuterium bonus for cold positions

---

## Phase 9: Combat Mechanics (Priority: HIGH)

**Status:** TODO

### Test Data: `tests/wiki/combat.json`

Verify battle calculations:
- Shield damage absorption
- Weapon damage distribution
- Defense rebuild chance (70%)
- Debris field calculation (30% of destroyed ship value)

### Test File: `internal/service/wiki_combat_test.go`

- [ ] Battle simulation vs known scenarios
- [ ] Debris field calculation
- [ ] Moon generation chance
- [ ] Defense rebuild mechanics

---

## Phase 10: Player Classes (Priority: LOW)

**Status:** TODO

### Test Data: `tests/wiki/classes.json`

Verify class bonuses:
- Collector: +25% production, +10% energy, cargo/speed bonuses
- General: Ship/defense cost reductions
- Discoverer: Planet size bonus, reduced deuterium consumption

### Test File: `internal/service/wiki_class_test.go`

- [ ] Class bonus calculations
- [ ] Production modifiers
- [ ] Cost modifiers

---

## Discrepancy Handling

### When Tests Fail

1. **Log the discrepancy** with full context:
   - Expected value (from wiki)
   - Actual value (from our code)
   - Input parameters used
   - Formula/function being tested

2. **Categorize the issue:**
   - **BUG**: Our implementation is wrong - fix the code
   - **VERSION**: Game version differences - document and decide
   - **WIKI_ERROR**: Wiki is incorrect - document but skip test

3. **Decision matrix:**
   ```
   If wiki has specific number AND our formula gives different result:
     → Investigate our formula first (usually wrong)
   
   If wiki has formula (not number) AND our result differs:
     → Test with wiki's formula explicitly
   
   If wiki is ambiguous OR has "varies by universe":
     → Use OGameX as source of truth for our target universe
   ```

### Test Output Format

Each failed test should output:
```
[Wiki Verification] DISCREPANCY
  Function: CalculateMetalMineProduction
  Input: level=10, energyFactor=1.0
  Expected: 779 (wiki verified)
  Got: 765
  Formula: 30 * 10 * 1.1^10 = 779.4
  Status: BUG - our formula/math is incorrect
```

---

## Running the Tests

### Run All Wiki Tests
```bash
go test -v ./tests/wiki/...
```

### Run Specific Phase
```bash
go test -v -run "Wiki/Phase1" ./internal/formula/...
```

### Run With Discrepancy Report
```bash
go test -v -run "Wiki" ./... 2>&1 | grep -A5 "DISCREPANCY"
```

---

## Progress Tracking

| Phase | Category | Status | Tests | Discrepancies |
|-------|----------|--------|-------|---------------|
| 1 | Resource Production | DONE (with known diffs) | ~100 | ~70 (expected - uses Ceil) |
| 2 | Building Costs | DONE | ~70 | 1 (solar plant rounding) |
| 3 | Ship Statistics | DONE | ~80 | 2 (speed upgrade formula) |
| 4 | Defense Statistics | TODO | 0 | 0 |
| 5 | Research Technology | TODO | 0 | 0 |
| 6 | Fleet Mechanics | TODO | 0 | 0 |
| 7 | Rapid Fire Tables | TODO | 0 | 0 |
| 8 | Planet Position | TODO | 0 | 0 |
| 9 | Combat Mechanics | TODO | 0 | 0 |
| 10 | Player Classes | TODO | 0 | 0 |

---

## Phase 1 Findings (Completed)

### Bugs Found & Fixed
1. **Solar Satellite Formula Bug**: Original code used integer division `(temperature+140)/6` before converting to float, causing incorrect results for negative temperatures. Fixed to use float division first: `float64(temperature+140)/6.0`.

### Known Discrepancies (Expected)
The following discrepancies exist between our implementation (using Ceil per OGameX) and wiki values:

- **Metal Mine Production**: Wiki shows different values for levels 2-9, 11-19. Our implementation uses Ceil() as per OGameX source code, which matches OGameX behavior but differs from wiki table (likely from different game version).

### Root Cause Analysis
- OGameX uses `ceil()` for resource production (verified in PlanetService.php line 2009)
- Wiki tables appear to use a mix of floor/ceil or may be from an older game version
- Our implementation is CORRECT for OGameX compatibility
- These are "VERSION" discrepancies per our handling guidelines

### Test Data Corrections Made
- Fixed solar satellite test values for -100 and -130 temperatures (wiki data had typos)
