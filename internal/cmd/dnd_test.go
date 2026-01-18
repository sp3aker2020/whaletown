package cmd

import (
	"testing"
)

func TestAddressToAgentBeadID(t *testing.T) {
	tests := []struct {
		address  string
		expected string
	}{
		// Mayor and deacon use hq- prefix (town-level)
		{"mayor", "hq-mayor"},
		{"deacon", "hq-deacon"},
		{"whaletown/witness", "wt-whaletown-witness"},
		{"whaletown/refinery", "wt-whaletown-refinery"},
		{"whaletown/alpha", "wt-whaletown-polecat-alpha"},
		{"whaletown/crew/max", "wt-whaletown-crew-max"},
		{"beads/witness", "wt-beads-witness"},
		{"beads/beta", "wt-beads-polecat-beta"},
		// Invalid addresses should return empty string
		{"invalid", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.address, func(t *testing.T) {
			got := addressToAgentBeadID(tt.address)
			if got != tt.expected {
				t.Errorf("addressToAgentBeadID(%q) = %q, want %q", tt.address, got, tt.expected)
			}
		})
	}
}
