package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/speaker20/whaletown/internal/agents/common"
	"github.com/speaker20/whaletown/internal/trader"
	"github.com/speaker20/whaletown/internal/web"
	"github.com/spf13/cobra"
)

var (
	dashboardPort       int
	dashboardOpen       bool
	dashboardWithAgents bool
)

var dashboardCmd = &cobra.Command{
	Use:     "dashboard",
	GroupID: GroupDiag,
	Short:   "Start the convoy tracking web dashboard",
	Long: `Start a web server that displays the convoy tracking dashboard.

The dashboard shows real-time convoy status with:
- Convoy list with status indicators
- Progress tracking for each convoy
- Last activity indicator (green/yellow/red)
- Auto-refresh every 30 seconds via htmx

Trading agents can be auto-started with --with-agents flag.

Example:
  wt dashboard              # Start on default port 8080
  wt dashboard --port 3000  # Start on port 3000
  wt dashboard --open       # Start and open browser
  wt dashboard --with-agents # Also start trading agents`,
	RunE: runDashboard,
}

func init() {
	dashboardCmd.Flags().IntVar(&dashboardPort, "port", 8080, "HTTP port to listen on")
	dashboardCmd.Flags().BoolVar(&dashboardOpen, "open", false, "Open browser automatically")
	dashboardCmd.Flags().BoolVar(&dashboardWithAgents, "with-agents", false, "Auto-start trading agents (researcher + copytrade)")
	rootCmd.AddCommand(dashboardCmd)
}

func runDashboard(cmd *cobra.Command, args []string) error {
	// Try to create a live fetcher (may fail if not in workspace)
	var fetcher web.ConvoyFetcher
	var onTrade func(common.Trade)

	liveFetcher, err := web.NewLiveConvoyFetcher()
	if err != nil {
		// Not in a workspace - use demo fetcher with sample whale data
		demoFetcher := web.NewDemoConvoyFetcher()
		fetcher = demoFetcher
		onTrade = demoFetcher.AddTrade
	} else {
		fetcher = liveFetcher
	}

	// Auto-start trading agents if requested or if HELIUS_API_KEY is set
	if dashboardWithAgents || os.Getenv("HELIUS_API_KEY") != "" {
		startTradingAgents(onTrade)
	}

	// Create the handler
	handler, err := web.NewConvoyHandler(fetcher)
	if err != nil {
		return fmt.Errorf("creating convoy handler: %w", err)
	}

	// Build the URL
	url := fmt.Sprintf("http://localhost:%d", dashboardPort)

	// Open browser if requested
	if dashboardOpen {
		go openBrowser(url)
	}

	// Start the server with timeouts
	fmt.Printf("🚚 Whale Town Dashboard starting at %s\n", url)
	fmt.Printf("   Press Ctrl+C to stop\n")

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", dashboardPort),
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	return server.ListenAndServe()
}

// startTradingAgents starts the researcher and copytrade agents.
func startTradingAgents(onTrade func(common.Trade)) {
	mgr := trader.NewManager()

	// Hook up callback
	mgr.OnTrade = onTrade

	// Start researcher (discovers wallets)
	if err := mgr.Start(trader.AgentTypeResearcher); err != nil {
		fmt.Printf("⚠️  Failed to start researcher: %v\n", err)
	} else {
		fmt.Println("🔬 Researcher agent started (wallet discovery)")
	}

	// Start copytrade (tracks whale trades)
	if err := mgr.Start(trader.AgentTypeCopyTrade); err != nil {
		fmt.Printf("⚠️  Failed to start copytrade: %v\n", err)
	} else {
		fmt.Println("📈 Copy Trade agent started (tracking whales)")
	}
}

// openBrowser opens the specified URL in the default browser.
func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return
	}
	_ = cmd.Start()
}
