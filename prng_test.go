// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow_test

import (
	"github.com/mdhender/aow"
	"math/rand/v2"
	"testing"
)

func TestPRNG_RollD6(t *testing.T) {
	p := aow.NewPRNG(rand.NewPCG(0xcafe, 0xcafe)) // Use a fixed seed for reproducibility
	for _, tc := range []struct {
		name string
		n    int
		min  float64
		max  float64
	}{
		{"Single die", 1, 1, 6},
		{"Three dice", 3, 3, 18},
		{"Ten dice", 10, 10, 60},
	} {
		result := p.RollD6(tc.n)
		if result < tc.min || result > tc.max {
			t.Errorf("RollD6(%d) = %f, want between %f and %f", tc.n, result, tc.min, tc.max)
		}
	}
}

func TestPRNG_RollD10(t *testing.T) {
	p := aow.NewPRNG(rand.NewPCG(0xcafe, 0xcafe)) // Use a fixed seed for reproducibility
	for _, tc := range []struct {
		name string
		n    int
		min  float64
		max  float64
	}{
		{"Single die", 1, 1, 10},
		{"Three dice", 3, 3, 30},
		{"Ten dice", 10, 10, 100},
	} {
		result := p.RollD10(tc.n)
		if result < tc.min || result > tc.max {
			t.Errorf("RollD10(%d) = %f, want between %f and %f", tc.n, result, tc.min, tc.max)
		}
	}
}

func TestPRNG_RollPercentile(t *testing.T) {
	p := aow.NewPRNG(rand.NewPCG(0xcafe, 0xcafe)) // Use a fixed seed for reproducibility
	for i := 0; i < 1_000; i++ {
		result := p.RollPercentile()
		if result < 0 || result >= 1 {
			t.Errorf("RollPercentile() = %f, want between 0 and 1", result)
		}
	}
}

func TestPRNG_Vary5pct(t *testing.T) {
	p := aow.NewPRNG(rand.NewPCG(0xcafe, 0xcafe)) // Use a fixed seed for reproducibility
	minp, val, maxp := 0.95*10_000, 10_000.0, 1.05*10_000
	for n := 0; n < 1_000; n++ {
		result := p.Vary5Pct(val)
		if result < minp || result > maxp {
			t.Errorf("Vary5Pct(%f) = %f, want between %f and %f", val, result, minp, maxp)
		}
	}
}

func TestPRNG_Vary10pct(t *testing.T) {
	p := aow.NewPRNG(rand.NewPCG(0xcafe, 0xcafe)) // Use a fixed seed for reproducibility
	minp, val, maxp := 0.90*10_000, 10_000.0, 1.10*10_000
	for n := 0; n < 1_000; n++ {
		result := p.Vary10Pct(val)
		if result < minp || result > maxp {
			t.Errorf("Vary10Pct(%f) = %f, want between %f and %f", val, result, minp, maxp)
		}
	}
}
