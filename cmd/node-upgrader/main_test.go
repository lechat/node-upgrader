package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_Success(t *testing.T) {
	// Mock command-line arguments
	args := []string{"node-restarter", "--log-level=info", "--config-path=/path/to/config"}

	// Run the main function
	exitCode := run(args)

	// Assert exit code
	assert.Equal(t, 0, exitCode)
}

func TestRun_ArgumentParsingError(t *testing.T) {
	// Mock command-line arguments with invalid syntax
	args := []string{"node-restarter", "--invalid"}

	// Run the main function
	exitCode := run(args)

	// Assert exit code
	assert.Equal(t, 1, exitCode)
}

// Similar tests can be written for getConfig and other helper functions
