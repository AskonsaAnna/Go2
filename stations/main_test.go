package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		expectError  bool
		errorMessage string
	}{
		{"ExtraArg", []string{"program", "london.map", "waterloo", "st_pancras", "1", "extraArg"}, true, "Error: Incorrect number of command line arguments."},
		{"MissingArg", []string{"program", "london.map", "waterloo", "st_pancras"}, true, "Error: Incorrect number of command line arguments."},
		{"NonexistentStartStation", []string{"program", "london.map", "nonexistent", "st_pancras", "1"}, true, "Error: Start station does not exist."},
		{"NonexistentEndStation", []string{"program", "london.map", "waterloo", "nonexistent", "1"}, true, "Error: End station does not exist."},
		{"SameStartEndStation", []string{"program", "london.map", "waterloo", "waterloo", "1"}, true, "Error: Start and end station are the same."},
	}

	for _, test := range tests {

		cmd := exec.Command("./program", test.args[1:]...)
		cmd.Env = os.Environ()
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		err := cmd.Run()

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for args %v, but got none", test.args)
			} else if !strings.Contains(stderr.String(), test.errorMessage) {
				t.Errorf("Expected error message '%s' for args %v, but got '%s'", test.errorMessage, test.args, stderr.String())
			}
		} else if err != nil {
			t.Errorf("Unexpected error for args %v: %v", test.args, err)
		}

	}
}

func TestMovementInSixTurnsBeethovenToPart(t *testing.T) {
	args := []string{"program", "composers.map", "beethoven", "part", "9"}
	checkTurns(t, args, 6)
}

func TestMovementInEightTurnsSmallToLarge(t *testing.T) {
	args := []string{"program", "digits.map", "small", "large", "9"}
	checkTurns(t, args, 8)
}

func TestMovementInSixTurnsTwoToFour(t *testing.T) {
	args := []string{"program", "numbers.map", "two", "four", "4"}
	checkTurns(t, args, 6)
}

func TestMovementInEightTurnsJungleToDesert(t *testing.T) {
	args := []string{"program", "terrains.map", "jungle", "desert", "10"}
	checkTurns(t, args, 8)
}

func TestMovementInSixTurnsBondSquareToSpacePort(t *testing.T) {
	args := []string{"program", "bond.map", "bond_square", "space_port", "4"}
	checkTurns(t, args, 6)
}

func TestMovementInElevenTurnsBeginningToTerminus(t *testing.T) {
	args := []string{"program", "terminus.map", "beginning", "terminus", "20"}
	checkTurns(t, args, 11)
}

func checkTurns(t *testing.T, args []string, maxTurns int) {
	cmd := exec.Command("./program", args[1:]...)
	cmd.Env = os.Environ()
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run command: %v", err)
	}

	output := stdout.String()
	expectedOutput := fmt.Sprintf("Moved %s trains from %s to %s in %d turns\n", args[4], args[2], args[3], maxTurns)
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected '%s', but got '%s'", expectedOutput, output)
	} else {
		t.Logf("\n\n%s", output)
	}
}
