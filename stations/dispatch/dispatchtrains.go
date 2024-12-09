package dispatch

import (
	"fmt"
	"sync"
)

type Train struct {
	Number int
	Path   []string
}

func DistributeTrains(paths [][]string, startStation, endStation string, trains int) {

	// set amount of trains to be scheduled per go routine
	blockSize := 100
	// calculate the amount of blocks
	numBlocks := (trains + blockSize - 1) / blockSize
	turn := 0

	// Create channels for communication between goroutines
	scheduleCh := make(chan []Train, numBlocks)
	blockChs := make([]chan struct{}, numBlocks)

	var wg sync.WaitGroup

	// Schedule movements in blocks
	for i := 0; i < numBlocks; i++ {
		blockChs[i] = make(chan struct{})
		wg.Add(1)

		// schedule train movements for a block
		go func(blockNum int) {
			defer wg.Done()
			start := blockNum * blockSize
			end := start + blockSize
			if end > trains {
				end = trains
			}
			trainsInBlock := end - start

			// Initialize locked map for this block
			locked := make(map[int]map[string]bool)

			// use clean version of paths for each block, without accumulated wait time from other blocks
			cleanPaths := make([][]string, len(paths))
			copy(cleanPaths, paths)

			// collect scheduled movements for the trains in block
			distributed := scheduleTrainMovements(trainsInBlock, blockNum, startStation, endStation, cleanPaths, locked)
			if blockNum > 0 {
				<-blockChs[blockNum-1]
			}

			// save scheduled movements for trains in block to the channel
			scheduleCh <- distributed
			close(blockChs[blockNum])
		}(i)
	}

	go func() {
		wg.Wait()
		close(scheduleCh)
	}()

	// Collect scheduled blocks in order
	scheduledBlocks := make([][]Train, numBlocks)
	for i := range scheduledBlocks {
		scheduledBlocks[i] = <-scheduleCh
	}

	// Display movements for each block in order
	for _, distributed := range scheduledBlocks {
		turn += displayTrainMovements(distributed, startStation)
	}

	fmt.Println()
	fmt.Printf("Moved %d trains from %s to %s in %d turns\n", trains, startStation, endStation, turn)
}

// schedules train movements for trains in a block
func scheduleTrainMovements(trainsRemaining, blockNumber int, startStation, endStation string, paths [][]string, locked map[int]map[string]bool) []Train {
	distributed := []Train{}

	// schedule each train to be sent down the shortest path available
	// if a path is used, or has a scheduling conflict, that path gets wait time added to it
	// and is moved down the order
	for trainsRemaining > 0 {
		shortest := paths[0]

		// check for scheduling conflicts

		if pathFree(shortest, locked) {
			currentTrain := Train{
				Number: len(distributed) + 1 + blockNumber*100,
				Path:   shortest,
			}
			distributed = append(distributed, currentTrain)
			trainPath(shortest, startStation, endStation, locked)
			trainsRemaining--
		}

		// add start station to the beginnig of path to simulate wait time
		shortest = append([]string{shortest[0]}, shortest...)
		paths[0] = shortest

		// move path down the order

		for i := 0; i < len(paths)-1; i++ {
			if len(paths[i]) >= len(paths[i+1]) {
				paths[i], paths[i+1] = paths[i+1], paths[i]
			}
		}

	}

	// return block of scheduled train movements
	return distributed
}

func displayTrainMovements(distributed []Train, startStation string) int {
	turn := -1
	for len(distributed) > 0 {

		// use variable to keep track of trains not already at the end station
		var remainingTrains []Train
		moveMade := false

		for _, train := range distributed {

			// the first station in slice is removed each round, so the first station in slice is the current station
			currentStation := train.Path[0]

			// trains at the start station shouldn't be displayed
			if currentStation != startStation {

				fmt.Printf("T%d-%s ", train.Number, currentStation)

				// change flag value
				moveMade = true
			}

			// if slice length is greater than 1, the train is not at the end station
			if len(train.Path) > 1 {
				train.Path = train.Path[1:]
				remainingTrains = append(remainingTrains, train)
			}
		}

		if moveMade {
			fmt.Printf("\n")
		}

		distributed = remainingTrains
		turn++
	}

	// to add rounds spent for train movements in block to total
	return turn
}
