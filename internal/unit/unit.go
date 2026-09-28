package unit

import (
	"errors"
	"strings"
)

var ErrIncompatibleUnits = errors.New("incompatible units")

type Unit string

const (
	G  Unit = "g"
	KG Unit = "kg"
	MG Unit = "mg"
	OZ Unit = "oz"
	LB Unit = "lb"
	ML Unit = "ml"
	L  Unit = "L"
	DL Unit = "dL"
	UN Unit = "un"
)

type Dimension string

const (
	Mass   Dimension = "mass"
	Volume Dimension = "volume"
	Count  Dimension = "count"
)

type definition struct {
	dimension Dimension
	toBase    float64
}

var All = []Unit{G, KG, MG, OZ, LB, ML, L, DL, UN}

var definitions = map[Unit]definition{
	G:  {Mass, 1},
	KG: {Mass, 1000},
	MG: {Mass, 0.001},
	OZ: {Mass, 28.349523125},
	LB: {Mass, 453.59237},
	ML: {Volume, 1},
	L:  {Volume, 1000},
	DL: {Volume, 100},
	UN: {Count, 1},
}

func Parse(s string) (Unit, bool) {
	for _, u := range All {
		if strings.EqualFold(string(u), s) {
			return u, true
		}
	}
	return "", false
}

func (u Unit) Valid() bool {
	_, ok := definitions[u]
	return ok
}

func (u Unit) Dimension() Dimension {
	return definitions[u].dimension
}

func (u Unit) CompatibleWith(other Unit) bool {
	return u.Valid() && other.Valid() && u.Dimension() == other.Dimension()
}

func (u Unit) Compatible() []Unit {
	var units []Unit
	for _, other := range All {
		if u.CompatibleWith(other) {
			units = append(units, other)
		}
	}
	return units
}

func Convert(amount float64, from, to Unit) (float64, error) {
	if !from.CompatibleWith(to) {
		return 0, ErrIncompatibleUnits
	}
	return amount * definitions[from].toBase / definitions[to].toBase, nil
}

func Names() []string {
	names := make([]string, len(All))
	for i, u := range All {
		names[i] = string(u)
	}
	return names
}
