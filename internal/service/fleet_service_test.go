package service

import (
	"testing"

	"ogamex-go/internal/schema"
)

func TestColonizeName(t *testing.T) {
	tests := []struct {
		name     string
		galaxy   int
		system   int
		position int
		want     string
	}{
		{"Basic", 1, 1, 1, "Planet"},
		{"Different coords", 2, 50, 7, "Planet"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := colonizeName(tt.galaxy, tt.system, tt.position)
			if got != tt.want {
				t.Errorf("colonizeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateMoonChance(t *testing.T) {
	svc := &FleetService{}

	tests := []struct {
		name          string
		destroyed     int64
		wantMaxChance int
	}{
		{"Below threshold", 50000, 0},
		{"At threshold", 100000, 0},
		{"Above threshold", 200000, 1},
		{"High destruction", 2000000, 20},
		{"Very high destruction", 5000000, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.calculateMoonChance(tt.destroyed)
			if got > tt.wantMaxChance {
				t.Errorf("calculateMoonChance(%d) = %d, want <= %d", tt.destroyed, got, tt.wantMaxChance)
			}
		})
	}
}

func TestCalculateDestroyedResources(t *testing.T) {
	svc := &FleetService{}

	attackerLosses := map[int16]int16{
		202: 10,
		207: 5,
	}
	defenderLosses := map[int16]int16{
		204: 20,
		401: 50,
	}

	total := svc.calculateDestroyedResources(attackerLosses, defenderLosses)

	if total <= 0 {
		t.Errorf("calculateDestroyedResources() = %d, want > 0", total)
	}
}

func TestPlanetModelDefaults(t *testing.T) {
	planet := schema.Planet{
		Name:     "Test Planet",
		Galaxy:   1,
		System:   1,
		Position: 1,
	}

	if planet.MetalCapacity == 0 {
		planet.MetalCapacity = 10000
	}
	if planet.CrystalCapacity == 0 {
		planet.CrystalCapacity = 10000
	}
	if planet.DeuteriumCapacity == 0 {
		planet.DeuteriumCapacity = 10000
	}

	if planet.MetalCapacity != 10000 {
		t.Errorf("MetalCapacity = %d, want 10000", planet.MetalCapacity)
	}
}
