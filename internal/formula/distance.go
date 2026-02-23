package formula

import (
	"math"
)

func CalculateDistance(originX, originY, originZ, destX, destY, destZ int) int {
	dx := originX - destX
	if dx < 0 {
		dx = -dx
	}

	dy := originY - destY
	if dy < 0 {
		dy = -dy
	}

	dz := originZ - destZ
	if dz < 0 {
		dz = -dz
	}

	return dx*20000 + dy*5 + dz*20000
}

func CalculateFlightTime(distance int, speed int, universeSpeed int) int64 {
	if speed <= 0 || distance <= 0 {
		return 0
	}
	// time = 35000 * distance / (speed * 10 * universeSpeed)
	time := float64(35000*distance) / (float64(speed) * 10 * float64(universeSpeed))
	return int64(math.Ceil(time))
}

func CalculateFleetSpeed(baseSpeed int, techLevel int, numShips int, universeSpeed int) int {
	// Simplified: actual speed depends on ship composition
	// Each ship type has a base speed and speed upgrades
	return baseSpeed * universeSpeed
}

func CalculatePositionBonus(position int) (metalBonus, crystalBonus, deuteriumBonus float64) {
	bonuses := map[int][3]float64{
		1:  {1.0, 1.0, 1.0},
		2:  {1.0, 1.0, 1.0},
		3:  {1.0, 1.0, 1.0},
		4:  {1.0, 1.1, 1.0},
		5:  {1.0, 1.0, 1.0},
		6:  {1.0, 1.0, 1.0},
		7:  {1.0, 1.0, 1.0},
		8:  {1.0, 1.2, 1.0},
		9:  {1.0, 1.3, 1.0},
		10: {1.0, 1.4, 1.0},
		11: {1.0, 1.6, 1.0},
		12: {1.0, 1.7, 1.0},
		13: {1.0, 1.0, 1.0},
		14: {1.0, 1.0, 1.0},
		15: {1.0, 1.0, 1.0},
	}

	if b, ok := bonuses[position]; ok {
		return b[0], b[1], b[2]
	}

	return 1.0, 1.0, 1.0
}
