//go:build !cgo
// +build !cgo

package rustbattle

import "fmt"

type BattleInput struct {
	AttackerFleets []FleetInput `json:"attacker_fleets"`
	DefenderFleets []FleetInput `json:"defender_fleets"`
}

type FleetInput struct {
	FleetMissionID uint32            `json:"fleet_mission_id"`
	OwnerID       uint32            `json:"owner_id"`
	Units         map[int16]UnitInfo `json:"units"`
}

type UnitInfo struct {
	UnitID       int16   `json:"unit_id"`
	Amount       uint32  `json:"amount"`
	AttackPower  float32 `json:"attack_power"`
	ShieldPoints float32 `json:"shield_points"`
	HullPlating  float32 `json:"hull_plating"`
	Rapidfire    map[int16]uint16 `json:"rapidfire"`
}

type BattleOutput struct {
	Results []RoundResult `json:"results"`
	Winner  string        `json:"winner"`
}

type RoundResult struct {
	AttackerShips   map[int16]UnitCount `json:"attacker_ships"`
	DefenderShips   map[int16]UnitCount `json:"defender_ships"`
	AttackerLosses map[int16]UnitCount `json:"attacker_losses"`
	DefenderLosses map[int16]UnitCount `json:"defender_losses"`
}

type UnitCount struct {
	UnitID int16  `json:"unit_id"`
	Amount uint32 `json:"amount"`
}

func SimulateBattle(input BattleInput) (BattleOutput, error) {
	return BattleOutput{}, fmt.Errorf("battle engine not available: CGO disabled")
}
