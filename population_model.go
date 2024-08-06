// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

import (
	"math"
)

// PopulationModelForEarthLikeSystems returns the smallest volume of space (in cubic parsecs)
// that would likely contain the requested number of systems with Earth-like planets.
//
// Parameters:
//   - n: The number of systems Earth-like planets to target.
//   - tweak: A tweak factor to adjust the volume of space returned.
//
// If the tweak is between 0 and 5, the volume returned is adjusted by the tweak value.
//
// Returns:
//   - The suggested population model to use for the systems. This includes the stellar population
//     and the smallest volume of space that would likely contain the requested number of Earth-like planets.
func PopulationModelForEarthLikeSystems(n int, tweak float64) PopulationModel_t {
	pm := BasicPopulationModelTable()
	cubicParsesPerSolLikeSystem := 150.0
	if 0 < tweak && tweak <= 5 {
		cubicParsesPerSolLikeSystem += tweak
	}
	// the formula from p24 of the book
	pm.Volume = float64(n) * 2.0 * cubicParsesPerSolLikeSystem
	return pm
}

// PopulationModelForSolLikeNeighborhood returns the smallest volume of space (in cubic parsecs)
// that would likely contain the requested number of systems based on the population
// of Sol's neighborhood.
//
// Parameters:
//   - n: The number of systems to target.
//   - tweak: A tweak factor to adjust the volume of space returned.
//
// If the tweak is between 0 and 1, the volume returned is adjusted by the tweak value.
//
// Returns:
//   - The suggested population model to use for the systems. This includes the stellar population
//     and the smallest volume of space that would likely contain the requested number of systems.
func PopulationModelForSolLikeNeighborhood(n int, tweak float64) PopulationModel_t {
	pm := BasicPopulationModelTable()
	cubicParsecsPerStarSystem := 12.0
	if 0 < tweak && tweak <= 1 {
		cubicParsecsPerStarSystem += tweak
	}
	// the formula from p24 of the book
	pm.Volume = float64(n) * cubicParsecsPerStarSystem
	return pm
}

// PopulationModelForOtherNeighborhoods returns the smallest volume of space (in cubic parsecs)
// that would likely contain the requested number of systems for a neighborhood that is
// the given distance (in parsecs) from the center of the galaxy, and the given distance
// (in parsecs) above or below the galactic plane.
//
// Parameters:
//   - n: The number of systems to target.
//   - r: The distance (in parsecs) from the center of the galaxy.
//   - h: The distance (in parsecs) above or below the galactic plane.
//   - tweak: A tweak factor to adjust the volume of space returned.
//
// If the tweak is between 0 and 1, the volume returned is adjusted by the tweak value.
//
// Returns:
//   - The suggested population model to use for the systems. This includes the stellar population
//     and the smallest volume of space that would likely contain the requested number of systems.
func PopulationModelForOtherNeighborhoods(n int, r, h float64, tweak float64) PopulationModel_t {
	pm := AdvancedPopulationModelTable(math.Abs(r), math.Abs(h))
	cubicParsecsPerStarSystem := 1 / pm.CombinedDensity
	if 0 < tweak && tweak <= 1 {
		cubicParsecsPerStarSystem = cubicParsecsPerStarSystem + cubicParsecsPerStarSystem*tweak
	}
	// the formula from p24 of the book
	pm.Volume = float64(n) * cubicParsecsPerStarSystem
	return pm
}

type PopulationModel_t struct {
	YoungPopulationI, IntermediatePopulationI, OldPopulationI, DiskPopulationII, HaloPopulationII populationModel_t
	CombinedDensity                                                                               float64
	Volume                                                                                        float64 // volume of the population in cubic parsecs
}

type populationModel_t struct {
	Density  float64 // star systems per cubic parsec
	BaseAge  float64
	AgeRange float64
}

// BasicPopulationModelTable returns a population model table for a region of space similar to Sol's neighborhood.
// It uses the values from p25 of the book.
func BasicPopulationModelTable() PopulationModel_t {
	return PopulationModel_t{
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
func AdvancedPopulationModelTable(r, h float64) PopulationModel_t {
	var pm = PopulationModel_t{
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
