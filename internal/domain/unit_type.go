package domain

type UnitType int

const (
	UnitSmallCargo UnitType = iota + 1
	UnitLargeCargo
	UnitLightFighter
	UnitHeavyFighter
	UnitCruiser
	UnitBattleship
	UnitBattlecruiser
	UnitBomber
	UnitDestroyer
	UnitDeathstar
	UnitRecycler
	UnitEspionageProbe
	UnitSolarSatellite
	UnitCrawler
	UnitColonyShip
	UnitReaper
	UnitPathfinder
	UnitRocketLauncher
	UnitLightLaser
	UnitHeavyLaser
	UnitGaussCannon
	UnitIonCannon
	UnitPlasmaTurret
	UnitSmallShieldDome
	UnitLargeShieldDome
	UnitAntiBallisticMissiles
	UnitInterplanetaryMissiles
)

func (u UnitType) IsShip() bool {
	return u >= UnitSmallCargo && u <= UnitPathfinder
}

func (u UnitType) IsDefense() bool {
	return u >= UnitRocketLauncher && u <= UnitLargeShieldDome
}

func (u UnitType) IsMissile() bool {
	return u == UnitAntiBallisticMissiles || u == UnitInterplanetaryMissiles
}

func (u UnitType) String() string {
	switch u {
	case UnitSmallCargo:
		return "Small Cargo"
	case UnitLargeCargo:
		return "Large Cargo"
	case UnitLightFighter:
		return "Light Fighter"
	case UnitHeavyFighter:
		return "Heavy Fighter"
	case UnitCruiser:
		return "Cruiser"
	case UnitBattleship:
		return "Battleship"
	case UnitBattlecruiser:
		return "Battlecruiser"
	case UnitBomber:
		return "Bomber"
	case UnitDestroyer:
		return "Destroyer"
	case UnitDeathstar:
		return "Deathstar"
	case UnitRecycler:
		return "Recycler"
	case UnitEspionageProbe:
		return "Espionage Probe"
	case UnitSolarSatellite:
		return "Solar Satellite"
	case UnitCrawler:
		return "Crawler"
	case UnitColonyShip:
		return "Colony Ship"
	case UnitReaper:
		return "Reaper"
	case UnitPathfinder:
		return "Pathfinder"
	case UnitRocketLauncher:
		return "Rocket Launcher"
	case UnitLightLaser:
		return "Light Laser"
	case UnitHeavyLaser:
		return "Heavy Laser"
	case UnitGaussCannon:
		return "Gauss Cannon"
	case UnitIonCannon:
		return "Ion Cannon"
	case UnitPlasmaTurret:
		return "Plasma Turret"
	case UnitSmallShieldDome:
		return "Small Shield Dome"
	case UnitLargeShieldDome:
		return "Large Shield Dome"
	case UnitAntiBallisticMissiles:
		return "Anti-Ballistic Missiles"
	case UnitInterplanetaryMissiles:
		return "Interplanetary Missiles"
	default:
		return "Unknown"
	}
}
