package cmd

import (
	"os"
	"testing"
)

func TestDeriveSessionName(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		expected string
	}{
		{
			name: "polecat session",
			envVars: map[string]string{
				"WT_ROLE":    "polecat",
				"WT_RIG":     "whaletown",
				"WT_POLECAT": "toast",
			},
			expected: "wt-whaletown-toast",
		},
		{
			name: "crew session",
			envVars: map[string]string{
				"WT_ROLE": "crew",
				"WT_RIG":  "whaletown",
				"WT_CREW": "max",
			},
			expected: "wt-whaletown-crew-max",
		},
		{
			name: "witness session",
			envVars: map[string]string{
				"WT_ROLE": "witness",
				"WT_RIG":  "whaletown",
			},
			expected: "wt-whaletown-witness",
		},
		{
			name: "refinery session",
			envVars: map[string]string{
				"WT_ROLE": "refinery",
				"WT_RIG":  "whaletown",
			},
			expected: "wt-whaletown-refinery",
		},
		{
			name: "mayor session",
			envVars: map[string]string{
				"WT_ROLE": "mayor",
				"WT_TOWN": "ai",
			},
			expected: "wt-ai-mayor",
		},
		{
			name: "deacon session",
			envVars: map[string]string{
				"WT_ROLE": "deacon",
				"WT_TOWN": "ai",
			},
			expected: "wt-ai-deacon",
		},
		{
			name: "mayor session without WT_TOWN",
			envVars: map[string]string{
				"WT_ROLE": "mayor",
			},
			expected: "wt-mayor",
		},
		{
			name: "deacon session without WT_TOWN",
			envVars: map[string]string{
				"WT_ROLE": "deacon",
			},
			expected: "wt-deacon",
		},
		{
			name:     "no env vars",
			envVars:  map[string]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and clear relevant env vars
			saved := make(map[string]string)
			envKeys := []string{"WT_ROLE", "WT_RIG", "WT_POLECAT", "WT_CREW", "WT_TOWN"}
			for _, key := range envKeys {
				saved[key] = os.Getenv(key)
				os.Unsetenv(key)
			}
			defer func() {
				// Restore env vars
				for key, val := range saved {
					if val != "" {
						os.Setenv(key, val)
					}
				}
			}()

			// Set test env vars
			for key, val := range tt.envVars {
				os.Setenv(key, val)
			}

			result := deriveSessionName()
			if result != tt.expected {
				t.Errorf("deriveSessionName() = %q, want %q", result, tt.expected)
			}
		})
	}
}
