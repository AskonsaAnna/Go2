package dispatch

// var mu sync.Mutex

// check if any other trains have been scheduled to be at a given station
// on a given turn
func pathFree(path []string, locked map[int]map[string]bool) bool {
	// mu.Lock()
	// defer mu.Unlock()

	for turn, station := range path {
		if locked[turn] == nil {
			locked[turn] = make(map[string]bool)
		}

		if locked[turn][station] {
			return false
		}
	}
	return true
}

// mark which stations the train is scheduled to be at,
// and on which turn
func trainPath(path []string, start, end string, locked map[int]map[string]bool) {
	// mu.Lock()
	// defer mu.Unlock()

	for turn, station := range path {
		if locked[turn] == nil {
			locked[turn] = make(map[string]bool)
		}
		if station != start && station != end {
			locked[turn][station] = true
		}
	}
}
