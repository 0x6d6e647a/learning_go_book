package main

import (
	"fmt"
	"math"
	"sync"
)

const limit = 100_000
const parallel = false

// NOTE :: This was just for fun and much slower.
func calcSquareRootsParallel() map[int]float64 {
	squareRoots := make(map[int]float64, limit)

	var mutex sync.RWMutex
	var waitGroup sync.WaitGroup
	waitGroup.Add(limit)

	for i := range limit {
		go func() {
			defer waitGroup.Done()
			mutex.Lock()
			defer mutex.Unlock()
			squareRoots[i] = math.Sqrt(float64(i))
		}()
	}

	waitGroup.Wait()

	return squareRoots
}

func calcSquareRootsSequential() map[int]float64 {
	squareRoots := make(map[int]float64, limit)

	for i := range limit {
		squareRoots[i] = math.Sqrt(float64(i))
	}

	return squareRoots
}

func calcSquareRoots() map[int]float64 {
	if parallel {
		return calcSquareRootsParallel()
	} else {
		return calcSquareRootsSequential()
	}
}

var squareRoots = sync.OnceValue(calcSquareRoots)

func main() {
	squareRoots := squareRoots()

	for i := 0; i < limit; i += 1_000 {
		fmt.Printf("%06d => %0.4f\n", i, squareRoots[i])
	}
}
