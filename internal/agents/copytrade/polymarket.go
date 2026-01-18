// Package copytrade provides wallet tracking and copy trading agents.
package copytrade

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/speaker20/whaletown/internal/agents/common"
)

// PolymarketTracker monitors Polymarket for whale positions.
type PolymarketTracker struct {
	config  *common.Config
	wallets []common.TrackedWallet
	client  *http.Client
}

// NewPolymarketTracker creates a new Polymarket position tracker.
func NewPolymarketTracker(config *common.Config, wallets []common.TrackedWallet) *PolymarketTracker {
	// Filter to only Polymarket wallets
	polyWallets := []common.TrackedWallet{}
	for _, w := range wallets {
		if w.Platform == "polymarket" {
			polyWallets = append(polyWallets, w)
		}
	}

	return &PolymarketTracker{
		config:  config,
		wallets: polyWallets,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PolymarketMarket represents a prediction market.
type PolymarketMarket struct {
	ID       string  `json:"id"`
	Question string  `json:"question"`
	OutcomeA string  `json:"outcome_a"`
	OutcomeB string  `json:"outcome_b"`
	PriceA   float64 `json:"price_a"`
	PriceB   float64 `json:"price_b"`
	Volume   float64 `json:"volume"`
	EndDate  string  `json:"end_date"`
}

// PolymarketPosition represents a user's position on a market.
type PolymarketPosition struct {
	MarketID     string  `json:"market_id"`
	Outcome      string  `json:"outcome"`
	Shares       float64 `json:"shares"`
	AvgPrice     float64 `json:"avg_price"`
	CurrentValue float64 `json:"current_value"`
}

// FetchRecentBets fetches recent betting activity for tracked wallets.
func (t *PolymarketTracker) FetchRecentBets() ([]common.PredictionBet, error) {
	// Polymarket's public API is limited; this uses mock data for demo
	// In production, you'd need to scrape or use their GraphQL API
	return t.mockBets(), nil
}

// FetchTrendingMarkets fetches currently popular markets.
func (t *PolymarketTracker) FetchTrendingMarkets() ([]PolymarketMarket, error) {
	url := fmt.Sprintf("%s/markets?limit=10&order=volume", t.config.PolymarketBaseURL)

	resp, err := t.client.Get(url)
	if err != nil {
		return t.mockMarkets(), nil // Fallback to mock
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return t.mockMarkets(), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return t.mockMarkets(), nil
	}

	var markets []PolymarketMarket
	if err := json.Unmarshal(body, &markets); err != nil {
		return t.mockMarkets(), nil
	}

	return markets, nil
}

// mockBets returns demo prediction bets.
func (t *PolymarketTracker) mockBets() []common.PredictionBet {
	now := time.Now()
	return []common.PredictionBet{
		{
			Timestamp:   now.Add(-30 * time.Minute),
			Wallet:      "0x1234...5678",
			WalletAlias: "Poly Prophet",
			Market:      "Will BTC reach $150k by June 2026?",
			Outcome:     "YES",
			Amount:      5000.0,
			Platform:    "polymarket",
		},
		{
			Timestamp:   now.Add(-2 * time.Hour),
			Wallet:      "0x1234...5678",
			WalletAlias: "Poly Prophet",
			Market:      "Will ETH flip BTC in 2026?",
			Outcome:     "NO",
			Amount:      2500.0,
			Platform:    "polymarket",
		},
		{
			Timestamp:   now.Add(-4 * time.Hour),
			Wallet:      "0xabcd...efgh",
			WalletAlias: "Election Expert",
			Market:      "Who will win 2028 US Election?",
			Outcome:     "Republican",
			Amount:      10000.0,
			Platform:    "polymarket",
		},
	}
}

// mockMarkets returns demo trending markets.
func (t *PolymarketTracker) mockMarkets() []PolymarketMarket {
	return []PolymarketMarket{
		{
			ID:       "btc-150k-2026",
			Question: "Will BTC reach $150k by June 2026?",
			OutcomeA: "YES",
			OutcomeB: "NO",
			PriceA:   0.65,
			PriceB:   0.35,
			Volume:   2500000,
		},
		{
			ID:       "eth-flip-btc-2026",
			Question: "Will ETH market cap exceed BTC in 2026?",
			OutcomeA: "YES",
			OutcomeB: "NO",
			PriceA:   0.12,
			PriceB:   0.88,
			Volume:   1200000,
		},
		{
			ID:       "solana-top-3",
			Question: "Will Solana be top 3 by market cap EOY 2026?",
			OutcomeA: "YES",
			OutcomeB: "NO",
			PriceA:   0.45,
			PriceB:   0.55,
			Volume:   800000,
		},
	}
}

// GetTrackedWallets returns the list of tracked wallets.
func (t *PolymarketTracker) GetTrackedWallets() []common.TrackedWallet {
	return t.wallets
}
