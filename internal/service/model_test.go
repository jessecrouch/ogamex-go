package service

import (
	"testing"
	"time"

	"ogamex-go/internal/schema"
)

func TestBuildingServiceErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"ErrBuildingNotFound", ErrBuildingNotFound, "building not found"},
		{"ErrInsufficientFunds", ErrInsufficientFunds, "insufficient resources"},
		{"ErrBuildingNotAvailable", ErrBuildingNotAvailable, "building not available on this planet"},
		{"ErrQueueFull", ErrQueueFull, "building queue is full"},
		{"ErrTechNotFound", ErrTechNotFound, "technology not found"},
		{"ErrLabBusy", ErrLabBusy, "research lab is busy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.want {
				t.Errorf("Error = %v, want %v", tt.err.Error(), tt.want)
			}
		})
	}
}

func TestFleetServiceErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"ErrFleetNotFound", ErrFleetNotFound, "fleet not found"},
		{"ErrInvalidFleetParams", ErrInvalidFleetParams, "invalid fleet parameters"},
		{"ErrInsufficientFleet", ErrInsufficientFleet, "insufficient ships"},
		{"ErrDestinationOccupied", ErrDestinationOccupied, "destination occupied"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.want {
				t.Errorf("Error = %v, want %v", tt.err.Error(), tt.want)
			}
		})
	}
}

func TestResearchServiceErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"ErrInsufficientFundsResearch", ErrInsufficientFunds, "insufficient resources"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.want {
				t.Errorf("Error = %v, want %v", tt.err.Error(), tt.want)
			}
		})
	}
}

func TestBuildingQueueModel(t *testing.T) {
	now := time.Now()
	endTime := now.Add(1 * time.Hour)

	queue := &schema.BuildingQueue{
		ID:         1,
		PlanetID:   1,
		BuildingID: 1,
		Level:      5,
		StartTime:  now,
		EndTime:    endTime,
		IsCancelled: false,
	}

	if queue.IsCancelled {
		t.Error("Expected IsCancelled to be false")
	}

	queue.IsCancelled = true
	if !queue.IsCancelled {
		t.Error("Expected IsCancelled to be true")
	}
}

func TestResearchQueueModel(t *testing.T) {
	now := time.Now()
	endTime := now.Add(2 * time.Hour)

	queue := &schema.ResearchQueue{
		ID:         1,
		UserID:     1,
		ResearchID: 106,
		Level:      3,
		StartTime:  now,
		EndTime:    endTime,
	}

	if queue.ResearchID != 106 {
		t.Errorf("ResearchID = %d, want 106", queue.ResearchID)
	}
}

func TestUnitQueueModel(t *testing.T) {
	now := time.Now()
	endTime := now.Add(5 * time.Minute)

	queue := &schema.UnitQueue{
		ID:        1,
		PlanetID:  1,
		UnitID:    202,
		Amount:    10,
		StartTime: now,
		EndTime:   endTime,
	}

	if queue.UnitID != 202 {
		t.Errorf("UnitID = %d, want 202", queue.UnitID)
	}
	if queue.Amount != 10 {
		t.Errorf("Amount = %d, want 10", queue.Amount)
	}
}

func TestFleetMissionModel(t *testing.T) {
	now := time.Now()
	arrivalTime := now.Add(1 * time.Hour)
	returnTime := now.Add(2 * time.Hour)

	mission := &schema.FleetMission{
		ID:                1,
		UserID:            1,
		MissionType:       1,
		OriginGalaxy:      1,
		OriginSystem:     1,
		OriginPosition:   1,
		OriginPlanetType: 1,
		TargetGalaxy:     1,
		TargetSystem:     50,
		TargetPosition:   5,
		TargetPlanetType: 1,
		LaunchTime:       now,
		ArrivalTime:      arrivalTime,
		ReturnTime:       returnTime,
		Metal:            1000,
		Crystal:          500,
		Deuterium:        100,
		Status:            0,
	}

	if mission.MissionType != 1 {
		t.Errorf("MissionType = %d, want 1", mission.MissionType)
	}
	if mission.Status != 0 {
		t.Errorf("Status = %d, want 0", mission.Status)
	}
}

func TestUserTechModel(t *testing.T) {
	tech := &schema.UserTech{
		ID:                    1,
		UserID:                1,
		EnergyTechnology:     5,
		LaserTechnology:      5,
		IonTechnology:        3,
		HyperspaceTechnology: 2,
		PlasmaTechnology:     4,
		CombatDrive:          6,
		ImpulseDrive:         5,
		HyperspaceDrive:      4,
		EspionageTechnology:  3,
		ComputerTechnology:   8,
		Astrophysics:         2,
		WeaponsTechnology:    5,
		ShieldingTechnology: 5,
		ArmorTechnology:     5,
		AssemblyTechnology:   3,
	}

	if tech.ComputerTechnology != 8 {
		t.Errorf("ComputerTechnology = %d, want 8", tech.ComputerTechnology)
	}
	if tech.WeaponsTechnology != 5 {
		t.Errorf("WeaponsTechnology = %d, want 5", tech.WeaponsTechnology)
	}
}

func TestPlanetModel(t *testing.T) {
	planet := &schema.Planet{
		ID:        1,
		UserID:    1,
		Name:      "Home Planet",
		Galaxy:    1,
		System:    1,
		Position:  1,
		IsMoon:    false,
		PlanetType: 1,
		Metal:     500000,
		Crystal:   250000,
		Deuterium: 100000,
		MetalCapacity:     1000000,
		CrystalCapacity:   1000000,
		DeuteriumCapacity: 1000000,
		TempMin:      -50,
		TempMax:      50,
		FieldsUsed:   50,
		FieldsMax:    163,
	}

	if planet.Metal != 500000 {
		t.Errorf("Metal = %d, want 500000", planet.Metal)
	}
	if planet.FieldsMax != 163 {
		t.Errorf("FieldsMax = %d, want 163", planet.FieldsMax)
	}
	if planet.IsMoon {
		t.Error("Expected IsMoon to be false")
	}
}
