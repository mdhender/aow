// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow_test

import (
	"github.com/mdhender/aow"
	"math/rand/v2"
	"testing"
)

func TestGenerator_RollD6(t *testing.T) {
	g, err := aow.New(0, rand.NewPCG(0xcafe, 0xcafe), aow.ReferenceCatalog) // Use a fixed seed for reproducibility
	if err != nil {
		t.Fatalf("New() = %v, want nil", err)
	}

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
		result := g.RollD6(tc.n)
		if result < tc.min || result > tc.max {
			t.Errorf("RollD6(%d) = %f, want between %f and %f", tc.n, result, tc.min, tc.max)
		}
	}
}

func TestGenerator_RollD10(t *testing.T) {
	g, err := aow.New(0, rand.NewPCG(0xcafe, 0xcafe), aow.ReferenceCatalog) // Use a fixed seed for reproducibility
	if err != nil {
		t.Fatalf("New() = %v, want nil", err)
	}

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
		result := g.RollD10(tc.n)
		if result < tc.min || result > tc.max {
			t.Errorf("RollD10(%d) = %f, want between %f and %f", tc.n, result, tc.min, tc.max)
		}
	}
}

func TestGenerator_RollPercentile(t *testing.T) {
	g, err := aow.New(0, rand.NewPCG(0xcafe, 0xcafe), aow.ReferenceCatalog) // Use a fixed seed for reproducibility
	if err != nil {
		t.Fatalf("New() = %v, want nil", err)
	}

	for i := 0; i < 1_000; i++ {
		result := g.RollPercentile()
		if result < 0 || result >= 1 {
			t.Errorf("RollPercentile() = %f, want between 0 and 1", result)
		}
	}
}

func TestGenerator_Vary5pct(t *testing.T) {
	g, err := aow.New(0, rand.NewPCG(0xcafe, 0xcafe), aow.ReferenceCatalog) // Use a fixed seed for reproducibility
	if err != nil {
		t.Fatalf("New() = %v, want nil", err)
	}

	minp, val, maxp := 0.95*10_000, 10_000.0, 1.05*10_000
	for n := 0; n < 1_000; n++ {
		result := g.Vary5pct(val)
		if result < minp || result > maxp {
			t.Errorf("Vary5pct(%f) = %f, want between %f and %f", val, result, minp, maxp)
		}
	}
}

func TestGenerator_Vary10pct(t *testing.T) {
	g, err := aow.New(0, rand.NewPCG(0xcafe, 0xcafe), aow.ReferenceCatalog) // Use a fixed seed for reproducibility
	if err != nil {
		t.Fatalf("New() = %v, want nil", err)
	}

	minp, val, maxp := 0.90*10_000, 10_000.0, 1.10*10_000
	for n := 0; n < 1_000; n++ {
		result := g.Vary10pct(val)
		if result < minp || result > maxp {
			t.Errorf("Vary10pct(%f) = %f, want between %f and %f", val, result, minp, maxp)
		}
	}
}
