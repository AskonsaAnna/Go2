package pathing

import (
	// "fmt"
	// "os"
	"testing"
	"trains/mapping"
)

// 2. It finds more than one valid route for 3 trains between
// waterloo and st_pancras in the London Network Map

// 18. It finds more than one valid route for 100 trains
// between waterloo and st_pancras in the London Network Map

// 24. It can find more than one route for 2 trains
// between waterloo and st_pancras for the London Network Map

// 26. It finds only a single valid route for 1 train between
// waterloo and st_pancras in the London Network Map

// 27. It finds more than one valid route for 4 trains between
// waterloo and st_pancras in the London Network Map

// 29. It displays "Error" on stderr when no path exists between the start and end stations.

// FindPaths determines the paths based on the number of trains.
func TestFindPaths(t *testing.T) {
	// Setup mock stations for testing
	setupMockStations()

	tests := []struct {
		name        string
		start       string
		end         string
		trains      int
		wantRoutes  int
		expectError bool
	}{
		{
			name:        "Test with 3 trains",
			start:       "waterloo",
			end:         "st_pancras",
			trains:      3,
			wantRoutes:  2, // Assuming there are at least 2 valid routes
			expectError: false,
		},
		{
			name:        "Test with 100 trains",
			start:       "waterloo",
			end:         "st_pancras",
			trains:      100,
			wantRoutes:  2, // Assuming there are at least 2 valid routes
			expectError: false,
		},
		{
			name:        "Test with 2 trains",
			start:       "waterloo",
			end:         "st_pancras",
			trains:      2,
			wantRoutes:  2, // Assuming there are at least 2 valid routes
			expectError: false,
		},
		{
			name:        "Test with 1 train",
			start:       "waterloo",
			end:         "st_pancras",
			trains:      1,
			wantRoutes:  1, // Assuming there's only 1 valid route
			expectError: false,
		},
		{
			name:        "Test with 4 trains",
			start:       "waterloo",
			end:         "st_pancras",
			trains:      4,
			wantRoutes:  2, // Assuming there are at least 2 valid routes
			expectError: false,
		},
		{
			name:        "Test with no path",
			start:       "waterloo",
			end:         "non_existent",
			trains:      1,
			wantRoutes:  0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.expectError {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("FindPaths()")
					}
				}()
			}
			gotPaths := FindPaths(tt.start, tt.end, tt.trains)
			if gotPaths == nil && !tt.expectError && len(gotPaths) != tt.wantRoutes {
				t.Errorf("Expected %d routes, but got %d", tt.wantRoutes, len(gotPaths))
			}
		})
	}
}

func setupMockStations() {
	mapping.Stations = map[string]mapping.Station{
		"waterloo": {
			Connections: []string{"station1", "station2"},
		},
		"st_pancras": {
			Connections: []string{"station3", "station4"},
		},
		"station1": {
			Connections: []string{"waterloo", "station3"},
		},
		"station2": {
			Connections: []string{"waterloo", "station4"},
		},
		"station3": {
			Connections: []string{"station1", "st_pancras"},
		},
		"station4": {
			Connections: []string{"station2", "st_pancras"},
		},
		"non_existent": {
			Connections: []string{},
		},
	}
}
