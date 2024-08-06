// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

import (
	"math"
)

type PopulationModel struct {
	YoungPopulationI, IntermediatePopulationI, OldPopulationI, DiskPopulationII, HaloPopulationII populationModel_t
	CombinedDensity                                                                               float64
}

type populationModel_t struct {
	Density  float64 // star systems per cubic parsec
	BaseAge  float64
	AgeRange float64
}

// BasicPopulationModelTable returns a population model table for a region of space similar to Sol's neighborhood.
// It uses the values from p25 of the book.
func BasicPopulationModelTable() PopulationModel {
	return PopulationModel{
		YoungPopulationI:        populationModel_t{Density: 0.0344, BaseAge: 0.0, AgeRange: 2.0},
		IntermediatePopulationI: populationModel_t{Density: 0.0272, BaseAge: 2.0, AgeRange: 3.0},
		OldPopulationI:          populationModel_t{Density: 0.0158, BaseAge: 5.0, AgeRange: 3.0},
		DiskPopulationII:        populationModel_t{Density: 0.00339, BaseAge: 8.0, AgeRange: 1.5},
		HaloPopulationII:        populationModel_t{Density: 0.000339, BaseAge: 9.5, AgeRange: 3.0},
		CombinedDensity:         0.081129,
	}
}

// AdvancedPopulationModelTable returns a population model table for a region of space that might differ from Sol's neighborhood.
//
// Parameters:
//   - r: The distance from the center of the galaxy in parsecs along the galactic plane
//   - h: The distance above or below the galactic plane in parsecs
func AdvancedPopulationModelTable(r, h float64) PopulationModel {
	var pm = PopulationModel{
		YoungPopulationI:        populationModel_t{BaseAge: 0.0, AgeRange: 2.0},
		IntermediatePopulationI: populationModel_t{BaseAge: 2.0, AgeRange: 3.0},
		OldPopulationI:          populationModel_t{BaseAge: 5.0, AgeRange: 3.0},
		DiskPopulationII:        populationModel_t{BaseAge: 8.0, AgeRange: 1.5},
		HaloPopulationII:        populationModel_t{BaseAge: 9.5, AgeRange: 3.0},
	}

	// h must be non-negative
	h = math.Abs(h)

	// calculate the density for each stellar population using the equation from p26 of the book.
	pm.YoungPopulationI.Density = 0.373 * math.Pow(math.E, -(r/3_500.0)) * math.Pow(math.E, -(h/200))
	pm.IntermediatePopulationI.Density = 0.280 * math.Pow(math.E, -(r/3_500)) * math.Pow(math.E, -(h/400))
	pm.OldPopulationI.Density = 0.160 * math.Pow(math.E, -(r/3_500)) * math.Pow(math.E, -(h/700))
	pm.DiskPopulationII.Density = 0.0339 * math.Pow(math.E, -(r/3_500)) * math.Pow(math.E, -(h/1_000))
	pm.HaloPopulationII.Density = 0.00339 * math.Pow(math.E, -(r/3_500)) * math.Pow(math.E, -(h/2_000))

	// the combined density is saved for future calculations.
	pm.CombinedDensity = pm.YoungPopulationI.Density + pm.IntermediatePopulationI.Density + pm.OldPopulationI.Density + pm.DiskPopulationII.Density + pm.HaloPopulationII.Density

	return pm
}

// StellarPopulation_e is a grouping of stellar systems that have similar characteristics.
type StellarPopulation_e int

const (
	YoungPopulationI StellarPopulation_e = iota
	IntermediatePopulationI
	OldPopulationI
	DiskPopulationII
	HaloPopulationII
)
