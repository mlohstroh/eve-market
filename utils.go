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
