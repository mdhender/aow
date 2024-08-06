// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package aow provides a simple API for creating star systems using the
// mechanics from the book "Architect of Worlds" by Jon F. Zeigler.
package aow

import (
	"github.com/mdhender/semver"
	"math/rand/v2"
)

var (
	version = semver.Version{Major: 0, Minor: 0, Patch: 1}
)

// Generator is the structure that manages the settings for the generator.
type Generator struct {
	prng *rand.Rand
}

// New returns a new Generator.
//
// The prng parameter is used to generate random numbers.
// For normal use, pass in a source like rand.NewPCG(rand.Uint64(), rand.Uint64()).
// If you're testing and need repeatability, use a source like rand.NewPCG(0xcafe, 0xcafe).
func New(prng rand.Source, options ...Option) (*Generator, error) {
	if prng == nil {
		return nil, ErrPRNGNil
	}
	g := &Generator{
		prng: rand.New(prng),
	}
	for _, option := range options {
		if err := option(g); err != nil {
			return nil, err
		}
	}
	return g, nil
}
