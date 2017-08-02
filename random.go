package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
 * Essentially a wrapper for rand.Intn which includes a seed
 * Something resembling a macro would be applicable here,
 * But I have no idea if such a thing exists
 *
 * returns an int within the range [0, max]
 */
func randInt(max int) (int, error) {
	if max < 0 {
		return 0, fmt.Errorf("rand: max input (%d) is negative!", max)
	}
	if max == 0 {
		return 0, nil
	}

	rand.Seed(time.Now().Unix())
	return rand.Intn(max), nil
}

/*
 * Finds a random int in range of [min, max]
 * Works with negatives as well
 */
func randRangeInt(min, max int) (random int, err error) {
	if max < min {
		return 0, fmt.Errorf("math: invalid range, %d < %d", max, min)
	}
	if max == min {
		random = max
		return
	}

	random, err = randInt(max - min + 1)
	if err != nil {
		return 0, err
	}
	random += min
	return
}
