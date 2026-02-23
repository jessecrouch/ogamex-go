# OGameX GitHub Issues - Bug Analysis for OGameX-Go

**Last Updated:** Feb 2026  
**Purpose:** Analyze OGameX bugs and verify whether they exist in OGameX-Go

---

## Summary

| Issue | Title | Status | Our Status | Fix? |
|-------|-------|--------|------------|------|
| #1208 | Espionage Probes in wreckage | OPEN | **BUG EXISTS** | Yes - need to exclude ship 210 |
| #1188 | Espionage report missing player class | OPEN | **BUG EXISTS** | Yes - add character class to report |
| #1207 | DM halving queue calculation | CLOSED | DIFFERENT | No - simpler implementation works |
| #1213 | Deathstar rapid fire missing | OPEN | OK | No - we have all ships covered |
| #1214 | Player defenses in NPC fleet | CLOSED | OK | No - we don't include defenses |
| #1203 | Admins can be attacked | OPEN | NEEDS VERIFY | - |
| #1186 | Crawler on moons | OPEN | NEEDS VERIFY | - |
| #1216 | Fleet tooltip "Own Fleet" | OPEN | N/A | Frontend only |

---

## Detailed Analysis

### 🔴 Bugs That Exist in OGameX-Go

#### 1. Issue #1208 - Espionage Probes Should Not Be Part of Wreckage

**OGameX Status:** OPEN  
**Description:** When wreckage is created, espionage probes are incorrectly included in the calculations. Since espionage probes cannot be recovered via Space Dock, they shouldn't be in wreckage.

**Our Implementation:** `internal/service/debris_service.go:48-54`

```go
costs := map[int16]int64{
    202: 4000, 203: 12000, 204: 18000, 205: 45000, 206: 40000,
    207: 90000, 208: 20000, 209: 18000, 210: 1000,  // 210 = EspionageProbe!
    // ...
}
```

**Verdict:** BUG EXISTS - Ship ID 210 (EspionageProbe) is included in debris calculation.

**Fix Required:** Remove 210 from the costs map in `CreateDebrisFromBattle()`.

---

#### 2. Issue #1188 - Espionage Report Not Passing Player Class

**OGameX Status:** OPEN  
**Description:** When spying on a player, the player class isn't displayed in the espionage report - it defaults to "Class: Unknown".

**Our Implementation:** `internal/service/espionage_service.go:77-82`

```go
user, _ := s.userRepo.GetByID(ctx, targetPlanet.UserID)
username := "Unknown"
if user != nil {
    username = user.Username
    // CharacterClass is NOT being fetched or included!
}
```

**Verdict:** BUG EXISTS - CharacterClass is not included in espionage reports.

**Fix Required:** Add CharacterClass field to the espionage report structure and populate it.

---

### 🟡 Different Implementation (Not a Bug)

#### Issue #1207/#1206 - Dark Matter Halving Queue

**OGameX Status:** CLOSED (Fixed)  
**Description:** When having multiple items queued and using DM halving, the remaining queue time was not properly adjusted.

**Our Implementation:** `internal/service/halving_service.go`

Our implementation is simpler - it just applies 50% time reduction but doesn't recalculate the full queue. This works for agents but doesn't match the exact OGame behavior.

**Verdict:** DIFFERENT IMPLEMENTATION - Not a critical bug for API usage. The halving still works (50% time reduction), just doesn't recalculate queue.

**Fix Optional:** Would require more complex queue management logic.

---

### 🟢 Not Applicable / Already Fixed

#### Issue #1213 - Deathstar Missing Rapid Fire

**OGameX Status:** OPEN  
**Description:** Deathstar doesn't have rapid fire against new ships from player classes.

**Our Implementation:** `internal/formula/ship_stats.go:156-178`

Our Deathstar rapidfire includes: EspionageProbe, SolarSatellite, SmallCargo, LargeCargo, LightFighter, HeavyFighter, Cruiser, Battleship, ColonyShip, Recycler, Bomber, Destroyer, Battlecruiser, Reaper, Pathfinder, Crawler, RocketLauncher, LightLaser, HeavyLaser, IonCannon, GaussCannon.

**Verdict:** OK - All ships we support are covered.

---

#### Issue #1214 - Player Defenses in NPC Fleet (CLOSED)

**OGameX Status:** CLOSED  
**Description:** Player defenses were incorrectly added to NPC fleet on expeditions.

**Our Implementation:** `internal/service/npc_service.go:186-215`

We only pass ship IDs (210, 202-219) to NPC fleet creation, no defense units.

**Verdict:** OK - We don't include defenses in NPC fleets.

---

#### Issue #1216 - Fleet Events Tooltip

**OGameX Status:** OPEN  
**Description:** Fleet events always shown as "Own Fleet" in tooltip.

**Verdict:** N/A - This is a frontend/UI bug. Our API returns raw data, not tooltips.

---

### ❓ Need Verification

#### Issue #1203 - Admins Can Be Attacked

**OGameX Status:** OPEN  
**Description:** Admin planets can be attacked via galaxy view or manual coordinates.

**Need to verify:** Check if `fleet_service.go` validates target planet isn't an admin (IsNPC or similar flag).

---

#### Issue #1186 - Crawler on Moons

**OGameX Status:** OPEN  
**Description:** Crawlers should not be buildable on moons but appear in shipyard.

**Need to verify:** Check if unit service filters out Crawlers when planet is moon.

---

## Fixes to Implement

### Priority 1 (High)

1. **Fix #1208:** Remove EspionageProbe (210) from debris calculation costs

2. **Fix #1188:** Add CharacterClass to espionage reports

### Priority 2 (Medium)

3. **Verify #1203:** Ensure admin planets can't be attacked

4. **Verify #1186:** Ensure Crawlers filtered from moon shipyard

---

## Test Commands

```bash
# Build
go build ./...

# Run formula tests
go test -v ./internal/formula/...
```
