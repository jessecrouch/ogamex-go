package domain

type MissionType int

const (
	MissionAttack MissionType = iota + 1
	MissionAttackHome
	MissionTransport
	MissionDeploy
	MissionHold
	MissionSpy
	MissionHarvest
	MissionColonize
	MissionRecycle
	MissionDestroy
	MissionMissile
	MissionACS
	MissionExpedition
)

func (m MissionType) String() string {
	switch m {
	case MissionAttack:
		return "Attack"
	case MissionAttackHome:
		return "Attack Home"
	case MissionTransport:
		return "Transport"
	case MissionDeploy:
		return "Deploy"
	case MissionHold:
		return "Hold"
	case MissionSpy:
		return "Spy"
	case MissionHarvest:
		return "Harvest"
	case MissionColonize:
		return "Colonize"
	case MissionRecycle:
		return "Recycle"
	case MissionDestroy:
		return "Destroy"
	case MissionMissile:
		return "Missile"
	case MissionACS:
		return "ACS"
	case MissionExpedition:
		return "Expedition"
	default:
		return "Unknown"
	}
}

type CharacterClass int

const (
	CharacterClassNone CharacterClass = iota
	CharacterClassCollector
	CharacterClassGeneral
	CharacterClassDiscoverer
)

func (c CharacterClass) String() string {
	switch c {
	case CharacterClassNone:
		return "None"
	case CharacterClassCollector:
		return "Collector"
	case CharacterClassGeneral:
		return "General"
	case CharacterClassDiscoverer:
		return "Discoverer"
	default:
		return "Unknown"
	}
}
