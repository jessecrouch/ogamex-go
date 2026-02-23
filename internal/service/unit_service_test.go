package service

import (
	"testing"

	"ogamex-go/internal/schema"
)

func TestGetUnitCost(t *testing.T) {
	svc := &UnitService{}

	tests := []struct {
		name       string
		unitID     int
		wantMetal  int64
		wantCrystal int64
		wantDeut   int64
	}{
		{"Small Cargo", 202, 2000, 2000, 0},
		{"Large Cargo", 203, 6000, 6000, 0},
		{"Light Fighter", 204, 10000, 6000, 2000},
		{"Rocket Launcher (Defense)", 401, 2000, 0, 0},
		{"Gauss Cannon (Defense)", 405, 20000, 15000, 2000},
		{"Invalid Unit", 999, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metal, crystal, deut, _ := svc.GetUnitCost(tt.unitID)
			if metal != tt.wantMetal || crystal != tt.wantCrystal || deut != tt.wantDeut {
				t.Errorf("GetUnitCost(%d) = (%d, %d, %d), want (%d, %d, %d)",
					tt.unitID, metal, crystal, deut, tt.wantMetal, tt.wantCrystal, tt.wantDeut)
			}
		})
	}
}

func TestIsShipUnit(t *testing.T) {
	svc := &UnitService{}

	tests := []struct {
		name   string
		unitID int
		want   bool
	}{
		{"Small Cargo is ship", 202, true},
		{"Battleship is ship", 207, true},
		{"Rocket Launcher is not ship", 401, false},
		{"Invalid is not ship", 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.IsShipUnit(tt.unitID); got != tt.want {
				t.Errorf("IsShipUnit(%d) = %v, want %v", tt.unitID, got, tt.want)
			}
		})
	}
}

func TestIsDefenseUnit(t *testing.T) {
	svc := &UnitService{}

	tests := []struct {
		name   string
		unitID int
		want   bool
	}{
		{"Rocket Launcher is defense", 401, true},
		{"Gauss Cannon is defense", 405, true},
		{"Small Cargo is not defense", 202, false},
		{"Invalid is not defense", 999, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.IsDefenseUnit(tt.unitID); got != tt.want {
				t.Errorf("IsDefenseUnit(%d) = %v, want %v", tt.unitID, got, tt.want)
			}
		})
	}
}

func TestIsValidUnit(t *testing.T) {
	svc := &UnitService{}

	tests := []struct {
		name   string
		unitID int
		want   bool
	}{
		{"Small Cargo valid", 202, true},
		{"Defense valid", 405, true},
		{"Invalid unit", 999, false},
		{"Zero", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.IsValidUnit(tt.unitID); got != tt.want {
				t.Errorf("IsValidUnit(%d) = %v, want %v", tt.unitID, got, tt.want)
			}
		})
	}
}

func TestCalculateLoot(t *testing.T) {
	svc := &FleetService{}

	tests := []struct {
		name      string
		metal     int64
		crystal   int64
		deuterium int64
		wantMax   int64
	}{
		{"Normal resources", 10000, 5000, 2000, 8500},
		{"No resources", 0, 0, 0, 0},
		{"Only metal", 10000, 0, 0, 5000},
		{"Only crystal", 0, 10000, 0, 5000},
		{"Only deuterium", 0, 0, 10000, 5000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			planet := &schema.Planet{
				Metal:     tt.metal,
				Crystal:   tt.crystal,
				Deuterium: tt.deuterium,
			}
			result := svc.calculateLoot(planet)
			total := result.Metal + result.Crystal + result.Deuterium

			if tt.wantMax > 0 && total > tt.wantMax {
				t.Errorf("calculateLoot() total = %d, want <= %d", total, tt.wantMax)
			}
		})
	}
}
