package web

import (
	"fmt"
	"os"
	"time"

	"github.com/speaker20/whaletown/internal/activity"
	"github.com/speaker20/whaletown/internal/agents/common"
	"github.com/speaker20/whaletown/internal/agents/copytrade"
)

// DemoConvoyFetcher returns sample whale-themed data for demo/showcase purposes.
// It also fetches LIVE whale trades when HELIUS_API_KEY is set.
type DemoConvoyFetcher struct {
	solanaTracker *copytrade.SolanaTracker
}

// NewDemoConvoyFetcher creates a demo fetcher with sample data.
func NewDemoConvoyFetcher() *DemoConvoyFetcher {
	config := common.DefaultConfig()
	wallets := common.DefaultTrackedWallets()

	return &DemoConvoyFetcher{
		solanaTracker: copytrade.NewSolanaTracker(config, wallets),
	}
}

// FetchConvoys returns sample whale-themed convoy data.
func (f *DemoConvoyFetcher) FetchConvoys() ([]ConvoyRow, error) {
	return []ConvoyRow{
		{
			ID:         "hq-cv-trade1",
			Title:      "🐋 High-Frequency Arbitrage",
			Status:     "open",
			Progress:   "3/5",
			Completed:  3,
			Total:      5,
			WorkStatus: "active",
			LastActivity: activity.Info{
				FormattedAge: "20ms ago",
				ColorClass:   activity.ColorGreen,
			},
			TrackedIssues: []TrackedIssue{
				{ID: "wt-scan-dex", Title: "Scan DEX liquidity pools", Status: "closed", Assignee: "market-maker/polecats/Alpha"},
				{ID: "wt-calc-spread", Title: "Calculate spread opportunities", Status: "closed", Assignee: "market-maker/polecats/Beta"},
				{ID: "wt-exec-flash", Title: "Execute flash loan", Status: "hooked", Assignee: "market-maker/polecats/Gamma"},
				{ID: "wt-rebalance", Title: "Rebalance assets", Status: "open", Assignee: ""},
				{ID: "wt-log-profit", Title: "Log trade profit", Status: "open", Assignee: ""},
			},
		},
		{
			ID:         "hq-cv-trade2",
			Title:      "🫧 Meme Coin Sniper",
			Status:     "open",
			Progress:   "1/3",
			Completed:  1,
			Total:      3,
			WorkStatus: "active",
			LastActivity: activity.Info{
				FormattedAge: "5s ago",
				ColorClass:   activity.ColorGreen,
			},
			TrackedIssues: []TrackedIssue{
				{ID: "wt-monitor-new", Title: "Monitor new pair listings", Status: "closed", Assignee: "sniper-lab/polecats/Doge"},
				{ID: "wt-check-honeypot", Title: "Verify contract (anti-rug)", Status: "hooked", Assignee: "sniper-lab/polecats/Shiba"},
				{ID: "wt-buy-entry", Title: "Execute entry buy", Status: "open", Assignee: ""},
			},
		},
		{
			ID:         "hq-cv-trade3",
			Title:      "🌊 Portfolio Rebalancer",
			Status:     "open",
			Progress:   "0/4",
			Completed:  0,
			Total:      4,
			WorkStatus: "waiting",
			LastActivity: activity.Info{
				FormattedAge: "1h ago",
				ColorClass:   activity.ColorYellow,
			},
			TrackedIssues: []TrackedIssue{
				{ID: "wt-fetch-balances", Title: "Fetch wallet balances", Status: "open", Assignee: ""},
				{ID: "wt-calc-allocation", Title: "Calculate target allocation", Status: "open", Assignee: ""},
				{ID: "wt-gen-swaps", Title: "Generate swap instructions", Status: "open", Assignee: ""},
				{ID: "wt-exec-batch", Title: "Execute batch transaction", Status: "open", Assignee: ""},
			},
		},
	}, nil
}

// FetchMergeQueue returns sample merge queue data.
func (f *DemoConvoyFetcher) FetchMergeQueue() ([]MergeQueueRow, error) {
	return []MergeQueueRow{
		{
			Number:     101,
			Repo:       "strategies",
			Title:      "feat: Add Solana MEV protection",
			URL:        "https://github.com/sp3aker2020/strategies/pull/101",
			CIStatus:   "pass",
			Mergeable:  "ready",
			ColorClass: "mq-green",
		},
		{
			Number:     102,
			Repo:       "strategies",
			Title:      "fix: Slippage calculation overflow",
			URL:        "https://github.com/sp3aker2020/strategies/pull/102",
			CIStatus:   "pending",
			Mergeable:  "ready",
			ColorClass: "mq-yellow",
		},
	}, nil
}

// FetchPolecats returns sample polecat/worker data.
func (f *DemoConvoyFetcher) FetchPolecats() ([]PolecatRow, error) {
	return []PolecatRow{
		{
			Name:      "Alpha",
			Rig:       "market-maker",
			SessionID: "wt-market-maker-Alpha",
			LastActivity: activity.Info{
				FormattedAge: "100ms ago",
				ColorClass:   activity.ColorGreen,
			},
			StatusHint: "Scanning ETH/USDC pools...",
		},
		{
			Name:      "Beta",
			Rig:       "market-maker",
			SessionID: "wt-market-maker-Beta",
			LastActivity: activity.Info{
				FormattedAge: "200ms ago",
				ColorClass:   activity.ColorGreen,
			},
			StatusHint: "Analyzing spread...",
		},
		{
			Name:      "Doge",
			Rig:       "sniper-lab",
			SessionID: "wt-sniper-lab-Doge",
			LastActivity: activity.Info{
				FormattedAge: "1s ago",
				ColorClass:   activity.ColorGreen,
			},
			StatusHint: "Watching mempool...",
		},
		{
			Name:      "refinery",
			Rig:       "whaletown",
			SessionID: "wt-whaletown-refinery",
			LastActivity: activity.Info{
				FormattedAge: "10s ago",
				ColorClass:   activity.ColorGreen,
			},
			StatusHint: "Backtesting strategy #42",
		},
	}, nil
}

// FetchWhaleTrades fetches live whale trades from the Solana tracker.
func (f *DemoConvoyFetcher) FetchWhaleTrades() ([]WhaleTradeRow, error) {
	apiKey := os.Getenv("HELIUS_API_KEY")
	if apiKey == "" {
		// Return mock data if no API key
		return f.mockWhaleTrades(), nil
	}

	// Fetch real trades
	trades, err := f.solanaTracker.FetchRecentTrades()
	if err != nil {
		return f.mockWhaleTrades(), nil
	}

	// Convert to WhaleTradeRow format
	rows := make([]WhaleTradeRow, 0, len(trades))
	for _, t := range trades {
		rows = append(rows, WhaleTradeRow{
			Timestamp:   formatTimeAgo(t.Timestamp),
			WalletAlias: t.WalletAlias,
			Type:        t.Type,
			TokenIn:     t.TokenIn,
			TokenOut:    t.TokenOut,
			AmountIn:    formatAmount(t.AmountIn),
			AmountOut:   formatAmount(t.AmountOut),
			TxHash:      shortenTx(t.TxHash),
			TxURL:       fmt.Sprintf("https://solscan.io/tx/%s", t.TxHash),
			Platform:    t.Platform,
		})
	}

	return rows, nil
}

// mockWhaleTrades returns demo trades when no API key is set.
func (f *DemoConvoyFetcher) mockWhaleTrades() []WhaleTradeRow {
	return []WhaleTradeRow{
		{
			Timestamp:   "2m ago",
			WalletAlias: "Memecoin Master",
			Type:        "swap",
			TokenIn:     "SOL",
			TokenOut:    "BONK",
			AmountIn:    "50.00",
			AmountOut:   "2.5B",
			TxHash:      "5abc...mock",
			TxURL:       "https://solscan.io/tx/mock1",
			Platform:    "solana",
		},
		{
			Timestamp:   "15m ago",
			WalletAlias: "TRUMP Whale",
			Type:        "swap",
			TokenIn:     "USDC",
			TokenOut:    "TRUMP",
			AmountIn:    "10,000",
			AmountOut:   "5,000",
			TxHash:      "7def...mock",
			TxURL:       "https://solscan.io/tx/mock2",
			Platform:    "solana",
		},
	}
}

// Helper functions
func formatTimeAgo(t time.Time) string {
	diff := time.Since(t)
	if diff < time.Minute {
		return fmt.Sprintf("%ds ago", int(diff.Seconds()))
	} else if diff < time.Hour {
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}
	return t.Format("Jan 2")
}

func formatAmount(amt float64) string {
	if amt >= 1_000_000_000 {
		return fmt.Sprintf("%.1fB", amt/1_000_000_000)
	} else if amt >= 1_000_000 {
		return fmt.Sprintf("%.1fM", amt/1_000_000)
	} else if amt >= 1_000 {
		return fmt.Sprintf("%.1fK", amt/1_000)
	} else if amt >= 1 {
		return fmt.Sprintf("%.2f", amt)
	}
	return fmt.Sprintf("%.4f", amt)
}

func shortenTx(tx string) string {
	if len(tx) <= 12 {
		return tx
	}
	return tx[:4] + "..." + tx[len(tx)-4:]
}
