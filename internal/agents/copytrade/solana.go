// Package copytrade provides wallet tracking and copy trading agents.
package copytrade

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/speaker20/whaletown/internal/agents/common"
)

// CacheDuration is how long to cache API results to avoid rate limits.
const CacheDuration = 60 * time.Second

// SolanaTracker monitors Solana whale wallets for swap transactions.
type SolanaTracker struct {
	config  *common.Config
	wallets []common.TrackedWallet
	client  *http.Client

	// Cache to avoid rate limits
	cacheMu    sync.RWMutex
	cachedData []common.Trade
	cacheTime  time.Time
}

// NewSolanaTracker creates a new Solana wallet tracker.
func NewSolanaTracker(config *common.Config, wallets []common.TrackedWallet) *SolanaTracker {
	// Filter to only Solana wallets
	solanaWallets := []common.TrackedWallet{}
	for _, w := range wallets {
		if w.Platform == "solana" {
			solanaWallets = append(solanaWallets, w)
		}
	}

	return &SolanaTracker{
		config:  config,
		wallets: solanaWallets,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// HeliusTransaction represents a parsed transaction from Helius API.
type HeliusTransaction struct {
	Signature       string                 `json:"signature"`
	Timestamp       int64                  `json:"timestamp"`
	Type            string                 `json:"type"`
	Description     string                 `json:"description"`
	Source          string                 `json:"source"`
	Fee             int                    `json:"fee"`
	TokenTransfers  []HeliusTokenTransfer  `json:"tokenTransfers"`
	NativeTransfers []HeliusNativeTransfer `json:"nativeTransfers"`
}

// HeliusTokenTransfer represents a token transfer in a transaction.
type HeliusTokenTransfer struct {
	FromUserAccount string  `json:"fromUserAccount"`
	ToUserAccount   string  `json:"toUserAccount"`
	TokenAmount     float64 `json:"tokenAmount"`
	Mint            string  `json:"mint"`
	TokenName       string  `json:"tokenName,omitempty"`
	TokenSymbol     string  `json:"tokenSymbol,omitempty"`
}

// HeliusNativeTransfer represents a SOL transfer.
type HeliusNativeTransfer struct {
	FromUserAccount string `json:"fromUserAccount"`
	ToUserAccount   string `json:"toUserAccount"`
	Amount          int64  `json:"amount"`
}

// FetchRecentTrades fetches recent swap transactions for all tracked wallets.
// Results are cached for 60 seconds to avoid rate limits.
func (t *SolanaTracker) FetchRecentTrades() ([]common.Trade, error) {
	// Check cache first
	t.cacheMu.RLock()
	if time.Since(t.cacheTime) < CacheDuration && len(t.cachedData) > 0 {
		cached := t.cachedData
		t.cacheMu.RUnlock()
		return cached, nil
	}
	t.cacheMu.RUnlock()

	// Fetch fresh data
	allTrades := []common.Trade{}

	for i, wallet := range t.wallets {
		// Add delay between requests to avoid rate limits
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		trades, err := t.fetchWalletTrades(wallet)
		if err != nil {
			// Log but continue with other wallets
			fmt.Printf("Error fetching trades for %s: %v\n", wallet.Alias, err)
			continue
		}
		allTrades = append(allTrades, trades...)
	}

	// Update cache
	t.cacheMu.Lock()
	t.cachedData = allTrades
	t.cacheTime = time.Now()
	t.cacheMu.Unlock()

	return allTrades, nil
}

// fetchWalletTrades fetches transactions for a single wallet.
func (t *SolanaTracker) fetchWalletTrades(wallet common.TrackedWallet) ([]common.Trade, error) {
	// Use Helius parsed transaction history API
	url := fmt.Sprintf(
		"https://api.helius.xyz/v0/addresses/%s/transactions?api-key=%s&limit=10",
		wallet.Address,
		t.config.HeliusAPIKey,
	)

	// If no API key, use a mock/demo mode
	if t.config.HeliusAPIKey == "" {
		return t.mockTrades(wallet), nil
	}

	resp, err := t.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Helius API error: %s - %s", resp.Status, string(body))
	}

	var txns []HeliusTransaction
	if err := json.NewDecoder(resp.Body).Decode(&txns); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return t.parseTrades(wallet, txns), nil
}

// parseTrades converts Helius transactions to our Trade format.
func (t *SolanaTracker) parseTrades(wallet common.TrackedWallet, txns []HeliusTransaction) []common.Trade {
	trades := []common.Trade{}

	for _, tx := range txns {
		// Only interested in SWAP transactions
		if tx.Type != "SWAP" && tx.Type != "TRANSFER" {
			continue
		}

		// Parse token transfers to determine what was swapped
		if len(tx.TokenTransfers) >= 2 {
			// Typically: first is token sent (sold), second is token received (bought)
			tokenIn := tx.TokenTransfers[0]
			tokenOut := tx.TokenTransfers[1]

			trade := common.Trade{
				Timestamp:   time.Unix(tx.Timestamp, 0),
				Wallet:      wallet.Address,
				WalletAlias: wallet.Alias,
				Type:        "swap",
				TokenIn:     tokenIn.TokenSymbol,
				TokenOut:    tokenOut.TokenSymbol,
				AmountIn:    tokenIn.TokenAmount,
				AmountOut:   tokenOut.TokenAmount,
				TxHash:      tx.Signature,
				Platform:    "solana",
			}

			// If symbols are empty, use mint addresses
			if trade.TokenIn == "" {
				trade.TokenIn = shortenAddress(tokenIn.Mint)
			}
			if trade.TokenOut == "" {
				trade.TokenOut = shortenAddress(tokenOut.Mint)
			}

			trades = append(trades, trade)
		}
	}

	return trades
}

// mockTrades returns demo trades when no API key is configured.
func (t *SolanaTracker) mockTrades(wallet common.TrackedWallet) []common.Trade {
	now := time.Now()
	return []common.Trade{
		{
			Timestamp:   now.Add(-5 * time.Minute),
			Wallet:      wallet.Address,
			WalletAlias: wallet.Alias,
			Type:        "swap",
			TokenIn:     "SOL",
			TokenOut:    "BONK",
			AmountIn:    50.0,
			AmountOut:   2500000000.0,
			TxHash:      "5abc...mock1",
			Platform:    "solana",
		},
		{
			Timestamp:   now.Add(-15 * time.Minute),
			Wallet:      wallet.Address,
			WalletAlias: wallet.Alias,
			Type:        "swap",
			TokenIn:     "USDC",
			TokenOut:    "JUP",
			AmountIn:    1000.0,
			AmountOut:   1250.0,
			TxHash:      "7def...mock2",
			Platform:    "solana",
		},
		{
			Timestamp:   now.Add(-1 * time.Hour),
			Wallet:      wallet.Address,
			WalletAlias: wallet.Alias,
			Type:        "swap",
			TokenIn:     "SOL",
			TokenOut:    "WIF",
			AmountIn:    100.0,
			AmountOut:   500.0,
			TxHash:      "9ghi...mock3",
			Platform:    "solana",
		},
	}
}

// shortenAddress returns a shortened version of an address.
func shortenAddress(addr string) string {
	if len(addr) <= 8 {
		return addr
	}
	return addr[:4] + "..." + addr[len(addr)-4:]
}

// GetTrackedWallets returns the list of tracked wallets.
func (t *SolanaTracker) GetTrackedWallets() []common.TrackedWallet {
	return t.wallets
}
