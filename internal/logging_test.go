package internal

import (
	"testing"
)

func TestInitializeLogger(t *testing.T) {
	tests := []struct {
		level      string
		wantLogger bool
		wantErr    bool
	}{
		{"debug", true, false},
		{"info", true, false},
		{"warn", true, false},
		{"error", true, false},
		{"invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			logger, err := InitializeLogger(tt.level)

			if (err != nil) != tt.wantErr {
				t.Errorf("InitializeLogger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantLogger && logger == nil {
				t.Error("InitializeLogger() did not return a logger, but it was expected")
			}

			if !tt.wantLogger && logger != nil {
				t.Error("InitializeLogger() returned a logger, but it was not expected")
			}
		})
	}
}
