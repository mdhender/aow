// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow_test

import (
	"github.com/mdhender/aow"
	"math"
	"testing"
)

func TestPopulationModelForEarthLikeSystems(t *testing.T) {
	for _, tc := range []struct {
		name     string
		n        int
		tweak    float64
		expected float64
	}{
		{"bob", 40, 0, 12_000},
	} {
		result := aow.PopulationModelForEarthLikeSystems(tc.n, tc.tweak)
		if math.Abs(result.Volume-tc.expected) > 0.001 {
			t.Errorf("PopulationModelForEarthLikeSystems(%d, %f) = %f, want %f", tc.n, tc.tweak, result, tc.expected)
		}
	}
}

func TestPopulationModelForSolLikeNeighborhood(t *testing.T) {
	for _, tc := range []struct {
		name     string
		n        int
		tweak    float64
		expected float64
	}{
		{"alice", 100, 1, 1_300},
	} {
		result := aow.PopulationModelForSolLikeNeighborhood(tc.n, tc.tweak)
		if math.Abs(result.Volume-tc.expected) > 0.001 {
			t.Errorf("PopulationModelForSolLikeNeighborhood(%d, %f) = %f, want %f", tc.n, tc.tweak, result, tc.expected)
		}
	}
}

func TestPopulationModelForOtherNeighborhoods(t *testing.T) {
	for _, tc := range []struct {
		name     string
		n        int
		tweak    float64
		expected float64
	}{
		{"sol-like", 100, 0, 1_235.378},
	} {
		result := aow.PopulationModelForOtherNeighborhoods(tc.n, 8_000, 20, tc.tweak)
		if math.Abs(result.Volume-tc.expected) > 0.001 {
			t.Errorf("PopulationModelForOtherNeighborhoods(%d, %f) = %f, want %f", tc.n, tc.tweak, result, tc.expected)
		}
	}
}
