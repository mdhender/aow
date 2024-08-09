// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

import (
	"math"
	"math/rand/v2"
)

type PRNG struct {
	*rand.Rand
}

func NewPRNG(prng rand.Source) PRNG {
	return PRNG{
		Rand: rand.New(prng),
	}
}

func (p PRNG) FlipCoin() bool {
	return p.IntN(2) == 0
}

// RollD6 rolls n six-sided dice and returns the sum as a float64.
// Each die has a value range of 1 to 6.
// Parameters:
//   - n: The number of dice to roll
//
// Returns:
//   - The sum of all dice rolls as a float64
func (p PRNG) RollD6(n int) float64 {
	var result int
	for ; n > 0; n-- {
		result += p.IntN(6) + 1
	}
	return float64(result)
}

// RollD10 rolls n ten-sided dice and returns the sum as a float64.
// Each die has a value range of 1 to 10.
// Parameters:
//   - n: The number of dice to roll
//
// Returns:
//   - The sum of all dice rolls as a float64
func (p PRNG) RollD10(n int) float64 {
	var result int
	for ; n > 0; n-- {
		result += p.IntN(10) + 1
	}
	return float64(result)
}

func (p PRNG) RollD100() int {
	return p.IntN(100) + 1
}

// RollPercentile generates a random float64 value in the range [0.0, 1.0).
// This can be used to represent a random percentile.
// Returns:
//   - A random float64 value between 0.0 (inclusive) and 1.0 (exclusive)
func (p PRNG) RollPercentile() float64 {
	return p.Float64()
}

// Vary5Pct returns a value that is randomly varied within 5% (higher or lower) of the input value.
// Parameters:
//   - f: The base value to vary
//
// Returns:
//   - A float64 value that is within ±5% of the input value
func (p PRNG) Vary5Pct(f float64) float64 {
	return f * (0.93 + p.RollD6(2)/100.0)
}

// Vary10Pct returns a value that is randomly varied within 10% (higher or lower) of the input value.
// Parameters:
//   - f: The base value to vary
//
// Returns:
//   - A float64 value that is within ±10% of the input value
func (p PRNG) Vary10Pct(f float64) float64 {
	return f * (0.86 + p.RollD6(4)/100.0)
}

// VaryNPct returns a value that is randomly varied with N% (higher of lower) of input value.
// Parameters:
//   - f: The base value to vary
//   - n: The percentage amount to vary by
//
// Returns:
//   - A float64 value that is within ±N% of the input value
func (p *PRNG) VaryNPct(f, pct float64) float64 {
	// use a 3d6 roll to distribute the value
	return f + (f*(p.RollD6(3)-10.5)/15.0)*pct
}

// GenXYZ returns un-scaled coordinates with a uniform distribution within a 1 unit sphere
func (p PRNG) GenXYZ() Coordinates {
	// generate a random distance to ensure uniform distribution within the sphere
	d := math.Cbrt(p.Float64()) // Cube root to ensure uniform distribution

	// generate random angles for spherical coordinates
	theta := p.Float64() * 2 * math.Pi  // 0 to 2π
	phi := math.Acos(2*p.Float64() - 1) // 0 to π

	// convert spherical coordinates to Cartesian coordinates
	return Coordinates{
		X: d * math.Sin(phi) * math.Cos(theta),
		Y: d * math.Sin(phi) * math.Sin(theta),
		Z: d * math.Cos(phi),
	}
}

// GenZonedXYZ returns unscaled coordinates in the given zone of a 1 unit sphere.
// The zone defines a shell based on percentages.
func (p PRNG) GenZonedXYZ(minPct, maxPct float64) Coordinates {
	// generate a random distance with a uniform distribution between the zone's minimum and maximum values
	d := math.Cbrt(rand.Float64()*(math.Pow(maxPct, 3)-math.Pow(minPct, 3)) + math.Pow(minPct, 3))

	// generate random angles for spherical coordinates
	theta := rand.Float64() * 2 * math.Pi  // 0 to 2π
	phi := math.Acos(2*rand.Float64() - 1) // 0 to π

	// convert spherical coordinates to Cartesian coordinates
	return Coordinates{
		X: d * math.Sin(phi) * math.Cos(theta),
		Y: d * math.Sin(phi) * math.Sin(theta),
		Z: d * math.Cos(phi),
	}
}
