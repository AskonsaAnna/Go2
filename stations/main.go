package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"trains/dispatch"
	"trains/mapping"
	pathing "trains/paths"
)

func main() {
	startTime := time.Now()

	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "Error: Incorrect number of command line arguments.")
		os.Exit(1)
	}

	mapFile := os.Args[1]
	startStation := os.Args[2]
	endStation := os.Args[3]
	numTrains, err := strconv.Atoi(os.Args[4])
	if err != nil || numTrains <= 0 {
		fmt.Fprintln(os.Stderr, "Error: Number of trains is not a valid positive integer.")
		os.Exit(1)
	}

	// create map
	err = mapping.ParseNetworkMap(mapFile)
	if err != nil {
		os.Exit(1)
	}

	if _, exists := mapping.Stations[startStation]; !exists {
		fmt.Fprintln(os.Stderr, "Error: Start station does not exist.")
		os.Exit(1)
	}

	if _, exists := mapping.Stations[endStation]; !exists {
		fmt.Fprintln(os.Stderr, "Error: End station does not exist.")
		os.Exit(1)
	}

	if startStation == endStation {
		fmt.Fprintln(os.Stderr, "Error: Start and end station are the same.")
		os.Exit(1)
	}

	// find all available paths between start and end stations
	paths := pathing.FindPaths(startStation, endStation, numTrains)

	if paths == nil {
		fmt.Fprintln(os.Stderr, "Error:  No path exists between the start and end stations.")
		os.Exit(1)

	}

	// calculate paths for trains, and display train movements
	dispatch.DistributeTrains(paths, startStation, endStation, numTrains)

	elapsed := time.Since(startTime)

	fmt.Println()
	log.Println("Program execution time: ", elapsed.Abs())
}

