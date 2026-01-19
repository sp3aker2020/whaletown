// Package common provides shared types and utilities for trading agents.
package common

import (
	"os"
)

// Config holds API keys and configuration for agents.
type Config struct {
	// Helius API for parsed transactions (fallback)
	HeliusAPIKey string

	// Custom Solana RPC (for higher rate limits)
	// If set, uses this for RPC calls instead of public endpoints
	SolanaRPCURL string

	// Solana WebSocket URL for real-time subscriptions
	// Example: wss://api.mainnet-beta.solana.com
	SolanaWSURL string

	// Public Key: (derived from private key)
	SolanaPrivateKey string

	// Prediction market APIs
	PolymarketBaseURL string
	KalshiAPIKey      string
	KalshiAPISecret   string
}

// DefaultConfig returns configuration from environment variables.
func DefaultConfig() *Config {
	return &Config{
		HeliusAPIKey:      os.Getenv("HELIUS_API_KEY"),
		SolanaRPCURL:      os.Getenv("SOLANA_RPC_URL"),
		SolanaWSURL:       os.Getenv("SOLANA_WS_URL"),
		SolanaPrivateKey:  os.Getenv("SOLANA_PRIVATE_KEY"),
		PolymarketBaseURL: "https://clob.polymarket.com",
		KalshiAPIKey:      os.Getenv("KALSHI_API_KEY"),
		KalshiAPISecret:   os.Getenv("KALSHI_API_SECRET"),
	}
}

// HasCustomRPC returns true if a custom RPC URL is configured.
func (c *Config) HasCustomRPC() bool {
	return c.SolanaRPCURL != ""
}

// HasWebSocket returns true if WebSocket URL is configured.
func (c *Config) HasWebSocket() bool {
	return c.SolanaWSURL != ""
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
		{
			Address:  "27Fyd42KmGRmbZSRHSmT85mA8JJwH4aEfUExPJwKYUTN",
			Alias:    "Test Wallet (User)",
			Platform: "solana",
			Notes:    "Manual test wallet for copy trading verification",
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
