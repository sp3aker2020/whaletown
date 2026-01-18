// Package common provides shared types and utilities for trading agents.
package common

import "time"

// Trade represents a swap/trade transaction.
type Trade struct {
	Timestamp   time.Time `json:"timestamp"`
	Wallet      string    `json:"wallet"`
	WalletAlias string    `json:"wallet_alias,omitempty"`
	Type        string    `json:"type"` // "swap", "buy", "sell"
	TokenIn     string    `json:"token_in"`
	TokenOut    string    `json:"token_out"`
	AmountIn    float64   `json:"amount_in"`
	AmountOut   float64   `json:"amount_out"`
	TxHash      string    `json:"tx_hash"`
	Platform    string    `json:"platform"` // "solana", "polymarket", "kalshi"
}

// Signal represents a trading signal from the Researcher agent.
type Signal struct {
	Timestamp  time.Time `json:"timestamp"`
	Token      string    `json:"token"`
	Action     string    `json:"action"`     // "BUY", "SELL", "HOLD"
	Confidence float64   `json:"confidence"` // 0.0 - 1.0
	Reason     string    `json:"reason"`
	Source     string    `json:"source"` // "researcher", "copytrade"
}

// TrackedWallet represents a whale wallet being monitored.
type TrackedWallet struct {
	Address  string `json:"address"`
	Alias    string `json:"alias"`
	Platform string `json:"platform"` // "solana", "polymarket", "kalshi"
	Notes    string `json:"notes,omitempty"`
}

// PredictionBet represents a bet on a prediction market.
type PredictionBet struct {
	Timestamp   time.Time `json:"timestamp"`
	Wallet      string    `json:"wallet"`
	WalletAlias string    `json:"wallet_alias,omitempty"`
	Market      string    `json:"market"`
	Outcome     string    `json:"outcome"` // "YES", "NO"
	Amount      float64   `json:"amount"`
	Platform    string    `json:"platform"` // "polymarket", "kalshi"
}
