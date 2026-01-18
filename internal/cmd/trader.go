package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/speaker20/whaletown/internal/trader"
	"github.com/spf13/cobra"
)

// Global manager instance (persists for the process lifetime)
var traderManager *trader.Manager

func init() {
	traderManager = trader.NewManager()
}

var traderCmd = &cobra.Command{
	Use:   "trader",
	Short: "Manage trading agents",
	Long: `Manage Whale Town trading agents for copy trading and research.

Available agents:
  copytrade   - Monitors whale wallets and copies their trades
  researcher  - Generates buy/sell signals from sentiment analysis

Examples:
  wt trader start copytrade    # Start the copy trade agent
  wt trader stop copytrade     # Stop the agent
  wt trader list               # List running agents
  wt trader status             # Show current trades/signals`,
}

var traderStartCmd = &cobra.Command{
	Use:   "start <agent>",
	Short: "Start a trading agent",
	Long: `Start a trading agent. Available agents:
  copytrade   - Monitors whale wallets via Helius API
  researcher  - Generates signals from sentiment (coming soon)

The agent runs in the background and continuously monitors for trades.
Set HELIUS_API_KEY environment variable for live Solana data.`,
	Args: cobra.ExactArgs(1),
	RunE: runTraderStart,
}

var traderStopCmd = &cobra.Command{
	Use:   "stop <agent>",
	Short: "Stop a trading agent",
	Args:  cobra.ExactArgs(1),
	RunE:  runTraderStop,
}

var traderListCmd = &cobra.Command{
	Use:   "list",
	Short: "List running trading agents",
	RunE:  runTraderList,
}

var traderStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show trading status (trades, signals)",
	RunE:  runTraderStatus,
}

var (
	traderJSON bool
)

func init() {
	rootCmd.AddCommand(traderCmd)
	traderCmd.AddCommand(traderStartCmd)
	traderCmd.AddCommand(traderStopCmd)
	traderCmd.AddCommand(traderListCmd)
	traderCmd.AddCommand(traderStatusCmd)

	traderListCmd.Flags().BoolVar(&traderJSON, "json", false, "Output as JSON")
	traderStatusCmd.Flags().BoolVar(&traderJSON, "json", false, "Output as JSON")
}

func runTraderStart(cmd *cobra.Command, args []string) error {
	agentName := args[0]

	var agentType trader.AgentType
	switch agentName {
	case "copytrade":
		agentType = trader.AgentTypeCopyTrade
	case "researcher":
		agentType = trader.AgentTypeResearcher
	default:
		return fmt.Errorf("unknown agent: %s (available: copytrade, researcher)", agentName)
	}

	if err := traderManager.Start(agentType); err != nil {
		return err
	}

	fmt.Printf("🐋 Started %s agent\n", agentName)

	if agentName == "copytrade" {
		if os.Getenv("HELIUS_API_KEY") == "" {
			fmt.Println("⚠️  HELIUS_API_KEY not set - using mock data")
		} else {
			fmt.Println("✅ Connected to Helius API for live Solana data")
		}
	}

	fmt.Println("\n📊 Agent running... Press Ctrl+C to stop")
	fmt.Println("   Use 'wt dashboard' in another terminal to see trades\n")

	// Block and show periodic status
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sigCh:
			fmt.Println("\n🛑 Stopping agent...")
			traderManager.Stop(agentName)
			return nil
		case <-ticker.C:
			agents := traderManager.List()
			for _, a := range agents {
				if a.Name == agentName {
					fmt.Printf("   📈 %s: %d trades tracked\n", a.Name, a.Trades)
				}
			}
		}
	}
}

func runTraderStop(cmd *cobra.Command, args []string) error {
	agentName := args[0]

	if err := traderManager.Stop(agentName); err != nil {
		return err
	}

	fmt.Printf("🛑 Stopped %s agent\n", agentName)
	return nil
}

func runTraderList(cmd *cobra.Command, args []string) error {
	agents := traderManager.List()

	if traderJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(agents)
	}

	if len(agents) == 0 {
		fmt.Println("No trading agents running")
		fmt.Println("\nStart one with: wt trader start copytrade")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tTYPE\tSTATUS\tTRADES\tSIGNALS\tSTARTED")
	for _, a := range agents {
		status := "stopped"
		if a.Running {
			status = "running"
		}
		started := ""
		if !a.StartedAt.IsZero() {
			started = a.StartedAt.Format("15:04:05")
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\t%s\n",
			a.Name, a.Type, status, a.Trades, a.Signals, started)
	}
	return w.Flush()
}

func runTraderStatus(cmd *cobra.Command, args []string) error {
	agents := traderManager.List()

	if len(agents) == 0 {
		fmt.Println("No trading agents running")
		return nil
	}

	// Try to get latest trades from copytrade agent
	trades, err := traderManager.FetchLatestTrades()
	if err == nil && len(trades) > 0 {
		fmt.Println("🐋 Recent Whale Trades:")
		fmt.Println()

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TIME\tWHALE\tACTION\tIN\tOUT")

		maxTrades := 10
		if len(trades) < maxTrades {
			maxTrades = len(trades)
		}

		for i := 0; i < maxTrades; i++ {
			t := trades[i]
			fmt.Fprintf(w, "%s\t%s\t%s\t%.4f %s\t%.4f %s\n",
				t.Timestamp.Format("15:04:05"),
				t.WalletAlias,
				t.Type,
				t.AmountIn, t.TokenIn,
				t.AmountOut, t.TokenOut)
		}
		w.Flush()
	} else {
		fmt.Println("No trades available (agent may be starting up)")
	}

	return nil
}
