package mapping

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseCoordinates(t *testing.T) {
	tests := []struct {
		name     string
		stations []string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid input",
			stations: []string{
				"station1,10,20",
				"station2,30,40",
			},
			wantErr: false,
		},
		{
			name: "duplicate station names",
			stations: []string{
				"station1,10,20",
				"station1,30,40",
			},
			wantErr: true,
			errMsg:  "duplicate station name",
		},
		{
			name: "invalid format",
			stations: []string{
				"station1,10,20",
				"station2,30",
			},
			wantErr: true,
			errMsg:  "invalid station format",
		},
		{
			name: "invalid coordinates",
			stations: []string{
				"station1,10,20",
				"station2,-30,40",
			},
			wantErr: true,
			errMsg:  "invalid coordinates",
		},
		{
			name: "duplicate coordinates",
			stations: []string{
				"station1,10,20",
				"station2,10,20",
			},
			wantErr: true,
			errMsg:  "duplicate coordinates for stations 'station2' and 'station1'",
		},
		{
			name: "more than 10000 stations",
			stations: func() []string {
				st := make([]string, 10001)
				for i := 0; i < 10001; i++ {
					st[i] = fmt.Sprintf("Station%d,%d,%d", i, i*10, i*10+1)
				}
				return st
			}(),
			wantErr: true,
			errMsg:  "map contains more than 10000 stations",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Stations = make(map[string]Station) // reset the Stations map before each test
			err := parseCoordinates(tt.stations)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCoordinates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("parseCoordinates() error = %v, expected error message to contain %v", err, tt.errMsg)
			}
		})
	}
}

func TestParseConnections(t *testing.T) {
	tests := []struct {
		name        string
		stations    []string
		connections []string
		wantErr     bool
		errMsg      string
	}{
		{
			name: "station does not exist",
			stations: []string{
				"station1,10,20",
				"station2,30,40",
			},
			connections: []string{
				"station3-station2",
			},
			wantErr: true,
			errMsg:  "station does not exist: station3",
		},
		{
			name: "station does not exist",
			stations: []string{
				"station1,10,20",
				"station2,30,40",
			},
			connections: []string{
				"station1-station3",
			},
			wantErr: true,
			errMsg:  "station does not exist: station3",
		},

		{
			name: "duplicate connection",
			stations: []string{
				"station1,10,20",
				"station2,30,40",
			},
			connections: []string{
				"station1-station2",
				"station1-station2",
			},
			wantErr: true,
			errMsg:  "duplicate connection between station1 and station2",
		},

		{
			name: "duplicate connection",
			stations: []string{
				"station1,10,20",
				"station2,30,40",
			},
			connections: []string{
				"station1-station2",
				"station2-station1",
			},
			wantErr: true,
			errMsg:  "duplicate connection between station2 and station1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Stations = make(map[string]Station) // reset the Stations map before each test

			err := parseCoordinates(tt.stations) // add stations first
			if err != nil {
				t.Fatalf("failed to parse stations: %v", err)
			}

			err = parseConnections(tt.connections)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConnections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("parseConnections() error = %v, expected error message to contain %v", err, tt.errMsg)
			}
		})
	}
}
