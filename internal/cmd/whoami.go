package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/speaker20/whaletown/internal/config"
	"github.com/speaker20/whaletown/internal/style"
	"github.com/speaker20/whaletown/internal/workspace"
)

var whoamiCmd = &cobra.Command{
	Use:     "whoami",
	GroupID: GroupDiag,
	Short:   "Show current identity for mail commands",
	Long: `Show the identity that will be used for mail commands.

Identity is determined by:
1. WT_ROLE env var (if set) - indicates an agent session
2. No WT_ROLE - you are the overseer (human)

Use --identity flag with mail commands to override.

Examples:
  wt whoami                      # Show current identity
  wt mail inbox                  # Check inbox for current identity
  wt mail inbox --identity mayor/  # Check Mayor's inbox instead`,
	RunE: runWhoami,
}

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func runWhoami(cmd *cobra.Command, args []string) error {
	// Get current identity using same logic as mail commands
	identity := detectSender()

	fmt.Printf("%s %s\n", style.Bold.Render("Identity:"), identity)

	// Show how it was determined
	gtRole := os.Getenv("WT_ROLE")
	if gtRole != "" {
		fmt.Printf("%s WT_ROLE=%s\n", style.Dim.Render("Source:"), gtRole)

		// Show additional env vars if present
		if rig := os.Getenv("WT_RIG"); rig != "" {
			fmt.Printf("%s WT_RIG=%s\n", style.Dim.Render("       "), rig)
		}
		if polecat := os.Getenv("WT_POLECAT"); polecat != "" {
			fmt.Printf("%s WT_POLECAT=%s\n", style.Dim.Render("       "), polecat)
		}
		if crew := os.Getenv("WT_CREW"); crew != "" {
			fmt.Printf("%s WT_CREW=%s\n", style.Dim.Render("       "), crew)
		}
	} else {
		fmt.Printf("%s no WT_ROLE set (human at terminal)\n", style.Dim.Render("Source:"))

		// If overseer, show their configured identity
		if identity == "overseer" {
			townRoot, err := workspace.FindFromCwd()
			if err == nil && townRoot != "" {
				if overseerConfig, err := config.LoadOverseerConfig(config.OverseerConfigPath(townRoot)); err == nil {
					fmt.Printf("\n%s\n", style.Bold.Render("Overseer Identity:"))
					fmt.Printf("  Name:  %s\n", overseerConfig.Name)
					if overseerConfig.Email != "" {
						fmt.Printf("  Email: %s\n", overseerConfig.Email)
					}
					if overseerConfig.Username != "" {
						fmt.Printf("  User:  %s\n", overseerConfig.Username)
					}
					fmt.Printf("  %s %s\n", style.Dim.Render("(detected via"), style.Dim.Render(overseerConfig.Source+")"))
				}
			}
		}
	}

	return nil
}
