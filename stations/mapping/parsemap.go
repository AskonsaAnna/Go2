package mapping

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Station struct {
	Coordinates []int    // x, y
	Connections []string // connected stations
}

var Stations map[string]Station

// 4. It displays "Error" on stderr when the map does not contain a "connections:" section.

//12. It displays "Error" on stderr when the map does not contain a "stations:" section.

// 19. It displays "Error" on stderr when station names are invalid.

func ParseNetworkMap(mapFile string) error {
	stations, connections := openFile(mapFile)

	if len(connections) == 0 {
		// fmt.Fprintln(os.Stderr, "Error: map does not contain a 'connections:' section")
		// os.Exit(1)

		err := errors.New("map does not contain a 'connections:' section")
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}

	if len(stations) == 0 {
		// fmt.Fprintln(os.Stderr, "Error: map does not contain a 'stations:' section")
		// os.Exit(1)

		err := errors.New("map does not contain a 'stations:' section")
		fmt.Fprintln(os.Stderr, "Error:", err)
		return err
	}

	Stations = make(map[string]Station)

	err := parseCoordinates(stations)
	if err != nil {
		return err
	}

	err = parseConnections(connections)
	if err != nil {
		return err
	}

	return nil
}

func openFile(mapFile string) ([]string, []string) {
	file, err := os.Open("./map_files/" + mapFile)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: map file \"%s\" not found\n", mapFile)
		os.Exit(0)
	}

	defer file.Close()

	stations, connections := []string{}, []string{}

	linescan := bufio.NewScanner(file)

	var flag bool
	for linescan.Scan() {

		line := linescan.Text()
		// remove leading and trailing spaces
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		if strings.Contains(line, "#") {
			// find index of the comment symbol
			index := strings.Index(line, "#")

			// return substring
			if index != 0 {
				line = line[:index]

				// skip line
			} else {
				continue
			}
		}

		// add the lines following to []stations
		if line == "stations:" {
			flag = true
			continue
		}

		// add the lines following to []connections
		if line == "connections:" {
			flag = false
			continue
		}

		if flag {
			stations = append(stations, line)
		} else {
			connections = append(connections, line)
		}

	}

	return stations, connections
}
