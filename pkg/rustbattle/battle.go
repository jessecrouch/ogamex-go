package rustbattle

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -L/home/bolt/Documents/ogamex-go/storage/rust-libs -lbattle_engine_ffi
#include <stdlib.h>
#include <string.h>

extern char* fight_battle_rounds(char* input_json);
*/
import "C"

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
	UnitID      int16   `json:"unit_id"`
	Amount      uint32  `json:"amount"`
	AttackPower float32 `json:"attack_power"`
	ShieldPoints float32 `json:"shield_points"`
	HullPlating float32 `json:"hull_plating"`
	Rapidfire   map[int16]uint16 `json:"rapidfire"`
}

type BattleOutput struct {
	Results []RoundResult `json:"results"`
	Winner  string        `json:"winner"`
}

type RoundResult struct {
	AttackerShips map[int16]UnitCount `json:"attacker_ships"`
	DefenderShips map[int16]UnitCount `json:"defender_ships"`
	AttackerLosses map[int16]UnitCount `json:"attacker_losses"`
	DefenderLosses map[int16]UnitCount `json:"defender_losses"`
}

type UnitCount struct {
	UnitID int16  `json:"unit_id"`
	Amount uint32 `json:"amount"`
}

func SimulateBattle(input BattleInput) (BattleOutput, error) {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return BattleOutput{}, fmt.Errorf("failed to marshal input: %w", err)
	}

	cInput := C.CString(string(inputJSON))
	defer C.free(unsafe.Pointer(cInput))

	cOutput := C.fight_battle_rounds(cInput)
	if cOutput == nil {
		return BattleOutput{}, fmt.Errorf("battle engine returned nil")
	}
	defer C.free(unsafe.Pointer(cOutput))

	outputJSON := C.GoString(cOutput)

	var output BattleOutput
	err = json.Unmarshal([]byte(outputJSON), &output)
	if err != nil {
		return BattleOutput{}, fmt.Errorf("failed to unmarshal output: %w", err)
	}

	return output, nil
}
