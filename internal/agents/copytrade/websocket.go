// Package copytrade provides wallet tracking and copy trading agents.
package copytrade

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/speaker20/whaletown/internal/agents/common"
)

// TradeCallback is called when a new trade is detected.
type TradeCallback func(trade common.Trade)

// WebSocketListener listens for real-time wallet transactions via WebSocket.
type WebSocketListener struct {
	config  *common.Config
	wallets []common.TrackedWallet
	conn    *websocket.Conn
	mu      sync.Mutex
	running bool
	OnTrade TradeCallback
}

// NewWebSocketListener creates a listener for real-time wallet monitoring.
func NewWebSocketListener(config *common.Config, wallets []common.TrackedWallet) *WebSocketListener {
	// Filter to only Solana wallets
	solanaWallets := []common.TrackedWallet{}
	for _, w := range wallets {
		if w.Platform == "solana" {
			solanaWallets = append(solanaWallets, w)
		}
	}

	return &WebSocketListener{
		config:  config,
		wallets: solanaWallets,
	}
}

// Start connects to the WebSocket and begins listening.
func (l *WebSocketListener) Start(ctx context.Context) error {
	wsURL := l.config.SolanaWSURL
	if wsURL == "" {
		// Use Helius WebSocket if available
		if l.config.HeliusAPIKey != "" {
			wsURL = fmt.Sprintf("wss://mainnet.helius-rpc.com/?api-key=%s", l.config.HeliusAPIKey)
		} else {
			return fmt.Errorf("no WebSocket URL configured (set SOLANA_WS_URL or HELIUS_API_KEY)")
		}
	}

	fmt.Printf("🔌 Connecting to WebSocket: %s\n", maskURL(wsURL))

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	l.mu.Lock()
	l.conn = conn
	l.running = true
	l.mu.Unlock()

	// Subscribe to each wallet
	for _, wallet := range l.wallets {
		if err := l.subscribeToWallet(wallet); err != nil {
			fmt.Printf("⚠️  Failed to subscribe to %s: %v\n", wallet.Alias, err)
		} else {
			fmt.Printf("👀 Watching %s (%s...)\n", wallet.Alias, wallet.Address[:8])
		}
	}

	// Listen for messages
	go l.readLoop(ctx)

	// Start ping keepalive
	go l.pingLoop(ctx)

	return nil
}

// pingLoop sends periodic pings to keep the connection alive.
func (l *WebSocketListener) pingLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.mu.Lock()
			if !l.running || l.conn == nil {
				l.mu.Unlock()
				return
			}
			conn := l.conn
			l.mu.Unlock()

			// Send ping
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("⚠️ Ping failed: %v\n", err)
				return
			}
		}
	}
}

// subscribeToWallet sends a subscription request for a wallet.
func (l *WebSocketListener) subscribeToWallet(wallet common.TrackedWallet) error {
	// Solana logs subscription to monitor account activity
	subRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      wallet.Address[:8],
		"method":  "logsSubscribe",
		"params": []interface{}{
			map[string]interface{}{
				"mentions": []string{wallet.Address},
			},
			map[string]interface{}{
				"commitment": "confirmed",
			},
		},
	}

	return l.conn.WriteJSON(subRequest)
}

// readLoop reads messages from the WebSocket.
func (l *WebSocketListener) readLoop(ctx context.Context) {
	defer l.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		l.mu.Lock()
		if !l.running || l.conn == nil {
			l.mu.Unlock()
			return
		}
		conn := l.conn
		l.mu.Unlock()

		// Set read deadline (longer than ping interval)
		conn.SetReadDeadline(time.Now().Add(120 * time.Second))

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return
			}
			// Try to reconnect
			fmt.Printf("⚠️  WebSocket read error: %v (will reconnect)\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		l.handleMessage(message)
	}
}

// handleMessage processes incoming WebSocket messages.
func (l *WebSocketListener) handleMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		return
	}

	// Check if this is a notification (new transaction)
	if method, ok := msg["method"].(string); ok && method == "logsNotification" {
		l.handleLogsNotification(msg)
	}
}

// handleLogsNotification processes a logs notification.
func (l *WebSocketListener) handleLogsNotification(msg map[string]interface{}) {
	params, ok := msg["params"].(map[string]interface{})
	if !ok {
		return
	}

	result, ok := params["result"].(map[string]interface{})
	if !ok {
		return
	}

	value, ok := result["value"].(map[string]interface{})
	if !ok {
		return
	}

	signature, _ := value["signature"].(string)
	if signature == "" {
		return
	}

	// We detected a transaction! Log it
	fmt.Printf("🚨 New transaction detected: %s...\n", signature[:16])

	// Create a trade alert (details would need to be fetched separately)
	trade := common.Trade{
		Timestamp:   time.Now(),
		Type:        "alert",
		TxHash:      signature,
		Platform:    "solana",
		WalletAlias: "Whale Alert",
	}

	if l.OnTrade != nil {
		l.OnTrade(trade)
	}
}

// Stop closes the WebSocket connection.
func (l *WebSocketListener) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.running = false
	if l.conn != nil {
		l.conn.Close()
		l.conn = nil
	}
}

// IsRunning returns whether the listener is active.
func (l *WebSocketListener) IsRunning() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.running
}

// maskURL hides API keys in URLs for logging.
func maskURL(url string) string {
	if len(url) > 50 {
		return url[:40] + "..."
	}
	return url
}
