// ABOUTME: Command to disable Whale Town system-wide.
// ABOUTME: Sets the global state to disabled so tools work vanilla.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/speaker20/whaletown/internal/shell"
	"github.com/speaker20/whaletown/internal/state"
	"github.com/speaker20/whaletown/internal/style"
)

var disableClean bool

var disableCmd = &cobra.Command{
	Use:     "disable",
	GroupID: GroupConfig,
	Short:   "Disable Whale Town system-wide",
	Long: `Disable Whale Town for all agentic coding tools.

When disabled:
  - Shell hooks become no-ops
  - Claude Code SessionStart hooks skip 'wt prime'
  - Tools work 100% vanilla (no Whale Town behavior)

The workspace (~/gt) is preserved. Use 'wt enable' to re-enable.

Flags:
  --clean  Also remove shell integration from ~/.zshrc/~/.bashrc

Environment overrides still work:
  GASTOWN_ENABLED=1   - Enable for current session only`,
	RunE: runDisable,
}

func init() {
	disableCmd.Flags().BoolVar(&disableClean, "clean", false,
		"Remove shell integration from RC files")
	rootCmd.AddCommand(disableCmd)
}

func runDisable(cmd *cobra.Command, args []string) error {
	if err := state.Disable(); err != nil {
		return fmt.Errorf("disabling Whale Town: %w", err)
	}

	if disableClean {
		if err := removeShellIntegration(); err != nil {
			fmt.Printf("%s Could not clean shell integration: %v\n",
				style.Warning.Render("!"), err)
		} else {
			fmt.Println("  Removed shell integration from RC files")
		}
	}

	fmt.Printf("%s Whale Town disabled\n", style.Success.Render("✓"))
	fmt.Println()
	fmt.Println("All agentic coding tools now work vanilla.")
	if !disableClean {
		fmt.Printf("Use %s to also remove shell hooks\n",
			style.Dim.Render("gt disable --clean"))
	}
	fmt.Printf("Use %s to re-enable\n", style.Dim.Render("gt enable"))

	return nil
}

func removeShellIntegration() error {
	return shell.Remove()
}
