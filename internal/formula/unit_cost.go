package formula

import (
	"ogamex-go/internal/domain"
)

type UnitBaseCost struct {
	Metal     int64
	Crystal   int64
	Deuterium int64
	BuildTime int64
}

var shipBaseCosts = map[domain.UnitType]UnitBaseCost{
	domain.UnitSmallCargo:       {Metal: 2000, Crystal: 2000, Deuterium: 0, BuildTime: 300},
	domain.UnitLargeCargo:       {Metal: 6000, Crystal: 6000, Deuterium: 0, BuildTime: 480},
	domain.UnitLightFighter:     {Metal: 10000, Crystal: 6000, Deuterium: 2000, BuildTime: 1200},
	domain.UnitHeavyFighter:     {Metal: 25000, Crystal: 15000, Deuterium: 5000, BuildTime: 2400},
	domain.UnitCruiser:          {Metal: 10000, Crystal: 20000, Deuterium: 10000, BuildTime: 600},
	domain.UnitBattleship:       {Metal: 50000, Crystal: 25000, Deuterium: 15000, BuildTime: 4800},
	domain.UnitBattlecruiser:    {Metal: 30000, Crystal: 40000, Deuterium: 15000, BuildTime: 4200},
	domain.UnitBomber:           {Metal: 50000, Crystal: 50000, Deuterium: 25000, BuildTime: 12000},
	domain.UnitDestroyer:        {Metal: 60000, Crystal: 50000, Deuterium: 15000, BuildTime: 9000},
	domain.UnitDeathstar:        {Metal: 1000000, Crystal: 1000000, Deuterium: 500000, BuildTime: 360000},
	domain.UnitRecycler:         {Metal: 10000, Crystal: 6000, Deuterium: 2000, BuildTime: 900},
	domain.UnitEspionageProbe:   {Metal: 0, Crystal: 1000, Deuterium: 0, BuildTime: 1800},
	domain.UnitSolarSatellite:   {Metal: 0, Crystal: 2000, Deuterium: 500, BuildTime: 180},
	domain.UnitCrawler:          {Metal: 8000, Crystal: 4000, Deuterium: 2000, BuildTime: 600},
	domain.UnitColonyShip:       {Metal: 10000, Crystal: 10000, Deuterium: 0, BuildTime: 3000},
	domain.UnitReaper:           {Metal: 140000, Crystal: 40000, Deuterium: 20000, BuildTime: 6000},
	domain.UnitPathfinder:       {Metal: 20000, Crystal: 10000, Deuterium: 10000, BuildTime: 4500},
}

var defenseBaseCosts = map[domain.UnitType]UnitBaseCost{
	domain.UnitRocketLauncher:     {Metal: 2000, Crystal: 0, Deuterium: 0, BuildTime: 600},
	domain.UnitLightLaser:        {Metal: 1500, Crystal: 500, Deuterium: 0, BuildTime: 660},
	domain.UnitHeavyLaser:        {Metal: 6000, Crystal: 2000, Deuterium: 0, BuildTime: 1320},
	domain.UnitGaussCannon:       {Metal: 20000, Crystal: 15000, Deuterium: 2000, BuildTime: 2700},
	domain.UnitIonCannon:         {Metal: 2000, Crystal: 6000, Deuterium: 0, BuildTime: 960},
	domain.UnitPlasmaTurret:     {Metal: 50000, Crystal: 50000, Deuterium: 30000, BuildTime: 5400},
	domain.UnitSmallShieldDome:   {Metal: 10000, Crystal: 10000, Deuterium: 0, BuildTime: 1200},
	domain.UnitLargeShieldDome:   {Metal: 50000, Crystal: 50000, Deuterium: 0, BuildTime: 6000},
	domain.UnitAntiBallisticMissiles: {Metal: 8000, Crystal: 0, Deuterium: 2000, BuildTime: 900},
	domain.UnitInterplanetaryMissiles: {Metal: 12500, Crystal: 2500, Deuterium: 10000, BuildTime: 1500},
}

func CalculateUnitCost(unitType domain.UnitType) (metal, crystal, deuterium int64, buildTime int64) {
	if cost, ok := shipBaseCosts[unitType]; ok {
		return cost.Metal, cost.Crystal, cost.Deuterium, cost.BuildTime
	}
	if cost, ok := defenseBaseCosts[unitType]; ok {
		return cost.Metal, cost.Crystal, cost.Deuterium, cost.BuildTime
	}
	return 0, 0, 0, 0
}

func CalculateShipCost(unitType domain.UnitType) (metal, crystal, deuterium int64, buildTime int64) {
	if cost, ok := shipBaseCosts[unitType]; ok {
		return cost.Metal, cost.Crystal, cost.Deuterium, cost.BuildTime
	}
	return 0, 0, 0, 0
}

func CalculateDefenseCost(unitType domain.UnitType) (metal, crystal, deuterium int64, buildTime int64) {
	if cost, ok := defenseBaseCosts[unitType]; ok {
		return cost.Metal, cost.Crystal, cost.Deuterium, cost.BuildTime
	}
	return 0, 0, 0, 0
}

func GetUnitCost(unitType domain.UnitType) UnitBaseCost {
	if cost, ok := shipBaseCosts[unitType]; ok {
		return cost
	}
	if cost, ok := defenseBaseCosts[unitType]; ok {
		return cost
	}
	return UnitBaseCost{}
}
