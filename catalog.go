// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

import (
	"log"
	"math"
	"sort"
)

type Catalog_t struct {
	Kind        Catalog_e
	Radius      float64     // in parsecs
	Coordinates Coordinates // relative to an arbitrary point
	StarSystems []*StarSystem_t
}

type Catalog_e int

const (
	// SurveyCatalog is the default catalog and that includes all systems.
	SurveyCatalog Catalog_e = iota

	// ReferenceCatalog maps only "interesting" systems.
	ReferenceCatalog
)

// NewBackgroundPopulation creates a catalog containing the background population of a neighborhood.
//
// Uses the population model to generate the initial set of star systems.
func NewBackgroundPopulation(pm PopulationModel_t, prng PRNG) (*Catalog_t, error) {
	var c Catalog_t

	for _, v := range []struct {
		key   StellarPopulation_e
		value populationModel_t
	}{
		{key: YoungPopulationI, value: pm.YoungPopulationI},
		{key: IntermediatePopulationI, value: pm.IntermediatePopulationI},
		{key: OldPopulationI, value: pm.OldPopulationI},
		{key: DiskPopulationII, value: pm.DiskPopulationII},
		{key: HaloPopulationII, value: pm.HaloPopulationII},
	} {
		numberOfStarSystems := int(math.Ceil(prng.Vary10Pct(v.value.Density * pm.Volume)))
		for i := 0; i < numberOfStarSystems; i++ {
			c.StarSystems = append(c.StarSystems, &StarSystem_t{
				Population: v.key,
				// generate a random age for the star system
				Age: v.value.BaseAge + v.value.AgeRange*prng.RollPercentile(),
				// generate a random position for the star system
				Coordinates: prng.GenXYZ().Scale(pm.Radius),
			})
		}
	}

	return &c, nil
}

const (
	minPctClusterCoreZone  float64 = 0.0
	maxPctClusterCoreZone  float64 = 0.05
	minPctTidalRadiusZone  float64 = 0.05
	maxPctTidalRadiusZone  float64 = 0.2
	minPctExtendedHaloZone float64 = 0.2
	maxPctExtendedHaloZone float64 = 1.0
)

func NewOpenCluster(prng PRNG) (*Catalog_t, error) {
	// cluster can be tightly or loosely bound.
	isTightlyBound := prng.RollD6(3) <= 5

	// determine the age of the cluster (in billions of years)
	var clusterAge float64
	if isTightlyBound {
		switch n := prng.RollD100(); {
		case n <= 2:
			clusterAge = 0.0 + 0.1*prng.RollPercentile()
		case n <= 4:
			clusterAge = 0.1 + 0.1*prng.RollPercentile()
		case n <= 6:
			clusterAge = 0.2 + 0.1*prng.RollPercentile()
		case n <= 8:
			clusterAge = 0.3 + 0.1*prng.RollPercentile()
		case n <= 10:
			clusterAge = 0.4 + 0.1*prng.RollPercentile()
		case n <= 12:
			clusterAge = 0.5 + 0.1*prng.RollPercentile()
		case n <= 14:
			clusterAge = 0.6 + 0.1*prng.RollPercentile()
		case n <= 16:
			clusterAge = 0.7 + 0.1*prng.RollPercentile()
		case n <= 18:
			clusterAge = 0.8 + 0.1*prng.RollPercentile()
		case n <= 20:
			clusterAge = 0.9 + 0.1*prng.RollPercentile()
		case n <= 45:
			clusterAge = 1.0 + 2.0*prng.RollPercentile()
		default:
			clusterAge = 3.0 + 5.0*prng.RollPercentile()
		}
	} else {
		switch n := prng.RollD100(); {
		case n <= 21:
			clusterAge = 0.0 + 0.1*prng.RollPercentile()
		case n <= 38:
			clusterAge = 0.1 + 0.1*prng.RollPercentile()
		case n <= 52:
			clusterAge = 0.2 + 0.1*prng.RollPercentile()
		case n <= 64:
			clusterAge = 0.3 + 0.1*prng.RollPercentile()
		case n <= 73:
			clusterAge = 0.4 + 0.1*prng.RollPercentile()
		case n <= 81:
			clusterAge = 0.5 + 0.1*prng.RollPercentile()
		case n <= 87:
			clusterAge = 0.6 + 0.1*prng.RollPercentile()
		case n <= 92:
			clusterAge = 0.7 + 0.1*prng.RollPercentile()
		case n <= 96:
			clusterAge = 0.8 + 0.1*prng.RollPercentile()
		default:
			clusterAge = 0.9 + 0.1*prng.RollPercentile()
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

	// generate the initial radius (in parsecs), give or take 0.25 parsecs
	clusterRadius := prng.RollD6(2) / 2
	log.Printf("cluster radius: %f\n", clusterRadius)
	clusterRadius += (prng.VaryNPct(1.0, 0.25) - 1)
	log.Printf("cluster radius: %f\n", clusterRadius)

	// initial number of star systems in the cluster
	numberOfStarSystems := prng.RollD6(2) / 2
	if isTightlyBound && numberOfStarSystems < 3.5 {
		numberOfStarSystems = 3.5
	}
	numberOfStarSystems = math.Floor(numberOfStarSystems * math.Pow(clusterRadius, 3))
	log.Printf("number of star systems: %f\n", numberOfStarSystems)

	// determine the effective age of the cluster for the evaporation table
	effectiveClusterAge := clusterAge
	if isTightlyBound {
		effectiveClusterAge = clusterAge / 10
	}
	log.Printf("effective cluster age: %f\n", effectiveClusterAge)

	// use the cluster evaporation table to get radius (as a percentage) for each zone
	var corePct, tidalPct, extendedHaloPct float64
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
	log.Printf("core: %f, tidal: %f, extended halo: %f\n", corePct, tidalPct, extendedHaloPct)
	coreCount, tidalCount, extendedHaloCount := int(corePct*numberOfStarSystems), int(tidalPct*numberOfStarSystems), int(extendedHaloPct*numberOfStarSystems)
	log.Printf("core count: %d, title count: %d, extended halo count: %d\n", coreCount, tidalCount, extendedHaloCount)
	if float64(coreCount+tidalCount+extendedHaloCount) < numberOfStarSystems {
		coreCount++
		if float64(coreCount+tidalCount+extendedHaloCount) < numberOfStarSystems {
			tidalCount++
		}
	}
	log.Printf("core count: %d, title count: %d, extended halo count: %d\n", coreCount, tidalCount, extendedHaloCount)

	// we have the information needed to create the catalog for the cluster
	var catalog Catalog_t

	// create star systems in the cluster core zone
	log.Printf("gen core %f %f %d/%d\n", minPctClusterCoreZone, maxPctClusterCoreZone, int(corePct*numberOfStarSystems), int(numberOfStarSystems))
	for ; coreCount > 0; coreCount-- {
		catalog.StarSystems = append(catalog.StarSystems, &StarSystem_t{
			Population: stpop,
			// generate a random age for the star system
			Age: prng.Vary5Pct(clusterAge),
			// generate a random position for the star system
			Coordinates: prng.GenZonedXYZ(minPctClusterCoreZone, maxPctClusterCoreZone).Scale(clusterRadius),
		})
	}

	// create star systems in the tidal radius zone
	log.Printf("gen tidal %f %f %d/%d\n", minPctTidalRadiusZone, maxPctTidalRadiusZone, int(tidalPct*numberOfStarSystems), int(numberOfStarSystems))
	for ; tidalCount > 0; tidalCount-- {
		catalog.StarSystems = append(catalog.StarSystems, &StarSystem_t{
			Population: stpop,
			// generate a random age for the star system
			Age: prng.Vary5Pct(clusterAge),
			// generate a random position for the star system
			Coordinates: prng.GenZonedXYZ(minPctTidalRadiusZone, maxPctTidalRadiusZone).Scale(clusterRadius),
		})
	}

	// create star systems in the extended halo zone
	log.Printf("gen halo %f %f %d/%d\n", minPctExtendedHaloZone, maxPctExtendedHaloZone, int(extendedHaloPct*numberOfStarSystems), int(numberOfStarSystems))
	for ; extendedHaloCount > 0; extendedHaloCount-- {
		catalog.StarSystems = append(catalog.StarSystems, &StarSystem_t{
			Population: stpop,
			// generate a random age for the star system
			Age: prng.Vary5Pct(clusterAge),
			// generate a random position for the star system
			Coordinates: prng.GenZonedXYZ(minPctExtendedHaloZone, maxPctExtendedHaloZone).Scale(clusterRadius),
		})
	}

	return &catalog, nil
}

func NewStellarAssociation(prng PRNG) (*Catalog_t, error) {
	return NewOpenCluster(prng)
}

func (c *Catalog_t) Length() int {
	return len(c.StarSystems)
}

func (c *Catalog_t) Sort() {
	sort.Slice(c.StarSystems, func(i, j int) bool {
		if c.StarSystems[i].Age == c.StarSystems[j].Age {
			return c.StarSystems[i].Population < c.StarSystems[j].Population
		}
		return c.StarSystems[i].Age < c.StarSystems[j].Age
	})
}

func (c *Catalog_t) SortByDistance(origin Coordinates) {
	for _, ss := range c.StarSystems {
		ss.distance = ss.Coordinates.DistanceTo(origin)
	}
	sort.Slice(c.StarSystems, func(i, j int) bool {
		return c.StarSystems[i].distance < c.StarSystems[j].distance
	})
}

func (c *Catalog_t) Merge(other *Catalog_t, offset Coordinates) {
	for _, ss := range other.StarSystems {
		c.StarSystems = append(c.StarSystems, &StarSystem_t{
			Population:  ss.Population,
			Age:         ss.Age,
			Coordinates: ss.Coordinates.Translate(offset),
		})
	}
}
