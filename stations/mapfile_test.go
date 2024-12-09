package main

import (
	"strings"
	"testing"
	"trains/mapping"
)

func TestParseNetworkMap(t *testing.T) {
	tests := []struct {
		name    string
		mapFile string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no connections section",
			mapFile: "no_connections.map",
			wantErr: true,
			errMsg:  "map does not contain a 'connections:' section",
		},
		{
			name:    "no stations section",
			mapFile: "no_stations.map",
			wantErr: true,
			errMsg:  "map does not contain a 'stations:' section",
		},

		{
			name:    "neither section",
			mapFile: "neither.map",
			wantErr: true,
			errMsg:  "map does not contain a 'stations:' section",
		},

		{
			name:    "valid map file",
			mapFile: "artists.map",
			wantErr: false,
			errMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapping.ParseNetworkMap(tt.mapFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNetworkMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("ParseNetworkMap() error = %v, expected error message to contain %v", err, tt.errMsg)
			}
		})
	}
}
