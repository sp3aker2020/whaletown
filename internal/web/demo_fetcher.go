package web

import (
	"github.com/speaker20/whaletown/internal/activity"
)

// DemoConvoyFetcher returns sample whale-themed data for demo/showcase purposes.
// This allows the dashboard to show sample data when not in a Whale Town workspace.
type DemoConvoyFetcher struct{}

// NewDemoConvoyFetcher creates a demo fetcher with sample data.
func NewDemoConvoyFetcher() *DemoConvoyFetcher {
	return &DemoConvoyFetcher{}
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
