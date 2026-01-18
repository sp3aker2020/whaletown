package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/speaker20/whaletown/internal/version"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	GroupID: GroupDiag,
	Short:   "Show Whale Town information and what's new",
	Long: `Display information about the current Whale Town installation.

This command shows:
  - Version information
  - What's new in recent versions (with --whats-new flag)

Examples:
  wt info
  wt info --whats-new
  wt info --whats-new --json`,
	Run: func(cmd *cobra.Command, args []string) {
		whatsNewFlag, _ := cmd.Flags().GetBool("whats-new")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		if whatsNewFlag {
			showWhatsNew(jsonFlag)
			return
		}

		// Default: show basic info
		info := map[string]interface{}{
			"version": Version,
			"build":   Build,
		}

		if commit := resolveCommitHash(); commit != "" {
			info["commit"] = version.ShortCommit(commit)
		}
		if branch := resolveBranch(); branch != "" {
			info["branch"] = branch
		}

		if jsonFlag {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(info)
			return
		}

		fmt.Printf("Whale Town v%s (%s)\n", Version, Build)
		if commit, ok := info["commit"].(string); ok {
			if branch, ok := info["branch"].(string); ok {
				fmt.Printf("  %s@%s\n", branch, commit)
			} else {
				fmt.Printf("  %s\n", commit)
			}
		}
		fmt.Println("\nUse 'wt info --whats-new' to see recent changes")
	},
}

// VersionChange represents agent-relevant changes for a specific version
type VersionChange struct {
	Version string   `json:"version"`
	Date    string   `json:"date"`
	Changes []string `json:"changes"`
}

// versionChanges contains agent-actionable changes for recent versions
var versionChanges = []VersionChange{
	{
		Version: "0.4.0",
		Date:    "2026-01-17",
		Changes: []string{
			"FIX: Orphan cleanup skips valid tmux sessions - Prevents false kills of witnesses/refineries/deacon during startup by checking gt-*/hq-* session membership",
		},
	},
	{
		Version: "0.3.1",
		Date:    "2026-01-17",
		Changes: []string{
			"FIX: Orphan cleanup on macOS - TTY comparison now handles macOS '??' format",
			"FIX: Session kill orphan prevention - wt done and wt crew stop use KillSessionWithProcesses",
		},
	},
	{
		Version: "0.3.0",
		Date:    "2026-01-17",
		Changes: []string{
			"NEW: wt show/cat - Inspect bead contents and metadata",
			"NEW: wt orphans list/kill - Detect and clean up orphaned Claude processes",
			"NEW: wt convoy close - Manual convoy closure command",
			"NEW: wt commit/trail - Git wrappers with bead awareness",
			"NEW: Plugin system - wt plugin run/history, wt dispatch --plugin",
			"NEW: Beads-native messaging - Queue, channel, and group beads",
			"NEW: wt mail claim - Claim messages from queues",
			"NEW: wt polecat identity show - Display CV summary",
			"NEW: whaletown-release molecule formula - Automated release workflow",
			"NEW: Parallel agent startup - Faster boot with concurrency limit",
			"NEW: Automatic orphan cleanup - Detect and kill orphaned processes",
			"NEW: Worktree setup hooks - Inject local configurations",
			"CHANGED: MR tracking via beads - Removed mrqueue package",
			"CHANGED: Desire-path commands - Agent ergonomics shortcuts",
			"CHANGED: Explicit escalation in polecat templates",
			"FIX: Kill process tree on shutdown - Prevents orphaned Claude processes",
			"FIX: Agent bead prefix alignment - Multi-hyphen IDs for consistency",
			"FIX: Idle Polecat Heresy warnings in templates",
			"FIX: Zombie session detection in doctor",
			"FIX: Windows build support with platform-specific handling",
		},
	},
	{
		Version: "0.2.0",
		Date:    "2026-01-04",
		Changes: []string{
			"NEW: Convoy Dashboard - Web UI for monitoring Whale Town (gt dashboard)",
			"NEW: Two-level beads architecture - hq-* prefix for town, rig prefixes for projects",
			"NEW: Multi-agent support with pluggable registry",
			"NEW: wt rig start/stop/restart/status - Multi-rig management commands",
			"NEW: Ephemeral polecat model - Immediate recycling after each work unit",
			"NEW: wt costs command - Session cost tracking and reporting",
			"NEW: Conflict resolution workflow for polecats with merge-slot gates",
			"NEW: wt convoy --tree and wt convoy check for cross-rig coordination",
			"NEW: Batch slinging - wt sling supports multiple beads at once",
			"NEW: spawn alias for start across all role subcommands",
			"NEW: wt mail archive supports multiple message IDs",
			"NEW: wt mail --all flag for clearing all mail",
			"NEW: Circuit breaker for stuck agents",
			"NEW: Binary age detection in wt status",
			"NEW: Shell completion installation instructions",
			"CHANGED: Handoff migrated to skills format",
			"CHANGED: Crew workers push directly to main (no PRs)",
			"CHANGED: Session names include town name",
			"FIX: Thread-safety for agent session resume",
			"FIX: Orphan daemon prevention via file locking",
			"FIX: Zombie tmux session cleanup",
			"FIX: Default branch detection (no longer hardcodes 'main')",
			"FIX: Enter key retry logic for reliable delivery",
			"FIX: Beads prefix routing for cross-rig operations",
		},
	},
	{
		Version: "0.1.1",
		Date:    "2026-01-02",
		Changes: []string{
			"FIX: Tmux keybindings scoped to Whale Town sessions only",
			"NEW: OSS project files - CHANGELOG.md, .golangci.yml, RELEASING.md",
			"NEW: Version bump script - scripts/bump-version.sh",
			"FIX: wt rig add and wt crew add CLI syntax documentation",
			"FIX: Rig prefix routing for agent beads",
			"FIX: Beads init targets correct database",
		},
	},
	{
		Version: "0.1.0",
		Date:    "2026-01-02",
		Changes: []string{
			"Initial public release of Whale Town",
			"NEW: Town structure - Hierarchical workspace with rigs, crews, and polecats",
			"NEW: Rig management - wt rig add/list/remove",
			"NEW: Crew workspaces - wt crew add for persistent developer workspaces",
			"NEW: Polecat workers - Transient agent workers managed by Witness",
			"NEW: Mayor - Global coordinator for cross-rig work",
			"NEW: Deacon - Town-level lifecycle patrol and heartbeat",
			"NEW: Witness - Per-rig polecat lifecycle manager",
			"NEW: Refinery - Merge queue processor with code review",
			"NEW: Convoy system - wt convoy create/list/status",
			"NEW: Sling workflow - wt sling <bead> <rig>",
			"NEW: Molecule workflows - Formula-based multi-step task execution",
			"NEW: Mail system - wt mail inbox/send/read",
			"NEW: Escalation protocol - wt escalate with severity levels",
			"NEW: Handoff mechanism - wt handoff for context-preserving session cycling",
			"NEW: Beads integration - Issue tracking via beads (bd commands)",
			"NEW: Tmux sessions with theming",
			"NEW: Status dashboard - wt status",
			"NEW: Activity feed - wt feed",
			"NEW: Nudge system - wt nudge for reliable message delivery",
		},
	},
}

// showWhatsNew displays agent-relevant changes from recent versions
func showWhatsNew(jsonOutput bool) {
	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(map[string]interface{}{
			"current_version": Version,
			"recent_changes":  versionChanges,
		})
		return
	}

	// Human-readable output
	fmt.Printf("\nWhat's New in Whale Town (Current: v%s)\n", Version)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	for _, vc := range versionChanges {
		// Highlight if this is the current version
		versionMarker := ""
		if vc.Version == Version {
			versionMarker = " <- current"
		}

		fmt.Printf("## v%s (%s)%s\n\n", vc.Version, vc.Date, versionMarker)

		for _, change := range vc.Changes {
			fmt.Printf("  * %s\n", change)
		}
		fmt.Println()
	}

	fmt.Println("Tip: Use 'wt info --whats-new --json' for machine-readable output")
	fmt.Println()
}

func init() {
	infoCmd.Flags().Bool("whats-new", false, "Show agent-relevant changes from recent versions")
	infoCmd.Flags().Bool("json", false, "Output in JSON format")
	rootCmd.AddCommand(infoCmd)
}
