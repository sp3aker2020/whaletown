// Package common provides shared types and utilities for trading agents.
package common

import (
	"os"
)

// Config holds API keys and configuration for agents.
type Config struct {
	HeliusAPIKey      string
	PolymarketBaseURL string
	KalshiAPIKey      string
	KalshiAPISecret   string
}

// DefaultConfig returns configuration from environment variables.
func DefaultConfig() *Config {
	return &Config{
		HeliusAPIKey:      os.Getenv("HELIUS_API_KEY"),
		PolymarketBaseURL: "https://gamma-api.polymarket.com",
		KalshiAPIKey:      os.Getenv("KALSHI_API_KEY"),
		KalshiAPISecret:   os.Getenv("KALSHI_API_SECRET"),
	}
}

// DefaultTrackedWallets returns a list of known whale wallets to track.
// These are example addresses - replace with real ones.
func DefaultTrackedWallets() []TrackedWallet {
	return []TrackedWallet{
		// Solana whales (example addresses - these are placeholders)
		{
			Address:  "5ZWj7a1f8tWkjBESHKgrLmXshuXxqeY9SYcfbshpAqPG",
			Alias:    "Memecoin Alpha",
			Platform: "solana",
			Notes:    "Known for early memecoin entries",
		},
		{
			Address:  "7xKpY5q7e8VaLJQQwPRNdPgRD9zTb3VcXmWfZqRaNkvP",
			Alias:    "DeFi Whale",
			Platform: "solana",
			Notes:    "Large DeFi positions",
		},
		// Polymarket (example - need to find real whale addresses)
		{
			Address:  "0x1234567890abcdef1234567890abcdef12345678",
			Alias:    "Poly Prophet",
			Platform: "polymarket",
			Notes:    "Top leaderboard bettor",
		},
	}
}
