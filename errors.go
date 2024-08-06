// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

// Error defines a constant error
type Error string

// Error implements the Errors interface
func (e Error) Error() string { return string(e) }

const (
	ErrPRNGNil = Error("PRNG cannot be nil")
)
