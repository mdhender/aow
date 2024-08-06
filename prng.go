// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

import "math/rand/v2"

// RollD6 rolls n six-sided dice and returns the sum as a float64.
func (g *Generator) RollD6(n int) float64 {
	return g.prng.RollD6(n)
}

// RollD10 rolls n ten-sided dice and returns the sum as a float64.
func (g *Generator) RollD10(n int) float64 {
	return g.prng.RollD10(n)
}

// RollPercentile generates a random float64 value in the range [0.0, 1.0).
// This can be used to represent a random percentile.
// Returns:
//   - A random float64 value between 0.0 (inclusive) and 1.0 (exclusive)
func (g *Generator) RollPercentile() float64 {
	return g.prng.RollPercentile()
}

// Vary5pct returns a value that is randomly varied within 5% (higher or lower) of the input value.
func (g *Generator) Vary5pct(f float64) float64 {
	return g.prng.Vary5pct(f)
}

// Vary10pct returns a value that is randomly varied within 10% (higher or lower) of the input value.
func (g *Generator) Vary10pct(f float64) float64 {
	return g.prng.Vary10pct(f)
}

type PRNG struct {
	*rand.Rand
}

func NewPRNG(prng rand.Source) PRNG {
	return PRNG{
		Rand: rand.New(prng),
	}
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

// RollPercentile generates a random float64 value in the range [0.0, 1.0).
// This can be used to represent a random percentile.
// Returns:
//   - A random float64 value between 0.0 (inclusive) and 1.0 (exclusive)
func (p PRNG) RollPercentile() float64 {
	return p.Float64()
}

// Vary5pct returns a value that is randomly varied within 5% (higher or lower) of the input value.
// Parameters:
//   - f: The base value to vary
//
// Returns:
//   - A float64 value that is within ±5% of the input value
func (p PRNG) Vary5pct(f float64) float64 {
	return f * (0.93 + p.RollD6(2)/100.0)
}

// Vary10pct returns a value that is randomly varied within 10% (higher or lower) of the input value.
// Parameters:
//   - f: The base value to vary
//
// Returns:
//   - A float64 value that is within ±10% of the input value
func (p PRNG) Vary10pct(f float64) float64 {
	return f * (0.86 + p.RollD6(4)/100.0)
}
