//go:build ignore
// +build ignore

// Quick test for the Solana tracker
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/speaker20/whaletown/internal/agents/common"
	"github.com/speaker20/whaletown/internal/agents/copytrade"
)

func main() {
	// Load API key from environment or .env file
	apiKey := os.Getenv("HELIUS_API_KEY")
	if apiKey == "" {
		apiKey = "f38d1544-b8e7-494e-a946-cb7255143abf" // Fallback for testing
	}

	config := &common.Config{
		HeliusAPIKey: apiKey,
	}

	// Use the real tracked wallets from config
	wallets := common.DefaultTrackedWallets()

	tracker := copytrade.NewSolanaTracker(config, wallets)

	fmt.Println("🐋 Testing Solana Tracker...")
	fmt.Printf("API Key: %s...%s\n", apiKey[:8], apiKey[len(apiKey)-4:])
	fmt.Printf("Tracking %d wallets\n\n", len(tracker.GetTrackedWallets()))

	trades, err := tracker.FetchRecentTrades()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d trades:\n\n", len(trades))
	for i, trade := range trades {
		fmt.Printf("%d. [%s] %s\n", i+1, trade.WalletAlias, trade.Type)
		fmt.Printf("   %s → %s\n", trade.TokenIn, trade.TokenOut)
		fmt.Printf("   Amount: %.2f → %.2f\n", trade.AmountIn, trade.AmountOut)
		fmt.Printf("   Time: %s\n", trade.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Tx: %s\n\n", trade.TxHash[:12]+"...")
	}

	// Pretty print as JSON
	jsonData, _ := json.MarshalIndent(trades, "", "  ")
	fmt.Println("Raw JSON:")
	fmt.Println(string(jsonData))
}
