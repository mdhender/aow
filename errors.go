// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

// Error defines a constant error
type Error string

// Error implements the Errors interface
func (e Error) Error() string { return string(e) }

const (
	ErrNotImplemented             = Error("not implemented")
	ErrNeighborhoodOffsetTooSmall = Error("galactic neighborhood offset too small")
	ErrNeighborhoodOffsetTooLarge = Error("galactic neighborhood offset too large")
	ErrPRNGNil                    = Error("PRNG cannot be nil")
)
