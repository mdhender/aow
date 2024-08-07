// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package aow

type StarSystem_t struct {
	Population StellarPopulation_e
	Age        float64 // in billions of years?
	X, Y, Z    float64 // relative to center of the catalog
}
