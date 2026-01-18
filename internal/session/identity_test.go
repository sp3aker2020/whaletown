package session

import (
	"testing"
)

func TestParseSessionName(t *testing.T) {
	tests := []struct {
		name     string
		session  string
		wantRole Role
		wantRig  string
		wantName string
		wantErr  bool
	}{
		// Town-level roles (hq-mayor, hq-deacon)
		{
			name:     "mayor",
			session:  "hq-mayor",
			wantRole: RoleMayor,
		},
		{
			name:     "deacon",
			session:  "hq-deacon",
			wantRole: RoleDeacon,
		},

		// Witness (simple rig)
		{
			name:     "witness simple rig",
			session:  "wt-whaletown-witness",
			wantRole: RoleWitness,
			wantRig:  "whaletown",
		},
		{
			name:     "witness hyphenated rig",
			session:  "wt-foo-bar-witness",
			wantRole: RoleWitness,
			wantRig:  "foo-bar",
		},

		// Refinery (simple rig)
		{
			name:     "refinery simple rig",
			session:  "wt-whaletown-refinery",
			wantRole: RoleRefinery,
			wantRig:  "whaletown",
		},
		{
			name:     "refinery hyphenated rig",
			session:  "wt-my-project-refinery",
			wantRole: RoleRefinery,
			wantRig:  "my-project",
		},

		// Crew (with marker)
		{
			name:     "crew simple",
			session:  "wt-whaletown-crew-max",
			wantRole: RoleCrew,
			wantRig:  "whaletown",
			wantName: "max",
		},
		{
			name:     "crew hyphenated rig",
			session:  "wt-foo-bar-crew-alice",
			wantRole: RoleCrew,
			wantRig:  "foo-bar",
			wantName: "alice",
		},
		{
			name:     "crew hyphenated name",
			session:  "wt-whaletown-crew-my-worker",
			wantRole: RoleCrew,
			wantRig:  "whaletown",
			wantName: "my-worker",
		},

		// Polecat (fallback)
		{
			name:     "polecat simple",
			session:  "wt-whaletown-morsov",
			wantRole: RolePolecat,
			wantRig:  "whaletown",
			wantName: "morsov",
		},
		{
			name:     "polecat hyphenated rig",
			session:  "wt-foo-bar-Toast",
			wantRole: RolePolecat,
			wantRig:  "foo-bar",
			wantName: "Toast",
		},

		// Error cases
		{
			name:    "missing prefix",
			session: "whaletown-witness",
			wantErr: true,
		},
		{
			name:    "empty after prefix",
			session: "wt-",
			wantErr: true,
		},
		{
			name:    "just prefix single segment",
			session: "wt-x",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSessionName(tt.session)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSessionName(%q) error = %v, wantErr %v", tt.session, err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.Role != tt.wantRole {
				t.Errorf("ParseSessionName(%q).Role = %v, want %v", tt.session, got.Role, tt.wantRole)
			}
			if got.Rig != tt.wantRig {
				t.Errorf("ParseSessionName(%q).Rig = %v, want %v", tt.session, got.Rig, tt.wantRig)
			}
			if got.Name != tt.wantName {
				t.Errorf("ParseSessionName(%q).Name = %v, want %v", tt.session, got.Name, tt.wantName)
			}
		})
	}
}

func TestAgentIdentity_SessionName(t *testing.T) {
	tests := []struct {
		name     string
		identity AgentIdentity
		want     string
	}{
		{
			name:     "mayor",
			identity: AgentIdentity{Role: RoleMayor},
			want:     "hq-mayor",
		},
		{
			name:     "deacon",
			identity: AgentIdentity{Role: RoleDeacon},
			want:     "hq-deacon",
		},
		{
			name:     "witness",
			identity: AgentIdentity{Role: RoleWitness, Rig: "whaletown"},
			want:     "wt-whaletown-witness",
		},
		{
			name:     "refinery",
			identity: AgentIdentity{Role: RoleRefinery, Rig: "my-project"},
			want:     "wt-my-project-refinery",
		},
		{
			name:     "crew",
			identity: AgentIdentity{Role: RoleCrew, Rig: "whaletown", Name: "max"},
			want:     "wt-whaletown-crew-max",
		},
		{
			name:     "polecat",
			identity: AgentIdentity{Role: RolePolecat, Rig: "whaletown", Name: "morsov"},
			want:     "wt-whaletown-morsov",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.identity.SessionName(); got != tt.want {
				t.Errorf("AgentIdentity.SessionName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentIdentity_Address(t *testing.T) {
	tests := []struct {
		name     string
		identity AgentIdentity
		want     string
	}{
		{
			name:     "mayor",
			identity: AgentIdentity{Role: RoleMayor},
			want:     "mayor",
		},
		{
			name:     "deacon",
			identity: AgentIdentity{Role: RoleDeacon},
			want:     "deacon",
		},
		{
			name:     "witness",
			identity: AgentIdentity{Role: RoleWitness, Rig: "whaletown"},
			want:     "whaletown/witness",
		},
		{
			name:     "refinery",
			identity: AgentIdentity{Role: RoleRefinery, Rig: "my-project"},
			want:     "my-project/refinery",
		},
		{
			name:     "crew",
			identity: AgentIdentity{Role: RoleCrew, Rig: "whaletown", Name: "max"},
			want:     "whaletown/crew/max",
		},
		{
			name:     "polecat",
			identity: AgentIdentity{Role: RolePolecat, Rig: "whaletown", Name: "Toast"},
			want:     "whaletown/polecats/Toast",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.identity.Address(); got != tt.want {
				t.Errorf("AgentIdentity.Address() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSessionName_RoundTrip(t *testing.T) {
	// Test that parsing then reconstructing gives the same result
	sessions := []string{
		"hq-mayor",
		"hq-deacon",
		"wt-whaletown-witness",
		"wt-foo-bar-refinery",
		"wt-whaletown-crew-max",
		"wt-whaletown-morsov",
	}

	for _, sess := range sessions {
		t.Run(sess, func(t *testing.T) {
			identity, err := ParseSessionName(sess)
			if err != nil {
				t.Fatalf("ParseSessionName(%q) error = %v", sess, err)
			}
			if got := identity.SessionName(); got != sess {
				t.Errorf("Round-trip failed: ParseSessionName(%q).SessionName() = %q", sess, got)
			}
		})
	}
}
