// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package aow provides a simple API for creating star systems using the
// mechanics from the book "Architect of Worlds" by Jon F. Zeigler.
package aow

import (
	"github.com/mdhender/semver"
	"log"
	"math"
	"math/rand/v2"
	"sort"
)

var (
	version = semver.Version{Major: 0, Minor: 0, Patch: 5}
)

// Generator is the structure that manages the settings for the generator.
type Generator struct {
	prng          PRNG
	typeOfCatalog Catalog_e // the type of catalog used to generate the star systems
	pm            PopulationModel_t
	Radius        float64 // the radius of the map in parsecs

	Catalog []*StarSystem_t
}

type Catalog_e int

const (
	// ReferenceCatalog is the default catalog and maps only "interesting" systems.
	ReferenceCatalog Catalog_e = iota

	// SurveyCatalog is a catalog that includes all systems.
	SurveyCatalog
)

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
	for _, v := range []struct {
		key   StellarPopulation_e
		value populationModel_t
	}{
		{key: YoungPopulationI, value: g.pm.YoungPopulationI},
		{key: IntermediatePopulationI, value: g.pm.IntermediatePopulationI},
		{key: OldPopulationI, value: g.pm.OldPopulationI},
		{key: DiskPopulationII, value: g.pm.DiskPopulationII},
		{key: HaloPopulationII, value: g.pm.HaloPopulationII},
	} {
		numberOfStarSystems := int(math.Ceil(g.Vary10pct(v.value.Density * g.pm.Volume)))
		for i := 0; i < numberOfStarSystems; i++ {
			var ss StarSystem_t
			ss.Population = v.key
			// generate a random age for the star system
			ss.Age = v.value.BaseAge + v.value.AgeRange*g.RollPercentile()
			// generate a random position for the star system
			ss.X, ss.Y, ss.Z = g.prng.GenXYZ()
			ss.X, ss.Y, ss.Z = ss.X*g.Radius, ss.Y*g.Radius, ss.Z*g.Radius
			g.Catalog = append(g.Catalog, &ss)
		}
	}
	return nil
}

func (g *Generator) OpenCluster(x, y, z float64) error {
	isTightlyBound := g.RollD6(3) <= 5

	var clusterAge float64
	if isTightlyBound {
		switch n := g.RollD100(); {
		case n <= 2:
			clusterAge = 0.0 + 0.1*g.RollPercentile()
		case n <= 4:
			clusterAge = 0.1 + 0.1*g.RollPercentile()
		case n <= 6:
			clusterAge = 0.2 + 0.1*g.RollPercentile()
		case n <= 8:
			clusterAge = 0.3 + 0.1*g.RollPercentile()
		case n <= 10:
			clusterAge = 0.4 + 0.1*g.RollPercentile()
		case n <= 12:
			clusterAge = 0.5 + 0.1*g.RollPercentile()
		case n <= 14:
			clusterAge = 0.6 + 0.1*g.RollPercentile()
		case n <= 16:
			clusterAge = 0.7 + 0.1*g.RollPercentile()
		case n <= 18:
			clusterAge = 0.8 + 0.1*g.RollPercentile()
		case n <= 20:
			clusterAge = 0.9 + 0.1*g.RollPercentile()
		case n <= 45:
			clusterAge = 1.0 + 2.0*g.RollPercentile()
		default:
			clusterAge = 3.0 + 5.0*g.RollPercentile()
		}
	} else {
		switch n := g.RollD100(); {
		case n <= 21:
			clusterAge = 0.0 + 0.1*g.RollPercentile()
		case n <= 38:
			clusterAge = 0.1 + 0.1*g.RollPercentile()
		case n <= 52:
			clusterAge = 0.2 + 0.1*g.RollPercentile()
		case n <= 64:
			clusterAge = 0.3 + 0.1*g.RollPercentile()
		case n <= 73:
			clusterAge = 0.4 + 0.1*g.RollPercentile()
		case n <= 81:
			clusterAge = 0.5 + 0.1*g.RollPercentile()
		case n <= 87:
			clusterAge = 0.6 + 0.1*g.RollPercentile()
		case n <= 92:
			clusterAge = 0.7 + 0.1*g.RollPercentile()
		case n <= 96:
			clusterAge = 0.8 + 0.1*g.RollPercentile()
		default:
			clusterAge = 0.9 + 0.1*g.RollPercentile()
		}
	}

	// population group depends on the age of the cluster
	var stpop StellarPopulation_e
	if clusterAge < 2.0 {
		stpop = YoungPopulationI
	} else if clusterAge < 5.0 {
		stpop = IntermediatePopulationI
	} else {
		stpop = OldPopulationI
	}

	// generate the initial radius with a random variation between -0.25 and +0.25
	initialRadius := float64(g.RollD6(2))/2 + ((g.prng.RollPercentile() * 0.5) - 0.25)

	// initial number of star systems in the cluster
	numberOfStarSystems := g.RollD6(2)
	if isTightlyBound && numberOfStarSystems < 7 {
		numberOfStarSystems = 7
	}
	numberOfStarSystems /= 2
	numberOfStarSystems = math.Floor(numberOfStarSystems * math.Pow(initialRadius, 3))

	// cluster evaporation table
	var corePct, tidalPct, extendedHaloPct float64
	effectiveClusterAge := clusterAge
	if isTightlyBound {
		effectiveClusterAge = clusterAge / 10
	}
	switch {
	case effectiveClusterAge < 0.1:
		corePct, tidalPct, extendedHaloPct = 1.00, 0.00, 0.00
	case effectiveClusterAge < 0.2:
		corePct, tidalPct, extendedHaloPct = 0.80, 0.20, 0.00
	case effectiveClusterAge < 0.3:
		corePct, tidalPct, extendedHaloPct = 0.64, 0.32, 0.04
	case effectiveClusterAge < 0.4:
		corePct, tidalPct, extendedHaloPct = 0.51, 0.38, 0.10
	case effectiveClusterAge < 0.5:
		corePct, tidalPct, extendedHaloPct = 0.41, 0.41, 0.15
	case effectiveClusterAge < 0.6:
		corePct, tidalPct, extendedHaloPct = 0.33, 0.41, 0.20
	case effectiveClusterAge < 0.7:
		corePct, tidalPct, extendedHaloPct = 0.26, 0.39, 0.25
	case effectiveClusterAge < 0.8:
		corePct, tidalPct, extendedHaloPct = 0.21, 0.37, 0.28
	case effectiveClusterAge < 0.9:
		corePct, tidalPct, extendedHaloPct = 0.17, 0.33, 0.29
	case effectiveClusterAge < 1.0:
		corePct, tidalPct, extendedHaloPct = 0.13, 0.30, 0.30
	default:
		corePct, tidalPct, extendedHaloPct = 0.11, 0.27, 0.30
	}

	// star systems in each zone
	// cluster core zone
	minRadius, maxRadius := 0.0, initialRadius
	log.Printf("gen core %f %f %d/%d\n", minRadius, maxRadius, int(corePct*numberOfStarSystems), int(numberOfStarSystems))
	for n := 0; n < int(corePct*numberOfStarSystems); n++ {
		var ss StarSystem_t
		ss.Population = stpop

		// generate a random age for the star system
		ss.Age = g.Vary5pct(clusterAge)

		// generate a random position for the star system
		cx, cy, cz := g.GenXYZ(minRadius, maxRadius)
		cx, cy, cz = cx+x, cy+y, cz+z

		// convert to the catalog's coordinate system
		ss.X, ss.Y, ss.Z = cx+x, cy+y, cz+z

		g.Catalog = append(g.Catalog, &ss)
	}
	// tidal radius zone
	minRadius, maxRadius = maxRadius, 4*initialRadius
	log.Printf("gen tidal %f %f %d/%d\n", minRadius, maxRadius, int(tidalPct*numberOfStarSystems), int(numberOfStarSystems))
	for n := 0; n < int(tidalPct*numberOfStarSystems); n++ {
		var ss StarSystem_t
		ss.Population = stpop

		// generate a random age for the star system
		ss.Age = g.Vary5pct(clusterAge)

		// generate a random position for the star system
		cx, cy, cz := g.GenXYZ(minRadius, maxRadius)
		cx, cy, cz = cx+x, cy+y, cz+z

		// convert to the catalog's coordinate system
		ss.X, ss.Y, ss.Z = cx+x, cy+y, cz+z

		g.Catalog = append(g.Catalog, &ss)
	}
	// extended halo zone
	minRadius, maxRadius = maxRadius, 20*initialRadius
	log.Printf("gen halo %f %f %d/%d\n", minRadius, maxRadius, int(extendedHaloPct*numberOfStarSystems), int(numberOfStarSystems))
	for n := 0; n < int(extendedHaloPct*numberOfStarSystems); n++ {
		var ss StarSystem_t
		ss.Population = stpop

		// generate a random age for the star system
		ss.Age = g.Vary5pct(clusterAge)

		// generate a random position for the star system
		cx, cy, cz := g.GenXYZ(minRadius, maxRadius)
		cx, cy, cz = cx+x, cy+y, cz+z

		// convert to the catalog's coordinate system
		ss.X, ss.Y, ss.Z = cx+x, cy+y, cz+z

		g.Catalog = append(g.Catalog, &ss)
	}

	return nil
}

func (g *Generator) StellarAssociations(x, y, z float64) error {
	return g.OpenCluster(x, y, z)
}

// SortCatalog sorts the catalog by age and population.
func (g *Generator) SortCatalog() {
	sort.Slice(g.Catalog, func(i, j int) bool {
		if g.Catalog[i].Age == g.Catalog[j].Age {
			return g.Catalog[i].Population < g.Catalog[j].Population
		}
		return g.Catalog[i].Age < g.Catalog[j].Age
	})
}
