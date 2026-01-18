package cmd

import (
	"testing"
)

func TestResolveNudgePattern(t *testing.T) {
	// Create test agent sessions (mayor/deacon use hq- prefix)
	agents := []*AgentSession{
		{Name: "hq-mayor", Type: AgentMayor},
		{Name: "hq-deacon", Type: AgentDeacon},
		{Name: "wt-whaletown-witness", Type: AgentWitness, Rig: "whaletown"},
		{Name: "wt-whaletown-refinery", Type: AgentRefinery, Rig: "whaletown"},
		{Name: "wt-whaletown-crew-max", Type: AgentCrew, Rig: "whaletown", AgentName: "max"},
		{Name: "wt-whaletown-crew-jack", Type: AgentCrew, Rig: "whaletown", AgentName: "jack"},
		{Name: "wt-whaletown-alpha", Type: AgentPolecat, Rig: "whaletown", AgentName: "alpha"},
		{Name: "wt-whaletown-beta", Type: AgentPolecat, Rig: "whaletown", AgentName: "beta"},
		{Name: "wt-beads-witness", Type: AgentWitness, Rig: "beads"},
		{Name: "wt-beads-gamma", Type: AgentPolecat, Rig: "beads", AgentName: "gamma"},
	}

	tests := []struct {
		name     string
		pattern  string
		expected []string
	}{
		{
			name:     "mayor special case",
			pattern:  "mayor",
			expected: []string{"hq-mayor"},
		},
		{
			name:     "deacon special case",
			pattern:  "deacon",
			expected: []string{"hq-deacon"},
		},
		{
			name:     "specific witness",
			pattern:  "whaletown/witness",
			expected: []string{"wt-whaletown-witness"},
		},
		{
			name:     "all witnesses",
			pattern:  "*/witness",
			expected: []string{"wt-whaletown-witness", "wt-beads-witness"},
		},
		{
			name:     "specific refinery",
			pattern:  "whaletown/refinery",
			expected: []string{"wt-whaletown-refinery"},
		},
		{
			name:     "all polecats in rig",
			pattern:  "whaletown/polecats/*",
			expected: []string{"wt-whaletown-alpha", "wt-whaletown-beta"},
		},
		{
			name:     "specific polecat",
			pattern:  "whaletown/polecats/alpha",
			expected: []string{"wt-whaletown-alpha"},
		},
		{
			name:     "all crew in rig",
			pattern:  "whaletown/crew/*",
			expected: []string{"wt-whaletown-crew-max", "wt-whaletown-crew-jack"},
		},
		{
			name:     "specific crew member",
			pattern:  "whaletown/crew/max",
			expected: []string{"wt-whaletown-crew-max"},
		},
		{
			name:     "legacy polecat format",
			pattern:  "whaletown/alpha",
			expected: []string{"wt-whaletown-alpha"},
		},
		{
			name:     "no matches",
			pattern:  "nonexistent/polecats/*",
			expected: nil,
		},
		{
			name:     "invalid pattern",
			pattern:  "invalid",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveNudgePattern(tt.pattern, agents)

			if len(got) != len(tt.expected) {
				t.Errorf("resolveNudgePattern(%q) returned %d results, want %d: got %v, want %v",
					tt.pattern, len(got), len(tt.expected), got, tt.expected)
				return
			}

			// Check each expected value is present
			gotMap := make(map[string]bool)
			for _, g := range got {
				gotMap[g] = true
			}
			for _, e := range tt.expected {
				if !gotMap[e] {
					t.Errorf("resolveNudgePattern(%q) missing expected %q, got %v",
						tt.pattern, e, got)
				}
			}
		})
	}
}
