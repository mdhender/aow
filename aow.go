// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package aow provides a simple API for creating star systems using the
// mechanics from the book "Architect of Worlds" by Jon F. Zeigler.
package aow

import (
	"github.com/mdhender/semver"
	"math/rand/v2"
)

var (
	version = semver.Version{Major: 0, Minor: 0, Patch: 2}
)

// Generator is the structure that manages the settings for the generator.
type Generator struct {
	prng    *rand.Rand
	catalog Catalog_e // the type of catalog used to generate the star systems
	pm      PopulationModel
}

type Catalog_e int

const (
	// ReferenceCatalog is the default catalog and maps only "interesting" systems.
	ReferenceCatalog Catalog_e = iota

	// SurveyCatalog is a catalog that includes all systems.
	SurveyCatalog
)

// New returns a new Generator that will generate a region with approximately the given number of stellar systems.
// You may provide additional options to modify the behavior of the generator.
//
// The n parameter is the target number of stellar systems to generate.
//
// The prng parameter is used to generate random numbers.
// For normal use, pass in a source like rand.NewPCG(rand.Uint64(), rand.Uint64()).
// If you're testing and need repeatability, use a source like rand.NewPCG(0xcafe, 0xcafe).
func New(n int, prng rand.Source, cat Catalog_e, options ...Option) (*Generator, error) {
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
	// the number of interesting systems drives the volume of space that we will catalog

	return g, nil
}
