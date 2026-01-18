package mail

import (
	"testing"
)

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		pattern string
		address string
		want    bool
	}{
		// Exact matches
		{"whaletown/witness", "whaletown/witness", true},
		{"mayor/", "mayor/", true},

		// Wildcard matches
		{"*/witness", "whaletown/witness", true},
		{"*/witness", "beads/witness", true},
		{"whaletown/*", "whaletown/witness", true},
		{"whaletown/*", "whaletown/refinery", true},
		{"whaletown/crew/*", "whaletown/crew/max", true},

		// Non-matches
		{"*/witness", "whaletown/refinery", false},
		{"whaletown/*", "beads/witness", false},
		{"whaletown/crew/*", "whaletown/polecats/Toast", false},

		// Different path lengths
		{"whaletown/*", "whaletown/crew/max", false},      // * matches single segment
		{"whaletown/*/*", "whaletown/crew/max", true},     // Multiple wildcards
		{"*/*", "whaletown/witness", true},              // Both wildcards
		{"*/*/*", "whaletown/crew/max", true},           // Three-level wildcard
	}

	for _, tt := range tests {
		t.Run(tt.pattern+"_"+tt.address, func(t *testing.T) {
			got := matchPattern(tt.pattern, tt.address)
			if got != tt.want {
				t.Errorf("matchPattern(%q, %q) = %v, want %v", tt.pattern, tt.address, got, tt.want)
			}
		})
	}
}

func TestAgentBeadIDToAddress(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		// Town-level agents
		{"wt-mayor", "mayor/"},
		{"wt-deacon", "deacon/"},

		// Rig singletons
		{"wt-whaletown-witness", "whaletown/witness"},
		{"wt-whaletown-refinery", "whaletown/refinery"},
		{"wt-beads-witness", "beads/witness"},

		// Named agents
		{"wt-whaletown-crew-max", "whaletown/crew/max"},
		{"wt-whaletown-polecat-Toast", "whaletown/polecat/Toast"},
		{"wt-beads-crew-wolf", "beads/crew/wolf"},

		// Agent with hyphen in name
		{"wt-whaletown-crew-max-v2", "whaletown/crew/max-v2"},
		{"wt-whaletown-polecat-my-agent", "whaletown/polecat/my-agent"},

		// Invalid
		{"invalid", ""},
		{"not-gt-prefix", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := agentBeadIDToAddress(tt.id)
			if got != tt.want {
				t.Errorf("agentBeadIDToAddress(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestResolverResolve_DirectAddresses(t *testing.T) {
	resolver := NewResolver(nil, "")

	tests := []struct {
		name    string
		address string
		want    RecipientType
		wantLen int
	}{
		// Direct agent addresses
		{"direct agent", "whaletown/witness", RecipientAgent, 1},
		{"direct crew", "whaletown/crew/max", RecipientAgent, 1},
		{"mayor", "mayor/", RecipientAgent, 1},

		// Legacy prefixes (pass-through)
		{"list prefix", "list:oncall", RecipientAgent, 1},
		{"announce prefix", "announce:alerts", RecipientAgent, 1},

		// Explicit type prefixes
		{"queue prefix", "queue:work", RecipientQueue, 1},
		{"channel prefix", "channel:alerts", RecipientChannel, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.Resolve(tt.address)
			if err != nil {
				t.Fatalf("Resolve(%q) error: %v", tt.address, err)
			}
			if len(got) != tt.wantLen {
				t.Errorf("Resolve(%q) returned %d recipients, want %d", tt.address, len(got), tt.wantLen)
			}
			if len(got) > 0 && got[0].Type != tt.want {
				t.Errorf("Resolve(%q)[0].Type = %v, want %v", tt.address, got[0].Type, tt.want)
			}
		})
	}
}

func TestResolverResolve_AtPatterns(t *testing.T) {
	// Without beads, @patterns are passed through for existing router
	resolver := NewResolver(nil, "")

	tests := []struct {
		address string
	}{
		{"@town"},
		{"@witnesses"},
		{"@rig/whaletown"},
		{"@overseer"},
	}

	for _, tt := range tests {
		t.Run(tt.address, func(t *testing.T) {
			got, err := resolver.Resolve(tt.address)
			if err != nil {
				t.Fatalf("Resolve(%q) error: %v", tt.address, err)
			}
			if len(got) != 1 {
				t.Errorf("Resolve(%q) returned %d recipients, want 1", tt.address, len(got))
			}
			// Without beads, @patterns pass through unchanged
			if got[0].Address != tt.address {
				t.Errorf("Resolve(%q) = %q, want pass-through", tt.address, got[0].Address)
			}
		})
	}
}

func TestResolverResolve_UnknownName(t *testing.T) {
	resolver := NewResolver(nil, "")

	// A bare name without prefix should fail if not found
	_, err := resolver.Resolve("unknown-name")
	if err == nil {
		t.Error("Resolve(\"unknown-name\") should return error for unknown name")
	}
}
