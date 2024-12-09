package mapping

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func parseCoordinates(stations []string) error {
	if len(stations) > 10000 {
		err := errors.New("map contains more than 10000 stations") // create error
		fmt.Fprintln(os.Stderr, "Error:", err)                     // display on stderr
		return err                                                 // return err
	}

	coordinatesMap := make(map[string]string)
	var validNameRegex = regexp.MustCompile(`^[a-z0-9_]+$`)

	for _, line := range stations {
		parts := strings.Split(line, ",")

		if len(parts) != 3 {
			err := errors.New("invalid station format")
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}

		name := strings.TrimSpace(parts[0])

		if !validNameRegex.MatchString(name) {
			err := fmt.Errorf("invalid station name: %s", name)
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}

		// Check for duplicate station names
		if _, exists := Stations[name]; exists {
			err := fmt.Errorf("duplicate station name: %s", name)
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}

		x, err1 := strconv.Atoi(strings.TrimSpace(parts[1]))
		y, err2 := strconv.Atoi(strings.TrimSpace(parts[2]))

		if err1 != nil || err2 != nil || x < 0 || y < 0 {
			err := errors.New("invalid coordinates")
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}

		coordKey := fmt.Sprintf("%d,%d", x, y)
		// Check for duplicate coordinates
		if existingStation, exists := coordinatesMap[coordKey]; exists {
			err := fmt.Errorf("duplicate coordinates for stations '%s' and '%s'", name, existingStation)
			fmt.Fprintln(os.Stderr, "Error:", err)
			return err
		}
		coordinatesMap[coordKey] = name

		Stations[name] = Station{
			Coordinates: []int{x, y},
			Connections: []string{},
		}
	}
	return nil
}

func parseConnections(connections []string) error {
	for _, connection := range connections {
		parts := strings.Split(connection, "-")
		name := strings.TrimSpace(parts[0])
		connected := strings.TrimSpace(parts[1])

		// add the connection to map

		station := Stations[name]

		if _, exists := Stations[name]; !exists {
			err := fmt.Errorf("station does not exist: %s", name)
			fmt.Fprintln(os.Stderr, "Error:", err)
			//os.Exit(1)
			return err
		}

		if _, exists := Stations[connected]; !exists {
			err := fmt.Errorf("station does not exist: %s", connected)
			fmt.Fprintln(os.Stderr, "Error:", err)
			//os.Exit(1)
			return err
		}

		for _, dooble := range station.Connections {
			if dooble == connected {
				err := fmt.Errorf("duplicate connection between %s and %s", name, connected)
				fmt.Fprintln(os.Stderr, "Error:", err)
				//os.Exit(1)
				return err
			}
		}

		station.Connections = append(station.Connections, connected)
		Stations[name] = station

		// add the reverse connection also
		station = Stations[connected]
		station.Connections = append(station.Connections, name)
		Stations[connected] = station
	}

	return nil

}
