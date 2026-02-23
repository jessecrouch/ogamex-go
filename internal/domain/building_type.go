package domain

type BuildingType int

const (
	BuildingMetalMine BuildingType = iota + 1
	BuildingCrystalMine
	BuildingDeuteriumSynthesizer
	BuildingSolarPlant
	BuildingFusionPlant
	BuildingMetalStorage
	BuildingCrystalStorage
	BuildingDeuteriumStorage
	BuildingMetalMine2
	BuildingCrystalMine2
	BuildingDeuteriumSynthesizer2
	BuildingSolarSatellite
	BuildingCrawler
	BuildingSpaceDock
	BuildingNaniteFactory
	BuildingTerraformer
	BuildingMissileSilo
)

type BuildingQueueItem struct {
	BuildingID BuildingType
	Level      int
	StartTime  int64
	EndTime    int64
}
