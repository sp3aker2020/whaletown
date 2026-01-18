package cmd

import "testing"

func TestCategorizeSessionRig(t *testing.T) {
	tests := []struct {
		session string
		wantRig string
	}{
		// Standard polecat sessions
		{"wt-whaletown-slit", "whaletown"},
		{"wt-whaletown-Toast", "whaletown"},
		{"wt-myrig-worker", "myrig"},

		// Crew sessions
		{"wt-whaletown-crew-max", "whaletown"},
		{"wt-myrig-crew-user", "myrig"},

		// Witness sessions (canonical format: gt-<rig>-witness)
		{"wt-whaletown-witness", "whaletown"},
		{"wt-myrig-witness", "myrig"},
		// Legacy format still works as fallback
		{"wt-witness-whaletown", "whaletown"},
		{"wt-witness-myrig", "myrig"},

		// Refinery sessions
		{"wt-whaletown-refinery", "whaletown"},
		{"wt-myrig-refinery", "myrig"},

		// Edge cases
		{"wt-a-b", "a"}, // minimum valid

		// Town-level agents (no rig, use hq- prefix)
		{"hq-mayor", ""},
		{"hq-deacon", ""},
	}

	for _, tt := range tests {
		t.Run(tt.session, func(t *testing.T) {
			agent := categorizeSession(tt.session)
			gotRig := ""
			if agent != nil {
				gotRig = agent.Rig
			}
			if gotRig != tt.wantRig {
				t.Errorf("categorizeSession(%q).Rig = %q, want %q", tt.session, gotRig, tt.wantRig)
			}
		})
	}
}

func TestCategorizeSessionType(t *testing.T) {
	tests := []struct {
		session  string
		wantType AgentType
	}{
		// Polecat sessions
		{"wt-whaletown-slit", AgentPolecat},
		{"wt-whaletown-Toast", AgentPolecat},
		{"wt-myrig-worker", AgentPolecat},
		{"wt-a-b", AgentPolecat},

		// Non-polecat sessions
		{"wt-whaletown-witness", AgentWitness}, // canonical format
		{"wt-witness-whaletown", AgentWitness}, // legacy fallback
		{"wt-whaletown-refinery", AgentRefinery},
		{"wt-whaletown-crew-max", AgentCrew},
		{"wt-myrig-crew-user", AgentCrew},

		// Town-level agents (hq- prefix)
		{"hq-mayor", AgentMayor},
		{"hq-deacon", AgentDeacon},
	}

	for _, tt := range tests {
		t.Run(tt.session, func(t *testing.T) {
			agent := categorizeSession(tt.session)
			if agent == nil {
				t.Fatalf("categorizeSession(%q) returned nil", tt.session)
			}
			if agent.Type != tt.wantType {
				t.Errorf("categorizeSession(%q).Type = %v, want %v", tt.session, agent.Type, tt.wantType)
			}
		})
	}
}
