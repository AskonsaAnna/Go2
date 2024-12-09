package pathing

import (
	"sort"
	"trains/mapping"
)

func FindPaths(start, end string, trains int) [][]string {
	graph := mapping.Stations

	var paths [][]string
	var currentPath []string

	// if numTrains == 1, find only the shortest path
	if trains == 1 {
		paths = append(paths, bfs(graph, start, end))
		return paths
	}

	// survive large maps WIP
	// find shortest paths from stations connected to start station
	if len(mapping.Stations) > 100 {

		for _, station := range mapping.Stations[start].Connections {

			path := bfs(graph, station, end)
			path = append([]string{start}, path...)

			paths = append(paths, path)

		}

		// find all paths in small maps, remove if make the above logic find all relevant routes
	} else {
		// dfs finds every possible path in the map, not feasible in large maps
		dfs(graph, start, start, end, &currentPath, &paths)
	}

	// sort the paths
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	return paths
}

// dfs finds all paths from start to end using Depth-First Search.
func dfs(graph map[string]mapping.Station, start, current, end string, currentPath *[]string, paths *[][]string) {
	*currentPath = append(*currentPath, current)

	// if the current station is the end station, add the current path to paths
	if current == end {

		pathCopy := make([]string, len(*currentPath))
		copy(pathCopy, *currentPath)
		*paths = append(*paths, pathCopy)

		// otherwise, explore each connection

	} else {
		for _, neighbor := range graph[current].Connections {

			// to avoid cycles, skip the neighbor if it's already in the current path
			if !contains(*currentPath, neighbor) || (contains(graph[start].Connections, neighbor) && neighbor != (*currentPath)[1]) {
				dfs(graph, start, neighbor, end, currentPath, paths)
			}
		}

	}

	*currentPath = (*currentPath)[:len(*currentPath)-1]
}

// bfs finds the shortest path from start to end using Breadth-First Search.
func bfs(graph map[string]mapping.Station, start, end string) []string {
	queue := [][]string{{start}}
	visited := make(map[string]bool)
	visited[start] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		lastNode := path[len(path)-1]

		if lastNode == end {
			return path
		}

		for _, neighbor := range graph[lastNode].Connections {
			if !visited[neighbor] {
				newPath := append([]string{}, path...)
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
				visited[neighbor] = true
			}
		}
	}

	return nil
}

// contains checks if a slice contains a particular item.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
