// ABOUTME: Command to enable Whale Town system-wide.
// ABOUTME: Sets the global state to enabled for all agentic coding tools.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/speaker20/whaletown/internal/state"
	"github.com/speaker20/whaletown/internal/style"
)

var enableCmd = &cobra.Command{
	Use:     "enable",
	GroupID: GroupConfig,
	Short:   "Enable Whale Town system-wide",
	Long: `Enable Whale Town for all agentic coding tools.

When enabled:
  - Shell hooks set WT_TOWN_ROOT and WT_RIG environment variables
  - Claude Code SessionStart hooks run 'wt prime' for context
  - Git repos are auto-registered as rigs (configurable)

Use 'wt disable' to turn off. Use 'wt status --global' to check state.

Environment overrides:
  GASTOWN_DISABLED=1  - Disable for current session only
  GASTOWN_ENABLED=1   - Enable for current session only`,
	RunE: runEnable,
}

func init() {
	rootCmd.AddCommand(enableCmd)
}

func runEnable(cmd *cobra.Command, args []string) error {
	if err := state.Enable(Version); err != nil {
		return fmt.Errorf("enabling Whale Town: %w", err)
	}

	fmt.Printf("%s Whale Town enabled\n", style.Success.Render("✓"))
	fmt.Println()
	fmt.Println("Whale Town will now:")
	fmt.Println("  • Inject context into Claude Code sessions")
	fmt.Println("  • Set WT_TOWN_ROOT and WT_RIG environment variables")
	fmt.Println("  • Auto-register git repos as rigs (if configured)")
	fmt.Println()
	fmt.Printf("Use %s to disable, %s to check status\n",
		style.Dim.Render("gt disable"),
		style.Dim.Render("gt status --global"))

	return nil
}
