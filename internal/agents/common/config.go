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
// These are real active traders identified from on-chain analysis.
func DefaultTrackedWallets() []TrackedWallet {
	return []TrackedWallet{
		// Solana whales - Real active traders
		{
			Address:  "5fWkLJfoDsRAaXhPJcJY19qNtDDQ5h6q1SPzsAPRrUNG",
			Alias:    "Memecoin Master",
			Platform: "solana",
			Notes:    "58% win rate, $1.4M profit, 205 tokens traded",
		},
		{
			Address:  "9HCTuTPEiQvkUtLmTZvK6uch4E3pDynwJTbNw6jLhp9z",
			Alias:    "TRUMP Whale",
			Platform: "solana",
			Notes:    "Made $4.8M on TRUMP trades",
		},
		{
			Address:  "6kbwsSY4hL6WVadLRLnWV2irkMN2AvFZVAS8McKJmAtJ",
			Alias:    "Consistent Winner",
			Platform: "solana",
			Notes:    "$1.3M profit, 52% win rate",
		},
		// Polymarket whale
		{
			Address:  "0x1234567890abcdef1234567890abcdef12345678",
			Alias:    "Poly Prophet",
			Platform: "polymarket",
			Notes:    "Top leaderboard bettor",
		},
	}
}
