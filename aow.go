// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package aow provides a simple API for creating star systems using the
// mechanics from the book "Architect of Worlds" by Jon F. Zeigler.
package aow

import (
	"fmt"
	"github.com/mdhender/semver"
	"log"
	"math"
	"math/rand/v2"
)

var (
	version = semver.Version{Major: 0, Minor: 0, Patch: 6}
)

// Generator is the structure that manages the settings for the generator.
type Generator struct {
	prng          PRNG
	typeOfCatalog Catalog_e // the type of catalog used to generate the star systems
	pm            PopulationModel_t
	Radius        float64 // the radius of the map in parsecs

	Catalog *Catalog_t
}

// New returns a new Generator that will generate a catalog with approximately
// the given number of stellar systems. You may provide additional options to
// modify the behavior of the generator.
//
// The n parameter is the target number of systems with Earth-like planets to generate.
//
// The prng parameter is used to generate random numbers.
// For normal use, pass in a source like rand.NewPCG(rand.Uint64(), rand.Uint64()).
// If you're testing and need repeatability, use a source like rand.NewPCG(0xcafe, 0xcafe).
func New(n int, prng rand.Source, cat Catalog_e, options ...Option) (*Generator, error) {
	if prng == nil {
		return nil, ErrPRNGNil
	}
	g := &Generator{
		prng: PRNG{Rand: rand.New(prng)},
	}
	for _, option := range options {
		if err := option(g); err != nil {
			return nil, err
		}
	}

	// this iteration of the generator uses the basic population model and
	// treats "n" as the number of Earth-like planets to target.
	g.pm = PopulationModelForSolLikeNeighborhood(n, 0)
	g.Radius = math.Ceil(math.Cbrt((3 * g.pm.Volume) / (4 * math.Pi)))

	return g, nil
}

// BackgroundPopulation creates the background population of the catalog.
func (g *Generator) BackgroundPopulation() error {
	log.Printf("pm %+v\n", g.pm)
	catalog, err := NewBackgroundPopulation(g.pm, g.prng)
	if err != nil {
		return err
	}
	g.Catalog = catalog
	return nil
}

// OpenCluster creates a new open cluster.
func (g *Generator) OpenCluster(origin Coordinates) (*Catalog_t, error) {
	catalog, err := NewOpenCluster(g.prng)
	if err != nil {
		return nil, err
	}
	return catalog, nil
}

// StellarAssociation creates a new stellar association.
func (g *Generator) StellarAssociation(origin Coordinates) (*Catalog_t, error) {
	return g.OpenCluster(origin)
}

// SortCatalog sorts the catalog by age and population.
func (g *Generator) SortCatalog() {
	g.Catalog.Sort()
}

// GenXYZ returns scaled coordinates (they've been adjusted for the radius)
func (g *Generator) GenXYZ() Coordinates {
	return g.prng.GenXYZ().Scale(g.Radius)
}

// GenZonedXYZ returns scaled coordinates (they've been adjusted for the radius)
func (g *Generator) GenZonedXYZ(minPct, maxPct float64) Coordinates {
	return g.prng.GenZonedXYZ(minPct, maxPct).Scale(g.Radius)
}

type Coordinates struct {
	X, Y, Z float64
}

func (c Coordinates) DistanceBetween(o Coordinates) float64 {
	return c.DistanceTo(o)
}

func (c Coordinates) DistanceTo(o Coordinates) float64 {
	dx, dy, dz := o.X-c.X, o.Y-c.Y, o.Z-c.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (c Coordinates) Scale(s float64) Coordinates {
	return Coordinates{
		X: c.X * s,
		Y: c.Y * s,
		Z: c.Z * s,
	}
}

func (c Coordinates) Translate(s Coordinates) Coordinates {
	return Coordinates{
		X: c.X + s.X,
		Y: c.Y + s.Y,
		Z: c.Z + s.Z,
	}
}

func (c Coordinates) String() string {
	return fmt.Sprintf("(%.2g %.2g %.2g)", c.X, c.Y, c.Z)
}
