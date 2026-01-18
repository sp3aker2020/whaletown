// Package trader provides trading agent management.
package trader

import (
	"fmt"
	"sync"
	"time"

	"github.com/speaker20/whaletown/internal/agents/common"
	"github.com/speaker20/whaletown/internal/agents/copytrade"
)

// AgentType represents a type of trading agent.
type AgentType string

const (
	AgentTypeCopyTrade  AgentType = "copytrade"
	AgentTypeResearcher AgentType = "researcher"
)

// AgentStatus represents the status of a trading agent.
type AgentStatus struct {
	Name      string    `json:"name"`
	Type      AgentType `json:"type"`
	Running   bool      `json:"running"`
	StartedAt time.Time `json:"started_at,omitempty"`
	Trades    int       `json:"trades,omitempty"`  // Number of trades tracked
	Signals   int       `json:"signals,omitempty"` // Number of signals generated
}

// Manager manages trading agent lifecycles.
type Manager struct {
	mu     sync.RWMutex
	agents map[string]*runningAgent
	config *common.Config
}

type runningAgent struct {
	status      AgentStatus
	stopCh      chan struct{}
	tracker     *copytrade.SolanaTracker
	polyTracker *copytrade.PolymarketTracker
}

// NewManager creates a new trading agent manager.
func NewManager() *Manager {
	return &Manager{
		agents: make(map[string]*runningAgent),
		config: common.DefaultConfig(),
	}
}

// Start starts a trading agent.
func (m *Manager) Start(agentType AgentType) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := string(agentType)
	if _, exists := m.agents[name]; exists {
		return fmt.Errorf("agent %s is already running", name)
	}

	agent := &runningAgent{
		status: AgentStatus{
			Name:      name,
			Type:      agentType,
			Running:   true,
			StartedAt: time.Now(),
		},
		stopCh: make(chan struct{}),
	}

	switch agentType {
	case AgentTypeCopyTrade:
		wallets := common.DefaultTrackedWallets()
		agent.tracker = copytrade.NewSolanaTracker(m.config, wallets)
		agent.polyTracker = copytrade.NewPolymarketTracker(m.config, wallets)
		go m.runCopyTradeLoop(agent)

	case AgentTypeResearcher:
		// TODO: Implement researcher agent
		go m.runResearcherLoop(agent)

	default:
		return fmt.Errorf("unknown agent type: %s", agentType)
	}

	m.agents[name] = agent
	return nil
}

// Stop stops a trading agent.
func (m *Manager) Stop(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	agent, exists := m.agents[name]
	if !exists {
		return fmt.Errorf("agent %s is not running", name)
	}

	close(agent.stopCh)
	delete(m.agents, name)
	return nil
}

// List returns status of all running agents.
func (m *Manager) List() []AgentStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]AgentStatus, 0, len(m.agents))
	for _, agent := range m.agents {
		result = append(result, agent.status)
	}
	return result
}

// GetAgent returns a specific agent if running.
func (m *Manager) GetAgent(name string) (*runningAgent, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	agent, exists := m.agents[name]
	return agent, exists
}

// runCopyTradeLoop runs the copy trade agent loop.
func (m *Manager) runCopyTradeLoop(agent *runningAgent) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Initial fetch
	m.fetchTrades(agent)

	for {
		select {
		case <-agent.stopCh:
			return
		case <-ticker.C:
			m.fetchTrades(agent)
		}
	}
}

// fetchTrades fetches trades from trackers.
func (m *Manager) fetchTrades(agent *runningAgent) {
	if agent.tracker != nil {
		trades, err := agent.tracker.FetchRecentTrades()
		if err == nil {
			m.mu.Lock()
			agent.status.Trades = len(trades)
			m.mu.Unlock()
		}
	}
}

// runResearcherLoop runs the researcher agent loop (placeholder).
func (m *Manager) runResearcherLoop(agent *runningAgent) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-agent.stopCh:
			return
		case <-ticker.C:
			// TODO: Generate signals
			m.mu.Lock()
			agent.status.Signals++
			m.mu.Unlock()
		}
	}
}

// FetchLatestTrades returns latest trades from the copy trade agent.
func (m *Manager) FetchLatestTrades() ([]common.Trade, error) {
	m.mu.RLock()
	agent, exists := m.agents["copytrade"]
	m.mu.RUnlock()

	if !exists || agent.tracker == nil {
		return nil, fmt.Errorf("copytrade agent not running")
	}

	return agent.tracker.FetchRecentTrades()
}
