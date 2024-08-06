// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

// routines to help users decide on the volume of space that their catalog will cover

import "math"

// VolumeForEarthLikeSystems returns the smallest volume of space (in cubic parsecs)
// that would likely contain the requested number of systems with Earth-like planets.
//
// Parameters:
//   - n: The number of systems Earth-like planets to target.
//   - tweak: A tweak factor to adjust the volume of space returned.
//
// If the tweak is between 0 and 5, the volume returned is adjusted by the tweak value.
//
// Returns:
//   - The smallest volume of space that would likely contain the requested number of Earth-like planets.
func VolumeForEarthLikeSystems(n int, tweak float64) float64 {
	cubicParsesPerSolLikeSystem := 150.0
	if 0 < tweak && tweak <= 5 {
		cubicParsesPerSolLikeSystem += tweak
	}
	// the formula from p24 of the book
	return float64(n) * 2.0 * cubicParsesPerSolLikeSystem
}

// VolumeForSolLikeNeighborhood returns the smallest volume of space (in cubic parsecs)
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
//   - The smallest volume of space that would likely contain the requested number of systems.
func VolumeForSolLikeNeighborhood(n int, tweak float64) float64 {
	cubicParsecsPerStarSystem := 12.0
	if 0 < tweak && tweak <= 1 {
		cubicParsecsPerStarSystem += tweak
	}
	// the formula from p24 of the book
	return float64(n) * cubicParsecsPerStarSystem
}

// VolumeForOtherNeighborhoods returns the smallest volume of space (in cubic parsecs)
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
//   - The smallest volume of space that would likely contain the requested number of systems.
func VolumeForOtherNeighborhoods(n int, r, h float64, tweak float64) float64 {
	pm := AdvancedPopulationModelTable(math.Abs(r), math.Abs(h))
	cubicParsecsPerStarSystem := 1 / pm.CombinedDensity
	if 0 < tweak && tweak <= 1 {
		cubicParsecsPerStarSystem = cubicParsecsPerStarSystem + cubicParsecsPerStarSystem*tweak
	}
	// the formula from p24 of the book
	return float64(n) * cubicParsecsPerStarSystem
}
