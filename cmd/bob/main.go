// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package main implements the "Bob" example from the book.
package main

import (
	"github.com/mdhender/aow"
	"log"
	"math/rand/v2"
	"time"
)

func main() {
	started := time.Now()
	if err := run(true); err != nil {
		log.Fatal(err)
	}
	log.Printf("Bob took %s", time.Since(started))
}

func run(addCluster bool) error {
	// create a generator for Bob's map.
	// Bob wants at least 40 Earth-like systems.
	g, err := aow.New(1000, aow.NewPRNG(rand.NewPCG(0xcafe, 0xcafe)), aow.ReferenceCatalog)
	if err != nil {
		return err
	}
	log.Printf("g: %+v", g)

	var origin aow.Coordinates

	err = g.BackgroundPopulation()
	if err != nil {
		return err
	}
	log.Printf("g: %d star systems, %f radius", g.Catalog.Length(), g.Catalog.Radius)
	if addCluster {
		clusterCenter := g.GenZonedXYZ(0.77, 0.89)
		log.Printf("x %g y %g z %g r %g %8.4f\n", clusterCenter.X, clusterCenter.Y, clusterCenter.Z, g.Radius, origin.DistanceTo(clusterCenter))
		clusterCatalog, err := g.OpenCluster(clusterCenter)
		if err != nil {
			return err
		}
		g.Catalog.Merge(clusterCatalog, clusterCenter)
	}
	g.Catalog.SortByDistance(origin)
	for n, ss := range g.Catalog.StarSystems {
		log.Printf("%4d: ss pop %v age %6.2f delta %8.4f %s\n", n+1, ss.Population, ss.Age, origin.DistanceTo(ss.Coordinates), ss.Coordinates)
	}

	return nil
}
