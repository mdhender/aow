// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

// RollD6 rolls n six-sided dice and returns the sum as a float64.
// Each die has a value range of 1 to 6.
// Parameters:
//   - n: The number of dice to roll
//
// Returns:
//   - The sum of all dice rolls as a float64
func (g *Generator) RollD6(n int) float64 {
	var result int
	for ; n > 0; n-- {
		result += g.prng.IntN(6) + 1
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
func (g *Generator) RollD10(n int) float64 {
	var result int
	for ; n > 0; n-- {
		result += g.prng.IntN(10) + 1
	}
	return float64(result)
}

// RollPercentile generates a random float64 value in the range [0.0, 1.0).
// This can be used to represent a random percentile.
// Returns:
//   - A random float64 value between 0.0 (inclusive) and 1.0 (exclusive)
func (g *Generator) RollPercentile() float64 {
	return g.prng.Float64()
}

// Vary5pct returns a value that is randomly varied within 5% (higher or lower) of the input value.
// Parameters:
//   - f: The base value to vary
//
// Returns:
//   - A float64 value that is within ±5% of the input value
func (g *Generator) Vary5pct(f float64) float64 {
	return f * (0.93 + g.RollD6(2)/100.0)
}

// Vary10pct returns a value that is randomly varied within 10% (higher or lower) of the input value.
// Parameters:
//   - f: The base value to vary
//
// Returns:
//   - A float64 value that is within ±10% of the input value
func (g *Generator) Vary10pct(f float64) float64 {
	return f * (0.86 + g.RollD6(4)/100.0)
}
