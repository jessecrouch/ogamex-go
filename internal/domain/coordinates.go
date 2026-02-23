package domain

import "fmt"

type Coordinates struct {
	Galaxy   int
	System   int
	Position int
}

func (c Coordinates) String() string {
	return fmt.Sprintf("[%d:%d:%d]", c.Galaxy, c.System, c.Position)
}

func (c Coordinates) IsValid() bool {
	return c.Galaxy >= 1 && c.Galaxy <= 10 &&
		c.System >= 1 && c.System <= 499 &&
		c.Position >= 1 && c.Position <= 15
}

func (c Coordinates) DistanceTo(other Coordinates) int {
	dx := c.Galaxy - other.Galaxy
	if dx < 0 {
		dx = -dx
	}

	dy := c.System - other.System
	if dy < 0 {
		dy = -dy
	}

	dz := c.Position - other.Position
	if dz < 0 {
		dz = -dz
	}

	return dx*20000 + dy*5 + dz*20000
}

func (c Coordinates) IsMoon() bool {
	return c.Position == 0 || c.Position > 15
}
