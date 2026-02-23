package domain

import "fmt"

type Resources struct {
	Metal     int64
	Crystal   int64
	Deuterium int64
	Energy    int64
}

func (r Resources) IsZero() bool {
	return r.Metal == 0 && r.Crystal == 0 && r.Deuterium == 0 && r.Energy == 0
}

func (r Resources) Add(other Resources) Resources {
	return Resources{
		Metal:     r.Metal + other.Metal,
		Crystal:   r.Crystal + other.Crystal,
		Deuterium: r.Deuterium + other.Deuterium,
		Energy:    r.Energy + other.Energy,
	}
}

func (r Resources) Sub(other Resources) Resources {
	return Resources{
		Metal:     r.Metal - other.Metal,
		Crystal:   r.Crystal - other.Crystal,
		Deuterium: r.Deuterium - other.Deuterium,
		Energy:    r.Energy - other.Energy,
	}
}

func (r Resources) Mul(factor float64) Resources {
	return Resources{
		Metal:     int64(float64(r.Metal) * factor),
		Crystal:   int64(float64(r.Crystal) * factor),
		Deuterium: int64(float64(r.Deuterium) * factor),
		Energy:    int64(float64(r.Energy) * factor),
	}
}

func (r Resources) CanAfford(cost Resources) bool {
	return r.Metal >= cost.Metal &&
		r.Crystal >= cost.Crystal &&
		r.Deuterium >= cost.Deuterium
}

func (r Resources) String() string {
	return fmt.Sprintf("M: %d, C: %d, D: %d, E: %d", r.Metal, r.Crystal, r.Deuterium, r.Energy)
}

var ZeroResources = Resources{}
