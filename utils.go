package main

import (
	"log"
	"time"
)

func profileFunction(name string, f func()) {
	startTime := time.Now()

	f()

	elapsedDuration := time.Now().Sub(startTime)
	log.Printf("Function [%s] took %s to run", name, elapsedDuration.String())
}

// ContainsI checks if the int value is contained within an int array
func ContainsI(val int, a []int) bool {
	for _, v := range a {
		if val == v {
			return true
		}
	}

	return false
}
