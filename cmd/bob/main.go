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
	if err := run(false); err != nil {
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

	err = g.BackgroundPopulation()
	if err != nil {
		return err
	}
	log.Printf("g: %d star systems, %f radius", len(g.Catalog), g.Radius)
	if addCluster {
		x, y, z := g.GenXYZ(g.Radius*2/3, g.Radius)
		err = g.OpenCluster(x, y, z)
		if err != nil {
			return err
		}
	}
	g.SortCatalog()
	for n, ss := range g.Catalog {
		log.Printf("%4d: ss %+v", n+1, *ss)
	}

	return nil
}
