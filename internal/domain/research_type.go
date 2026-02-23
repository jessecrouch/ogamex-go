package domain

type ResearchType int

const (
	ResearchEnergyTechnology ResearchType = iota + 1
	ResearchLaserTechnology
	ResearchIonTechnology
	ResearchHyperspaceTechnology
	ResearchPlasmaTechnology
	ResearchFusionDrive
	ResearchImpulseDrive
	ResearchHyperspaceDrive
	ResearchEspionageTechnology
	ResearchComputerTechnology
	ResearchAstrophysics
	ResearchIntergalacticResearchNetwork
	ResearchGravitonTechnology
	ResearchWeaponsTechnology
	ResearchShieldingTechnology
	ResearchArmorTechnology
)

type ResearchQueueItem struct {
	ResearchID ResearchType
	Level      int
	StartTime  int64
	EndTime    int64
}
