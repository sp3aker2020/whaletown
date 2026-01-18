// Package researcher provides wallet discovery and analysis.
package researcher

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/speaker20/whaletown/internal/agents/common"
)

// Watchlist represents the shared wallet watchlist.
type Watchlist struct {
	UpdatedAt time.Time     `json:"updated_at"`
	Wallets   []WalletEntry `json:"wallets"`
}

// WalletEntry represents a wallet in the watchlist.
type WalletEntry struct {
	Address  string  `json:"address"`
	Alias    string  `json:"alias"`
	Score    int     `json:"score"`     // 0-100 score based on performance
	Profit7d float64 `json:"profit_7d"` // Profit in last 7 days (USD)
	WinRate  float64 `json:"win_rate"`  // 0.0-1.0
	Trades   int     `json:"trades"`    // Number of trades tracked
	Source   string  `json:"source"`    // "dune", "nansen", "manual"
	Platform string  `json:"platform"`  // "solana", "polymarket"
}

// WatchlistPath returns the path to the watchlist file.
func WatchlistPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".whaletown/watchlist.json"
	}
	return filepath.Join(home, ".whaletown", "watchlist.json")
}

// LoadWatchlist loads the watchlist from disk.
func LoadWatchlist() (*Watchlist, error) {
	path := WatchlistPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No watchlist yet
		}
		return nil, err
	}

	var wl Watchlist
	if err := json.Unmarshal(data, &wl); err != nil {
		return nil, err
	}

	return &wl, nil
}

// SaveWatchlist saves the watchlist to disk.
func SaveWatchlist(wl *Watchlist) error {
	path := WatchlistPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(wl, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ToTrackedWallets converts watchlist to common.TrackedWallet slice.
func (wl *Watchlist) ToTrackedWallets() []common.TrackedWallet {
	result := make([]common.TrackedWallet, len(wl.Wallets))
	for i, w := range wl.Wallets {
		result[i] = common.TrackedWallet{
			Address:  w.Address,
			Alias:    w.Alias,
			Platform: w.Platform,
			Notes:    fmt.Sprintf("Score: %d, Win rate: %.0f%%", w.Score, w.WinRate*100),
		}
	}
	return result
}

// Researcher discovers and scores profitable wallets.
type Researcher struct {
	stopCh   chan struct{}
	interval time.Duration
	OnUpdate func(*Watchlist) // Callback when watchlist updates
}

// NewResearcher creates a new researcher agent.
func NewResearcher(interval time.Duration) *Researcher {
	return &Researcher{
		stopCh:   make(chan struct{}),
		interval: interval,
	}
}

// Start begins the discovery loop.
func (r *Researcher) Start() {
	// Initial discovery
	r.discover()

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-r.stopCh:
			return
		case <-ticker.C:
			r.discover()
		}
	}
}

// Stop stops the researcher.
func (r *Researcher) Stop() {
	close(r.stopCh)
}

// discover runs a discovery cycle.
func (r *Researcher) discover() {
	fmt.Println("🔬 Researcher: Scanning for profitable wallets...")

	// For now, use curated list from research
	// TODO: Integrate with Dune Analytics API
	wallets := r.getKnownProfitableWallets()

	// Score and rank wallets
	sort.Slice(wallets, func(i, j int) bool {
		return wallets[i].Score > wallets[j].Score
	})

	// Create watchlist
	wl := &Watchlist{
		UpdatedAt: time.Now(),
		Wallets:   wallets,
	}

	// Save to disk
	if err := SaveWatchlist(wl); err != nil {
		fmt.Printf("⚠️  Failed to save watchlist: %v\n", err)
		return
	}

	fmt.Printf("✅ Researcher: Updated watchlist with %d wallets\n", len(wallets))

	if r.OnUpdate != nil {
		r.OnUpdate(wl)
	}
}

// getKnownProfitableWallets returns curated list of known profitable wallets.
// In production, this would query Dune Analytics or Nansen APIs.
func (r *Researcher) getKnownProfitableWallets() []WalletEntry {
	return []WalletEntry{
		{
			Address:  "6kbwsSY4hL6WVadLRLnWV2irkMN2AvFZVAS8McKJmAtJ",
			Alias:    "Consistent Winner",
			Score:    92,
			Profit7d: 150000,
			WinRate:  0.52,
			Trades:   98,
			Source:   "nansen",
			Platform: "solana",
		},
		{
			Address:  "5fWkLJfoDsRAaXhPJcJY19qNtDDQ5h6q1SPzsAPRrUNG",
			Alias:    "Memecoin Master",
			Score:    88,
			Profit7d: 220000,
			WinRate:  0.58,
			Trades:   205,
			Source:   "nansen",
			Platform: "solana",
		},
		{
			Address:  "9HCTuTPEiQvkUtLmTZvK6uch4E3pDynwJTbNw6jLhp9z",
			Alias:    "TRUMP Whale",
			Score:    85,
			Profit7d: 480000,
			WinRate:  0.49,
			Trades:   45,
			Source:   "dune",
			Platform: "solana",
		},
		// Additional wallets to discover
		{
			Address:  "4ETAJ4ZLARUj6xnQrVTfCLpHKS7yVe2ySCz3q9Rny8Lo",
			Alias:    "High Stakes Trader",
			Score:    82,
			Profit7d: 95000,
			WinRate:  0.97, // Very high win rate but fewer trades
			Trades:   23,
			Source:   "dune",
			Platform: "solana",
		},
	}
}
