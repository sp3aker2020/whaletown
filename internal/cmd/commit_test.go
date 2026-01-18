package cmd

import "testing"

func TestIdentityToEmail(t *testing.T) {
	tests := []struct {
		name     string
		identity string
		domain   string
		want     string
	}{
		{
			name:     "crew member",
			identity: "whaletown/crew/jack",
			domain:   "whaletown.local",
			want:     "whaletown.crew.jack@whaletown.local",
		},
		{
			name:     "polecat",
			identity: "whaletown/polecats/max",
			domain:   "whaletown.local",
			want:     "whaletown.polecats.max@whaletown.local",
		},
		{
			name:     "witness",
			identity: "whaletown/witness",
			domain:   "whaletown.local",
			want:     "whaletown.witness@whaletown.local",
		},
		{
			name:     "refinery",
			identity: "whaletown/refinery",
			domain:   "whaletown.local",
			want:     "whaletown.refinery@whaletown.local",
		},
		{
			name:     "mayor with trailing slash",
			identity: "mayor/",
			domain:   "whaletown.local",
			want:     "mayor@whaletown.local",
		},
		{
			name:     "deacon with trailing slash",
			identity: "deacon/",
			domain:   "whaletown.local",
			want:     "deacon@whaletown.local",
		},
		{
			name:     "custom domain",
			identity: "myrig/crew/alice",
			domain:   "example.com",
			want:     "myrig.crew.alice@example.com",
		},
		{
			name:     "deeply nested",
			identity: "rig/polecats/nested/deep",
			domain:   "test.io",
			want:     "rig.polecats.nested.deep@test.io",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := identityToEmail(tt.identity, tt.domain)
			if got != tt.want {
				t.Errorf("identityToEmail(%q, %q) = %q, want %q",
					tt.identity, tt.domain, got, tt.want)
			}
		})
	}
}
