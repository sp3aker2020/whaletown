//go:build integration

package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// cleanGTEnv returns os.Environ() with all WT_* variables removed.
// This ensures tests don't inherit stale role environment from CI or previous tests.
func cleanGTEnv() []string {
	var clean []string
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "WT_") {
			clean = append(clean, env)
		}
	}
	return clean
}

// TestRoleHomeE2E validates that wt role home returns correct paths
// for all role types after a full wt install.
func TestRoleHomeE2E(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	rigName := "testrig"

	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "mayor",
			args:     []string{"role", "home", "mayor"},
			expected: filepath.Join(hqPath, "mayor"),
		},
		{
			name:     "deacon",
			args:     []string{"role", "home", "deacon"},
			expected: filepath.Join(hqPath, "deacon"),
		},
		{
			name:     "witness",
			args:     []string{"role", "home", "witness", "--rig", rigName},
			expected: filepath.Join(hqPath, rigName, "witness"),
		},
		{
			name:     "refinery",
			args:     []string{"role", "home", "refinery", "--rig", rigName},
			expected: filepath.Join(hqPath, rigName, "refinery", "rig"),
		},
		{
			name:     "polecat",
			args:     []string{"role", "home", "polecat", "--rig", rigName, "--polecat", "Toast"},
			expected: filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
		},
		{
			name:     "crew",
			args:     []string{"role", "home", "crew", "--rig", rigName, "--polecat", "worker1"},
			expected: filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, tt.args...)
			cmd.Dir = hqPath
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)

			// Use Output() to only capture stdout (warnings go to stderr)
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("gt %v failed: %v\nOutput: %s", tt.args, err, output)
			}

			got := strings.TrimSpace(string(output))
			if got != tt.expected {
				t.Errorf("gt %v = %q, want %q", tt.args, got, tt.expected)
			}
		})
	}
}

// TestRoleHomeMissingFlags validates that wt role home fails when required flags are missing.
func TestRoleHomeMissingFlags(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "witness without --rig",
			args: []string{"role", "home", "witness"},
		},
		{
			name: "refinery without --rig",
			args: []string{"role", "home", "refinery"},
		},
		{
			name: "polecat without --rig",
			args: []string{"role", "home", "polecat", "--polecat", "Toast"},
		},
		{
			name: "polecat without --polecat",
			args: []string{"role", "home", "polecat", "--rig", "testrig"},
		},
		{
			name: "crew without --rig",
			args: []string{"role", "home", "crew", "--polecat", "worker1"},
		},
		{
			name: "crew without --polecat",
			args: []string{"role", "home", "crew", "--rig", "testrig"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, tt.args...)
			cmd.Dir = hqPath
			// Use cleanGTEnv to ensure no stale WT_* vars affect the test
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)

			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Errorf("gt %v should have failed but succeeded with output: %s", tt.args, output)
			}
		})
	}
}


// TestRoleHomeCwdDetection validates wt role home without arguments detects role from cwd.
func TestRoleHomeCwdDetection(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	rigName := "testrig"

	// Create rig directory structure for cwd detection
	dirs := []string{
		filepath.Join(hqPath, rigName, "witness"),
		filepath.Join(hqPath, rigName, "refinery", "rig"),
		filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
		filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	tests := []struct {
		name     string
		cwd      string
		expected string
	}{
		{
			name:     "mayor from mayor dir",
			cwd:      filepath.Join(hqPath, "mayor"),
			expected: filepath.Join(hqPath, "mayor"),
		},
		{
			name:     "deacon from deacon dir",
			cwd:      filepath.Join(hqPath, "deacon"),
			expected: filepath.Join(hqPath, "deacon"),
		},
		{
			name:     "witness from witness dir",
			cwd:      filepath.Join(hqPath, rigName, "witness"),
			expected: filepath.Join(hqPath, rigName, "witness"),
		},
		{
			name:     "refinery from refinery/rig dir",
			cwd:      filepath.Join(hqPath, rigName, "refinery", "rig"),
			expected: filepath.Join(hqPath, rigName, "refinery", "rig"),
		},
		{
			name:     "polecat from polecats/Toast/rig dir",
			cwd:      filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
			expected: filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
		},
		{
			name:     "crew from crew/worker1/rig dir",
			cwd:      filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
			expected: filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, "role", "home")
			cmd.Dir = tt.cwd
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("gt role home failed: %v\nOutput: %s", err, output)
			}

			got := strings.TrimSpace(string(output))
			if got != tt.expected {
				t.Errorf("gt role home = %q, want %q", got, tt.expected)
			}
		})
	}
}

// TestRoleEnvCwdDetection validates wt role env without arguments detects role from cwd.
func TestRoleEnvCwdDetection(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	rigName := "testrig"

	// Create rig directory structure for cwd detection
	dirs := []string{
		filepath.Join(hqPath, rigName, "witness"),
		filepath.Join(hqPath, rigName, "refinery", "rig"),
		filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
		filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	tests := []struct {
		name string
		cwd  string
		want []string
	}{
		{
			name: "mayor from mayor dir",
			cwd:  filepath.Join(hqPath, "mayor"),
			want: []string{
				"export WT_ROLE=mayor",
				"export WT_ROLE_HOME=" + filepath.Join(hqPath, "mayor"),
			},
		},
		{
			name: "deacon from deacon dir",
			cwd:  filepath.Join(hqPath, "deacon"),
			want: []string{
				"export WT_ROLE=deacon",
				"export WT_ROLE_HOME=" + filepath.Join(hqPath, "deacon"),
			},
		},
		{
			name: "witness from witness dir",
			cwd:  filepath.Join(hqPath, rigName, "witness"),
			want: []string{
				"export WT_ROLE=witness",
				"export WT_RIG=" + rigName,
				"export BD_ACTOR=" + rigName + "/witness",
				"export WT_ROLE_HOME=" + filepath.Join(hqPath, rigName, "witness"),
			},
		},
		{
			name: "refinery from refinery/rig dir",
			cwd:  filepath.Join(hqPath, rigName, "refinery", "rig"),
			want: []string{
				"export WT_ROLE=refinery",
				"export WT_RIG=" + rigName,
				"export BD_ACTOR=" + rigName + "/refinery",
				"export WT_ROLE_HOME=" + filepath.Join(hqPath, rigName, "refinery", "rig"),
			},
		},
		{
			name: "polecat from polecats/Toast/rig dir",
			cwd:  filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
			want: []string{
				"export WT_ROLE=polecat",
				"export WT_RIG=" + rigName,
				"export WT_POLECAT=Toast",
				"export BD_ACTOR=" + rigName + "/polecats/Toast",
				"export WT_ROLE_HOME=" + filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
			},
		},
		{
			name: "crew from crew/worker1/rig dir",
			cwd:  filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
			want: []string{
				"export WT_ROLE=crew",
				"export WT_RIG=" + rigName,
				"export WT_CREW=worker1",
				"export BD_ACTOR=" + rigName + "/crew/worker1",
				"export WT_ROLE_HOME=" + filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, "role", "env")
			cmd.Dir = tt.cwd
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("gt role env failed: %v\nOutput: %s", err, output)
			}

			got := string(output)
			for _, w := range tt.want {
				if !strings.Contains(got, w) {
					t.Errorf("output missing %q\ngot: %s", w, got)
				}
			}
		})
	}
}

// TestRoleListE2E validates wt role list shows all roles.
func TestRoleListE2E(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	cmd = exec.Command(gtBinary, "role", "list")
	cmd.Dir = hqPath
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gt role list failed: %v\nOutput: %s", err, output)
	}

	got := string(output)

	// Check header
	if !strings.Contains(got, "Available roles:") {
		t.Errorf("output missing 'Available roles:' header\ngot: %s", got)
	}

	// Check all roles are listed
	roles := []string{"mayor", "deacon", "witness", "refinery", "polecat", "crew"}
	for _, role := range roles {
		if !strings.Contains(got, role) {
			t.Errorf("output missing role %q\ngot: %s", role, got)
		}
	}
}

// TestRoleShowE2E validates wt role show displays correct role info.
func TestRoleShowE2E(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	rigName := "testrig"

	// Create rig directory structure for cwd detection
	dirs := []string{
		filepath.Join(hqPath, rigName, "witness"),
		filepath.Join(hqPath, rigName, "refinery", "rig"),
		filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
		filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	tests := []struct {
		name       string
		cwd        string
		wantRole   string
		wantSource string
		wantHome   string
		wantRig    string
		wantWorker string
	}{
		{
			name:       "mayor from mayor dir",
			cwd:        filepath.Join(hqPath, "mayor"),
			wantRole:   "mayor",
			wantSource: "cwd",
			wantHome:   filepath.Join(hqPath, "mayor"),
		},
		{
			name:       "deacon from deacon dir",
			cwd:        filepath.Join(hqPath, "deacon"),
			wantRole:   "deacon",
			wantSource: "cwd",
			wantHome:   filepath.Join(hqPath, "deacon"),
		},
		{
			name:       "witness from witness dir",
			cwd:        filepath.Join(hqPath, rigName, "witness"),
			wantRole:   "witness",
			wantSource: "cwd",
			wantHome:   filepath.Join(hqPath, rigName, "witness"),
			wantRig:    rigName,
		},
		{
			name:       "polecat from polecats/Toast/rig dir",
			cwd:        filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
			wantRole:   "polecat",
			wantSource: "cwd",
			wantHome:   filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
			wantRig:    rigName,
			wantWorker: "Toast",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, "role", "show")
			cmd.Dir = tt.cwd
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("gt role show failed: %v\nOutput: %s", err, output)
			}

			got := string(output)

			if !strings.Contains(got, tt.wantRole) {
				t.Errorf("output missing role %q\ngot: %s", tt.wantRole, got)
			}

			if !strings.Contains(got, "Source: "+tt.wantSource) {
				t.Errorf("output missing 'Source: %s'\ngot: %s", tt.wantSource, got)
			}

			if !strings.Contains(got, "Home: "+tt.wantHome) {
				t.Errorf("output missing 'Home: %s'\ngot: %s", tt.wantHome, got)
			}

			if tt.wantRig != "" {
				if !strings.Contains(got, "Rig: "+tt.wantRig) {
					t.Errorf("output missing 'Rig: %s'\ngot: %s", tt.wantRig, got)
				}
			}

			if tt.wantWorker != "" {
				if !strings.Contains(got, "Worker: "+tt.wantWorker) {
					t.Errorf("output missing 'Worker: %s'\ngot: %s", tt.wantWorker, got)
				}
			}
		})
	}
}

// TestRoleShowMismatch validates wt role show displays mismatch warning.
func TestRoleShowMismatch(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	// Run from mayor dir but set WT_ROLE to deacon
	cmd = exec.Command(gtBinary, "role", "show")
	cmd.Dir = filepath.Join(hqPath, "mayor")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir, "WT_ROLE=deacon")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gt role show failed: %v\nOutput: %s", err, output)
	}

	got := string(output)

	// WT_ROLE takes precedence, so role should be deacon
	if !strings.Contains(got, "deacon") {
		t.Errorf("should show 'deacon' from WT_ROLE, got: %s", got)
	}

	// Source should be env
	if !strings.Contains(got, "Source: env") {
		t.Errorf("source should be 'env', got: %s", got)
	}

	// Should show mismatch warning
	if !strings.Contains(got, "ROLE MISMATCH") {
		t.Errorf("should show ROLE MISMATCH warning\ngot: %s", got)
	}

	// Should show both the env value and cwd suggestion
	if !strings.Contains(got, "WT_ROLE=deacon") {
		t.Errorf("should show WT_ROLE value\ngot: %s", got)
	}

	if !strings.Contains(got, "cwd suggests: mayor") {
		t.Errorf("should show cwd suggestion\ngot: %s", got)
	}
}

// TestRoleDetectE2E validates wt role detect uses cwd and ignores WT_ROLE.
func TestRoleDetectE2E(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	rigName := "testrig"

	// Create rig directory structure for cwd detection
	dirs := []string{
		filepath.Join(hqPath, rigName, "witness"),
		filepath.Join(hqPath, rigName, "refinery", "rig"),
		filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
		filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	tests := []struct {
		name        string
		cwd         string
		wantRole    string
		wantRig     string
		wantWorker  string
	}{
		{
			name:     "mayor from mayor dir",
			cwd:      filepath.Join(hqPath, "mayor"),
			wantRole: "mayor",
		},
		{
			name:     "deacon from deacon dir",
			cwd:      filepath.Join(hqPath, "deacon"),
			wantRole: "deacon",
		},
		{
			name:     "witness from witness dir",
			cwd:      filepath.Join(hqPath, rigName, "witness"),
			wantRole: "witness",
			wantRig:  rigName,
		},
		{
			name:     "refinery from refinery/rig dir",
			cwd:      filepath.Join(hqPath, rigName, "refinery", "rig"),
			wantRole: "refinery",
			wantRig:  rigName,
		},
		{
			name:       "polecat from polecats/Toast/rig dir",
			cwd:        filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
			wantRole:   "polecat",
			wantRig:    rigName,
			wantWorker: "Toast",
		},
		{
			name:       "crew from crew/worker1/rig dir",
			cwd:        filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
			wantRole:   "crew",
			wantRig:    rigName,
			wantWorker: "worker1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, "role", "detect")
			cmd.Dir = tt.cwd
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("gt role detect failed: %v\nOutput: %s", err, output)
			}

			got := string(output)

			// Check role is detected
			if !strings.Contains(got, tt.wantRole) {
				t.Errorf("output missing role %q\ngot: %s", tt.wantRole, got)
			}

			// Check "(from cwd)" marker
			if !strings.Contains(got, "(from cwd)") {
				t.Errorf("output missing '(from cwd)' marker\ngot: %s", got)
			}

			// Check rig if expected
			if tt.wantRig != "" {
				if !strings.Contains(got, "Rig: "+tt.wantRig) {
					t.Errorf("output missing 'Rig: %s'\ngot: %s", tt.wantRig, got)
				}
			}

			// Check worker if expected
			if tt.wantWorker != "" {
				if !strings.Contains(got, "Worker: "+tt.wantWorker) {
					t.Errorf("output missing 'Worker: %s'\ngot: %s", tt.wantWorker, got)
				}
			}
		})
	}
}

// TestRoleDetectIgnoresGTRole validates wt role detect ignores WT_ROLE env var.
func TestRoleDetectIgnoresGTRole(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	// Run from mayor dir but set WT_ROLE to deacon
	cmd = exec.Command(gtBinary, "role", "detect")
	cmd.Dir = filepath.Join(hqPath, "mayor")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir, "WT_ROLE=deacon")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("gt role detect failed: %v\nOutput: %s", err, output)
	}

	got := string(output)

	// Should detect mayor from cwd, not deacon from env
	if !strings.Contains(got, "mayor") {
		t.Errorf("should detect 'mayor' from cwd, got: %s", got)
	}

	// Should show mismatch warning
	if !strings.Contains(got, "Mismatch") {
		t.Errorf("should show mismatch warning when WT_ROLE disagrees\ngot: %s", got)
	}

	if !strings.Contains(got, "WT_ROLE=deacon") {
		t.Errorf("should show conflicting WT_ROLE value\ngot: %s", got)
	}
}

// TestRoleDetectInvalidPaths validates detection behavior for incomplete/invalid paths.
func TestRoleDetectInvalidPaths(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	rigName := "testrig"

	// Create incomplete directory structures
	dirs := []string{
		filepath.Join(hqPath, rigName),                        // rig root
		filepath.Join(hqPath, rigName, "polecats"),            // polecats without name
		filepath.Join(hqPath, rigName, "crew"),                // crew without name
		filepath.Join(hqPath, rigName, "refinery"),            // refinery without /rig
		filepath.Join(hqPath, rigName, "witness"),             // witness (valid - no /rig needed)
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	tests := []struct {
		name     string
		cwd      string
		wantRole string
	}{
		{
			name:     "rig root returns unknown",
			cwd:      filepath.Join(hqPath, rigName),
			wantRole: "unknown",
		},
		{
			name:     "polecats without name returns unknown",
			cwd:      filepath.Join(hqPath, rigName, "polecats"),
			wantRole: "unknown",
		},
		{
			name:     "crew without name returns unknown",
			cwd:      filepath.Join(hqPath, rigName, "crew"),
			wantRole: "unknown",
		},
		{
			name:     "refinery without /rig still detects refinery",
			cwd:      filepath.Join(hqPath, rigName, "refinery"),
			wantRole: "refinery",
		},
		{
			name:     "witness without /rig detects witness",
			cwd:      filepath.Join(hqPath, rigName, "witness"),
			wantRole: "witness",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, "role", "detect")
			cmd.Dir = tt.cwd
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("gt role detect failed: %v\nOutput: %s", err, output)
			}

			got := string(output)
			if !strings.Contains(got, tt.wantRole) {
				t.Errorf("expected role %q\ngot: %s", tt.wantRole, got)
			}
		})
	}
}

// TestRoleEnvIncompleteEnvVars validates wt role env fills gaps from cwd with warning.
func TestRoleEnvIncompleteEnvVars(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	rigName := "testrig"

	// Create rig directory structure for cwd detection
	dirs := []string{
		filepath.Join(hqPath, rigName, "witness"),
		filepath.Join(hqPath, rigName, "refinery", "rig"),
		filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
		filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	tests := []struct {
		name       string
		cwd        string
		envVars    []string
		wantExport []string // Expected exports in stdout
		wantStderr string   // Expected warning in stderr
	}{
		{
			name: "WT_ROLE=witness without WT_RIG, filled from cwd",
			cwd:  filepath.Join(hqPath, rigName, "witness"),
			envVars: []string{"WT_ROLE=witness"},
			wantExport: []string{
				"export WT_ROLE=witness",
				"export WT_RIG=" + rigName,
				"export BD_ACTOR=" + rigName + "/witness",
			},
			wantStderr: "env vars incomplete",
		},
		{
			name: "WT_ROLE=refinery without WT_RIG, filled from cwd",
			cwd:  filepath.Join(hqPath, rigName, "refinery", "rig"),
			envVars: []string{"WT_ROLE=refinery"},
			wantExport: []string{
				"export WT_ROLE=refinery",
				"export WT_RIG=" + rigName,
				"export BD_ACTOR=" + rigName + "/refinery",
			},
			wantStderr: "env vars incomplete",
		},
		{
			name: "WT_ROLE=polecat without WT_RIG or WT_POLECAT, filled from cwd",
			cwd:  filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
			envVars: []string{"WT_ROLE=polecat"},
			wantExport: []string{
				"export WT_ROLE=polecat",
				"export WT_RIG=" + rigName,
				"export WT_POLECAT=Toast",
				"export BD_ACTOR=" + rigName + "/polecats/Toast",
			},
			wantStderr: "env vars incomplete",
		},
		{
			name: "WT_ROLE=polecat with WT_RIG but no WT_POLECAT, filled from cwd",
			cwd:  filepath.Join(hqPath, rigName, "polecats", "Toast", "rig"),
			envVars: []string{"WT_ROLE=polecat", "WT_RIG=" + rigName},
			wantExport: []string{
				"export WT_ROLE=polecat",
				"export WT_RIG=" + rigName,
				"export WT_POLECAT=Toast",
				"export BD_ACTOR=" + rigName + "/polecats/Toast",
			},
			wantStderr: "env vars incomplete",
		},
		{
			name: "WT_ROLE=crew without WT_RIG or WT_CREW, filled from cwd",
			cwd:  filepath.Join(hqPath, rigName, "crew", "worker1", "rig"),
			envVars: []string{"WT_ROLE=crew"},
			wantExport: []string{
				"export WT_ROLE=crew",
				"export WT_RIG=" + rigName,
				"export WT_CREW=worker1",
				"export BD_ACTOR=" + rigName + "/crew/worker1",
			},
			wantStderr: "env vars incomplete",
		},
		{
			name: "Complete env vars - no warning",
			cwd:  filepath.Join(hqPath, rigName, "witness"),
			envVars: []string{"WT_ROLE=witness", "WT_RIG=" + rigName},
			wantExport: []string{
				"export WT_ROLE=witness",
				"export WT_RIG=" + rigName,
				"export BD_ACTOR=" + rigName + "/witness",
			},
			wantStderr: "", // No warning expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, "role", "env")
			cmd.Dir = tt.cwd
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
			cmd.Env = append(cmd.Env, tt.envVars...)

			// Use CombinedOutput to see stderr for debugging, but separate stdout/stderr
			stdout, _ := cmd.Output() // Only stdout
			// Re-run to get stderr
			cmd2 := exec.Command(gtBinary, "role", "env")
			cmd2.Dir = tt.cwd
			cmd2.Env = append(cleanGTEnv(), "HOME="+tmpDir)
			cmd2.Env = append(cmd2.Env, tt.envVars...)
			combined, _ := cmd2.CombinedOutput()
			stderr := strings.TrimPrefix(string(combined), string(stdout))

			// Check expected exports in stdout
			gotStdout := string(stdout)
			for _, w := range tt.wantExport {
				if !strings.Contains(gotStdout, w) {
					t.Errorf("stdout missing %q\ngot: %s", w, gotStdout)
				}
			}

			// Check expected warning in stderr
			if tt.wantStderr != "" {
				if !strings.Contains(stderr, tt.wantStderr) {
					t.Errorf("stderr should contain %q\ngot: %s\ncombined: %s", tt.wantStderr, stderr, combined)
				}
			} else {
				if strings.Contains(stderr, "incomplete") {
					t.Errorf("stderr should not contain 'incomplete' warning\ngot: %s", stderr)
				}
			}
		})
	}
}

// TestRoleEnvCwdMismatchFromIncompleteDir validates warnings when in incomplete directories.
func TestRoleEnvCwdMismatchFromIncompleteDir(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	rigName := "testrig"

	// Create incomplete directory structures (missing /rig)
	dirs := []string{
		filepath.Join(hqPath, rigName, "refinery"),            // refinery without /rig
		filepath.Join(hqPath, rigName, "polecats", "Toast"),   // polecat without /rig
		filepath.Join(hqPath, rigName, "crew", "worker1"),     // crew without /rig
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	tests := []struct {
		name       string
		cwd        string
		envVars    []string
		wantStderr string // Expected warning about cwd mismatch
	}{
		{
			name: "refinery without /rig shows cwd mismatch",
			cwd:  filepath.Join(hqPath, rigName, "refinery"),
			envVars: []string{"WT_ROLE=refinery", "WT_RIG=" + rigName},
			wantStderr: "cwd",
		},
		{
			name: "polecat without /rig shows cwd mismatch",
			cwd:  filepath.Join(hqPath, rigName, "polecats", "Toast"),
			envVars: []string{"WT_ROLE=polecat", "WT_RIG=" + rigName, "WT_POLECAT=Toast"},
			wantStderr: "cwd",
		},
		{
			name: "crew without /rig shows cwd mismatch",
			cwd:  filepath.Join(hqPath, rigName, "crew", "worker1"),
			envVars: []string{"WT_ROLE=crew", "WT_RIG=" + rigName, "WT_CREW=worker1"},
			wantStderr: "cwd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, "role", "env")
			cmd.Dir = tt.cwd
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
			cmd.Env = append(cmd.Env, tt.envVars...)

			combined, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("gt role env failed: %v\nOutput: %s", err, combined)
			}

			// Check for cwd mismatch warning
			if !strings.Contains(string(combined), tt.wantStderr) {
				t.Errorf("output should contain %q warning\ngot: %s", tt.wantStderr, combined)
			}
		})
	}
}

// TestRoleHomeInvalidPaths validates that commands fail gracefully for incomplete paths.
func TestRoleHomeInvalidPaths(t *testing.T) {
	tmpDir := t.TempDir()
	hqPath := filepath.Join(tmpDir, "test-hq")
	gtBinary := buildGT(t)

	cmd := exec.Command(gtBinary, "install", hqPath, "--no-beads")
	cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("wt install failed: %v\nOutput: %s", err, output)
	}

	rigName := "testrig"

	// Create incomplete directory structures
	dirs := []string{
		filepath.Join(hqPath, rigName),
		filepath.Join(hqPath, rigName, "polecats"),
		filepath.Join(hqPath, rigName, "crew"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}

	tests := []struct {
		name      string
		cwd       string
		shouldErr bool
	}{
		{
			name:      "rig root fails",
			cwd:       filepath.Join(hqPath, rigName),
			shouldErr: true,
		},
		{
			name:      "polecats without name fails",
			cwd:       filepath.Join(hqPath, rigName, "polecats"),
			shouldErr: true,
		},
		{
			name:      "crew without name fails",
			cwd:       filepath.Join(hqPath, rigName, "crew"),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(gtBinary, "role", "home")
			cmd.Dir = tt.cwd
			cmd.Env = append(cleanGTEnv(), "HOME="+tmpDir)

			_, err := cmd.CombinedOutput()
			if tt.shouldErr && err == nil {
				t.Errorf("expected error but command succeeded")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("expected success but got error: %v", err)
			}
		})
	}
}

