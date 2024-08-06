// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

import "math"

type Option func(*Generator) error

// WithOffset allows you to specify an offset (distance from the center of the galaxy),
// which triggers the generator to use the advanced population model.
//
// Parameters:
//   - r: The distance (in parsecs) from the center of the galaxy
//   - h: The distance (in parsecs) above or below the galactic plane
func WithOffset(r, h float64) Option {
	return func(g *Generator) error {
		r, h = math.Abs(r), math.Abs(h)
		if r < 300 {
			return ErrNeighborhoodOffsetTooSmall
		} else if r > 30_000 {
			return ErrNeighborhoodOffsetTooLarge
		} else if h > 1_250 {
			return ErrNeighborhoodOffsetTooLarge
		}
		return nil
	}
}
