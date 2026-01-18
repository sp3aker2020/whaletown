package web

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/speaker20/whaletown/internal/agents/common"
	"github.com/speaker20/whaletown/internal/agents/copytrade"
	"github.com/speaker20/whaletown/internal/agents/researcher"
)

// DemoConvoyFetcher returns sample whale-themed data for demo/showcase purposes.
// It also fetches LIVE whale trades when HELIUS_API_KEY is set.
type DemoConvoyFetcher struct {
	solanaTracker *copytrade.SolanaTracker
	startTime     time.Time

	// Real-time trades from WebSocket
	mu             sync.RWMutex
	realtimeTrades []common.Trade
}

// NewDemoConvoyFetcher creates a demo fetcher with sample data.
func NewDemoConvoyFetcher() *DemoConvoyFetcher {
	config := common.DefaultConfig()
	wallets := loadWalletsFromWatchlist()

	return &DemoConvoyFetcher{
		solanaTracker:  copytrade.NewSolanaTracker(config, wallets),
		startTime:      time.Now(),
		realtimeTrades: make([]common.Trade, 0),
	}
}

// AddTrade adds a real-time trade to the buffer.
func (f *DemoConvoyFetcher) AddTrade(trade common.Trade) {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Prepend new trade
	f.realtimeTrades = append([]common.Trade{trade}, f.realtimeTrades...)

	// Keep last 50
	if len(f.realtimeTrades) > 50 {
		f.realtimeTrades = f.realtimeTrades[:50]
	}
}

// loadWalletsFromWatchlist loads wallets from watchlist or defaults.
func loadWalletsFromWatchlist() []common.TrackedWallet {
	wl, err := researcher.LoadWatchlist()
	if err == nil && wl != nil && len(wl.Wallets) > 0 {
		return wl.ToTrackedWallets()
	}
	return common.DefaultTrackedWallets()
}

// FetchConvoys returns empty list (legacy demo data removed).
func (f *DemoConvoyFetcher) FetchConvoys() ([]ConvoyRow, error) {
	return []ConvoyRow{}, nil
}

// FetchMergeQueue returns empty list (legacy demo data removed).
func (f *DemoConvoyFetcher) FetchMergeQueue() ([]MergeQueueRow, error) {
	return []MergeQueueRow{}, nil
}

// FetchPolecats returns empty list (legacy demo data removed).
func (f *DemoConvoyFetcher) FetchPolecats() ([]PolecatRow, error) {
	return []PolecatRow{}, nil
}

// FetchWhaleTrades fetches live whale trades from the Solana tracker + WebSocket.
func (f *DemoConvoyFetcher) FetchWhaleTrades() ([]WhaleTradeRow, error) {
	// Get real-time trades first
	f.mu.RLock()
	trades := make([]common.Trade, len(f.realtimeTrades))
	copy(trades, f.realtimeTrades)
	f.mu.RUnlock()

	apiKey := os.Getenv("HELIUS_API_KEY")

	// If we have API key, also try to fetch historical/recent from REST
	if apiKey != "" {
		if apiTrades, err := f.solanaTracker.FetchRecentTrades(); err == nil {
			// Append API trades (ignoring duplicates ideally, but for now simple append)
			trades = append(trades, apiTrades...)
		}
	}

	// Logic: If we have NO data (WS or API), fall back to mocks ONLY if no key
	if len(trades) == 0 && apiKey == "" {
		return f.mockWhaleTrades(), nil
	}

	// Convert to WhaleTradeRow format
	rows := make([]WhaleTradeRow, 0, len(trades))
	for _, t := range trades {
		// Default type for WS alerts
		txType := t.Type
		if txType == "alert" {
			txType = "Detected 🚨"
		}

		txURL := fmt.Sprintf("https://solscan.io/tx/%s", t.TxHash)
		if t.TxHash == "" {
			txURL = "#"
		}

		rows = append(rows, WhaleTradeRow{
			Timestamp:   formatTimeAgo(t.Timestamp),
			WalletAlias: t.WalletAlias,
			Type:        txType,
			TokenIn:     t.TokenIn,
			TokenOut:    t.TokenOut,
			AmountIn:    formatAmount(t.AmountIn),
			AmountOut:   formatAmount(t.AmountOut),
			TxHash:      shortenTx(t.TxHash),
			TxURL:       txURL,
			Platform:    t.Platform,
		})
	}

	// Sort by time (newest first) implicitly or ensure it
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

// FetchAgentStatuses returns the status of trading agents.
func (f *DemoConvoyFetcher) FetchAgentStatuses() ([]AgentStatusRow, error) {
	// Calculate time since start for realistic display
	elapsed := time.Since(f.startTime)
	researcherNext := 5*time.Minute - (elapsed % (5 * time.Minute))
	copytradeNext := 30*time.Second - (elapsed % (30 * time.Second))

	// Count trades and wallets
	trades, _ := f.solanaTracker.FetchRecentTrades()
	tradeCount := len(trades)

	walletCount := len(f.solanaTracker.GetTrackedWallets())

	return []AgentStatusRow{
		{
			Name:        "researcher",
			DisplayName: "🔬 Researcher",
			Status:      "running",
			StatusClass: "agent-running",
			LastRun:     formatTimeAgo(f.startTime),
			NextRun:     formatDuration(researcherNext),
			ItemCount:   walletCount,
			ItemLabel:   "wallets",
		},
		{
			Name:        "copytrade",
			DisplayName: "📈 Copy Trade",
			Status:      "running",
			StatusClass: "agent-running",
			LastRun:     formatTimeAgo(f.startTime.Add(elapsed - (elapsed % (30 * time.Second)))),
			NextRun:     formatDuration(copytradeNext),
			ItemCount:   tradeCount,
			ItemLabel:   "trades",
		},
		// Upcoming Agents (Roadmap)
		{
			Name:        "sniper",
			DisplayName: "🎯 Sniper Agent",
			Status:      "coming soon",
			StatusClass: "agent-soon",
			LastRun:     "dev mode",
			NextRun:     "TBA",
			ItemCount:   0,
			ItemLabel:   "memes",
		},
		{
			Name:        "cda",
			DisplayName: "⚖️ CDA Strategy",
			Status:      "coming soon",
			StatusClass: "agent-soon",
			LastRun:     "planned",
			NextRun:     "TBA",
			ItemCount:   0,
			ItemLabel:   "arbs",
		},
		{
			Name:        "safety",
			DisplayName: "🛡️ Safety Checker",
			Status:      "coming soon",
			StatusClass: "agent-soon",
			LastRun:     "planned",
			NextRun:     "TBA",
			ItemCount:   0,
			ItemLabel:   "audits",
		},
		{
			Name:        "sentiment",
			DisplayName: "🧠 Sentiment Analysis",
			Status:      "coming soon",
			StatusClass: "agent-soon",
			LastRun:     "planned",
			NextRun:     "TBA",
			ItemCount:   0,
			ItemLabel:   "signals",
		},
	}, nil
}

// FetchTrackedWallets returns wallets from the researcher watchlist.
func (f *DemoConvoyFetcher) FetchTrackedWallets() ([]TrackedWalletRow, error) {
	wl, err := researcher.LoadWatchlist()
	if err != nil || wl == nil || len(wl.Wallets) == 0 {
		// Return default wallets if no watchlist
		defaults := common.DefaultTrackedWallets()
		rows := make([]TrackedWalletRow, len(defaults))
		for i, w := range defaults {
			rows[i] = TrackedWalletRow{
				Address:    shortenAddress(w.Address),
				Alias:      w.Alias,
				Score:      80,
				ScoreClass: "score-high",
				Profit7d:   "N/A",
				WinRate:    "N/A",
				Trades:     0,
				Platform:   w.Platform,
			}
		}
		return rows, nil
	}

	// Convert watchlist to rows
	rows := make([]TrackedWalletRow, len(wl.Wallets))
	for i, w := range wl.Wallets {
		scoreClass := "score-low"
		if w.Score >= 85 {
			scoreClass = "score-high"
		} else if w.Score >= 70 {
			scoreClass = "score-medium"
		}

		rows[i] = TrackedWalletRow{
			Address:    shortenAddress(w.Address),
			Alias:      w.Alias,
			Score:      w.Score,
			ScoreClass: scoreClass,
			Profit7d:   formatAmount(w.Profit7d),
			WinRate:    fmt.Sprintf("%.0f%%", w.WinRate*100),
			Trades:     w.Trades,
			Platform:   w.Platform,
		}
	}
	return rows, nil
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("in %ds", int(d.Seconds()))
	}
	return fmt.Sprintf("in %dm", int(d.Minutes()))
}

func shortenAddress(addr string) string {
	if len(addr) <= 12 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}
