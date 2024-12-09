package generating

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const numStations = 10000
const maxCoord = 100
const numConnections = 20000

func Generate() {
	// Create a new random source and a new rand.Rand instance
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// Create a file to write the stations and connections
	file, err := os.Create("stations_and_connections.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Generate station names
	stations := make([]string, numStations)
	for i := 0; i < numStations; i++ {
		stations[i] = fmt.Sprintf("station%d", i+1)
	}

	// Create a set to track used coordinates
	usedCoords := make(map[string]bool)

	// Write stations to file
	fmt.Fprintln(file, "stations:")
	for _, station := range stations {
		var x, y int
		var coord string
		for {
			x = r.Intn(maxCoord)
			y = r.Intn(maxCoord)
			coord = fmt.Sprintf("%d,%d", x, y)
			if !usedCoords[coord] {
				usedCoords[coord] = true
				break
			}
		}
		fmt.Fprintf(file, "%s,%d,%d\n", station, x, y)
	}

	// Write connections to file
	fmt.Fprintln(file, "\nconnections:")
	connections := make(map[string]bool)
	for i := 0; i < numConnections; i++ {
		station1 := stations[r.Intn(numStations)]
		station2 := stations[r.Intn(numStations)]

		// Ensure no self-connections and no duplicate connections
		if station1 != station2 {
			connection := fmt.Sprintf("%s-%s", station1, station2)
			reverseConnection := fmt.Sprintf("%s-%s", station2, station1)
			if !connections[connection] && !connections[reverseConnection] {
				fmt.Fprintf(file, "%s\n", connection)
				connections[connection] = true
			}
		}
	}
}
