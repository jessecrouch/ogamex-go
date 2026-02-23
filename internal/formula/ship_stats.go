package formula

import (
	"math"

	"ogamex-go/internal/domain"
)

type ShipStats struct {
	CargoCapacity       int64
	BaseSpeed           int64
	StructuralIntegrity int64
	Shield              int64
	Weapon              int64
	EngineType          string
}

var shipStats = map[domain.UnitType]ShipStats{
	domain.UnitSmallCargo:     {5000, 5000, 2000, 10, 5, "combustion"},
	domain.UnitLargeCargo:     {25000, 7500, 6000, 25, 12, "combustion"},
	domain.UnitLightFighter:  {50, 12500, 4000, 10, 50, "combustion"},
	domain.UnitHeavyFighter:   {100, 10000, 10000, 25, 150, "impulse"},
	domain.UnitCruiser:        {800, 15000, 27000, 50, 400, "impulse"},
	domain.UnitBattleship:     {1500, 10000, 60000, 200, 1000, "hyperspace"},
	domain.UnitBattlecruiser:  {750, 10000, 70000, 400, 700, "hyperspace"},
	domain.UnitBomber:         {500, 4000, 75000, 500, 1000, "impulse"},
	domain.UnitDestroyer:      {2000, 5000, 110000, 500, 2000, "hyperspace"},
	domain.UnitDeathstar:      {1000000, 100, 9000000, 50000, 200000, "hyperspace"},
	domain.UnitRecycler:       {20000, 6000, 16000, 10, 100, "hyperspace"},
	domain.UnitEspionageProbe: {5, 100000000, 1000, 1, 0, "combustion"},
	domain.UnitSolarSatellite: {0, 0, 2000, 1, 1, "none"},
	domain.UnitCrawler:        {0, 4000, 4000, 2, 8, "combustion"},
	domain.UnitColonyShip:     {7500, 2500, 30000, 100, 150, "impulse"},
	domain.UnitReaper:         {700, 10000, 140000, 700, 2800, "hyperspace"},
	domain.UnitPathfinder:     {500, 12000, 23000, 100, 200, "hyperspace"},
}

func GetShipStats(unitType domain.UnitType) (ShipStats, bool) {
	stats, exists := shipStats[unitType]
	return stats, exists
}

func CalculateShipCargoCapacity(unitType domain.UnitType) int64 {
	stats, exists := shipStats[unitType]
	if !exists {
		return 0
	}
	return stats.CargoCapacity
}

func CalculateShipBaseSpeed(unitType domain.UnitType) int64 {
	stats, exists := shipStats[unitType]
	if !exists {
		return 0
	}
	return stats.BaseSpeed
}

func CalculateShipSpeed(unitType domain.UnitType, combustionLevel, impulseLevel, hyperspaceLevel int) int64 {
	stats, exists := shipStats[unitType]
	if !exists || stats.BaseSpeed == 0 {
		return stats.BaseSpeed
	}

	var bonus float64 = 1.0
	switch stats.EngineType {
	case "combustion":
		if unitType == domain.UnitEspionageProbe || unitType == domain.UnitRecycler {
			bonus = 1.0 + float64(combustionLevel)*0.1
		} else {
			bonus = 1.0 + float64(combustionLevel)*0.02
		}
	case "impulse":
		bonus = 1.0 + float64(impulseLevel)*0.16
	case "hyperspace":
		bonus = 1.0 + float64(hyperspaceLevel)*0.3
	}

	speed := float64(stats.BaseSpeed) * bonus
	return int64(math.Floor(speed))
}

func CalculateShipStructuralIntegrity(unitType domain.UnitType) int64 {
	stats, exists := shipStats[unitType]
	if !exists {
		return 0
	}
	return stats.StructuralIntegrity
}

func CalculateShipShield(unitType domain.UnitType) int64 {
	stats, exists := shipStats[unitType]
	if !exists {
		return 0
	}
	return stats.Shield
}

func CalculateShipWeapon(unitType domain.UnitType) int64 {
	stats, exists := shipStats[unitType]
	if !exists {
		return 0
	}
	return stats.Weapon
}

type DefenseStats struct {
	StructuralIntegrity int64
	Shield             int64
	Weapon             int64
}

var defenseStats = map[domain.UnitType]DefenseStats{
	domain.UnitRocketLauncher:          {2000, 20, 80},
	domain.UnitLightLaser:             {2000, 25, 100},
	domain.UnitHeavyLaser:             {8000, 100, 250},
	domain.UnitGaussCannon:           {35000, 200, 1100},
	domain.UnitIonCannon:             {8000, 500, 150},
	domain.UnitPlasmaTurret:           {100000, 300, 3000},
	domain.UnitSmallShieldDome:       {20000, 2000, 1},
	domain.UnitLargeShieldDome:       {100000, 10000, 1},
	domain.UnitAntiBallisticMissiles: {8000, 1, 1},
	domain.UnitInterplanetaryMissiles: {15000, 1, 12000},
}

func GetDefenseStats(unitType domain.UnitType) (DefenseStats, bool) {
	stats, exists := defenseStats[unitType]
	return stats, exists
}

func CalculateDefenseStructuralIntegrity(unitType domain.UnitType) int64 {
	stats, exists := defenseStats[unitType]
	if !exists {
		return 0
	}
	return stats.StructuralIntegrity
}

func CalculateDefenseShield(unitType domain.UnitType) int64 {
	stats, exists := defenseStats[unitType]
	if !exists {
		return 0
	}
	return stats.Shield
}

func CalculateDefenseWeapon(unitType domain.UnitType) int64 {
	stats, exists := defenseStats[unitType]
	if !exists {
		return 0
	}
	return stats.Weapon
}

var rapidFireAgainst = map[domain.UnitType]map[domain.UnitType]int{
	domain.UnitDeathstar: {
		domain.UnitEspionageProbe:   1250,
		domain.UnitSolarSatellite:   1250,
		domain.UnitSmallCargo:       250,
		domain.UnitLargeCargo:       250,
		domain.UnitLightFighter:    200,
		domain.UnitHeavyFighter:     100,
		domain.UnitCruiser:          33,
		domain.UnitBattleship:       30,
		domain.UnitColonyShip:       250,
		domain.UnitRecycler:         250,
		domain.UnitBomber:           25,
		domain.UnitDestroyer:        5,
		domain.UnitBattlecruiser:    15,
		domain.UnitReaper:           30,
		domain.UnitPathfinder:       10,
		domain.UnitCrawler:          1250,
		domain.UnitRocketLauncher:  200,
		domain.UnitLightLaser:       200,
		domain.UnitHeavyLaser:      100,
		domain.UnitIonCannon:        100,
		domain.UnitGaussCannon:      50,
	},
	domain.UnitSolarSatellite: {
		domain.UnitEspionageProbe: 5,
		domain.UnitSmallCargo:     3,
		domain.UnitLargeCargo:     3,
		domain.UnitCrawler:         5,
	},
	domain.UnitSmallCargo: {
		domain.UnitEspionageProbe: 5,
		domain.UnitSolarSatellite: 5,
		domain.UnitSmallCargo:     3,
		domain.UnitLargeCargo:     3,
		domain.UnitCrawler:        5,
	},
	domain.UnitLightFighter: {
		domain.UnitEspionageProbe: 5,
		domain.UnitSolarSatellite: 5,
		domain.UnitSmallCargo:     3,
		domain.UnitLargeCargo:     3,
		domain.UnitCrawler:        5,
	},
	domain.UnitHeavyFighter: {
		domain.UnitEspionageProbe: 5,
		domain.UnitSolarSatellite: 5,
		domain.UnitSmallCargo:     4,
		domain.UnitLargeCargo:     4,
		domain.UnitCrawler:        5,
	},
	domain.UnitCruiser: {
		domain.UnitEspionageProbe: 5,
		domain.UnitSolarSatellite: 5,
		domain.UnitSmallCargo:     3,
		domain.UnitLargeCargo:     3,
		domain.UnitLightFighter:   6,
		domain.UnitHeavyFighter:   4,
		domain.UnitCrawler:        5,
	},
	domain.UnitBattleship: {
		domain.UnitEspionageProbe: 5,
		domain.UnitSolarSatellite: 5,
		domain.UnitSmallCargo:     3,
		domain.UnitLargeCargo:     3,
		domain.UnitCrawler:        5,
	},
	domain.UnitBattlecruiser: {
		domain.UnitEspionageProbe: 5,
		domain.UnitSolarSatellite: 5,
		domain.UnitSmallCargo:     3,
		domain.UnitLargeCargo:     3,
		domain.UnitLightFighter:   4,
		domain.UnitHeavyFighter:   7,
		domain.UnitCrawler:        5,
	},
	domain.UnitBomber: {
		domain.UnitEspionageProbe:   5,
		domain.UnitSolarSatellite:   5,
		domain.UnitSmallCargo:       3,
		domain.UnitLargeCargo:       3,
		domain.UnitCrawler:          5,
		domain.UnitRocketLauncher:   20,
		domain.UnitLightLaser:       20,
		domain.UnitHeavyLaser:       10,
		domain.UnitIonCannon:        10,
	},
	domain.UnitDestroyer: {
		domain.UnitEspionageProbe: 5,
		domain.UnitSolarSatellite: 5,
		domain.UnitSmallCargo:     3,
		domain.UnitLargeCargo:     3,
		domain.UnitCrawler:        5,
		domain.UnitDeathstar:      2,
	},
	domain.UnitReaper: {
		domain.UnitEspionageProbe:  5,
		domain.UnitSolarSatellite:  5,
		domain.UnitSmallCargo:      5,
		domain.UnitLargeCargo:      5,
		domain.UnitLightFighter:    5,
		domain.UnitHeavyFighter:    5,
		domain.UnitCruiser:         5,
		domain.UnitBattleship:       5,
		domain.UnitColonyShip:      5,
		domain.UnitRecycler:        5,
		domain.UnitBomber:          5,
		domain.UnitDestroyer:       5,
		domain.UnitDeathstar:       1250,
		domain.UnitReaper:          5,
		domain.UnitPathfinder:      5,
	},
	domain.UnitPathfinder: {
		domain.UnitEspionageProbe:  5,
		domain.UnitSolarSatellite:  5,
		domain.UnitSmallCargo:      5,
		domain.UnitLargeCargo:      5,
		domain.UnitLightFighter:    5,
		domain.UnitHeavyFighter:    5,
		domain.UnitCruiser:         5,
		domain.UnitBattleship:       5,
		domain.UnitColonyShip:      5,
		domain.UnitRecycler:        5,
		domain.UnitBomber:          5,
		domain.UnitDestroyer:       5,
		domain.UnitDeathstar:       1250,
		domain.UnitReaper:          5,
		domain.UnitPathfinder:       5,
	},
}

func GetRapidFireAgainst(attacker domain.UnitType, defender domain.UnitType) int {
	if attacks, ok := rapidFireAgainst[attacker]; ok {
		if rapidFire, ok := attacks[defender]; ok {
			return rapidFire
		}
	}
	return 1
}
