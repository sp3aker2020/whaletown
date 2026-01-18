package web

import (
	"html/template"
	"net/http"
)

// ConvoyFetcher defines the interface for fetching convoy data.
type ConvoyFetcher interface {
	FetchConvoys() ([]ConvoyRow, error)
	FetchMergeQueue() ([]MergeQueueRow, error)
	FetchPolecats() ([]PolecatRow, error)
	FetchWhaleTrades() ([]WhaleTradeRow, error)
	FetchAgentStatuses() ([]AgentStatusRow, error)
	FetchTrackedWallets() ([]TrackedWalletRow, error)
}

// ConvoyHandler handles HTTP requests for the convoy dashboard.
type ConvoyHandler struct {
	fetcher  ConvoyFetcher
	template *template.Template
}

// NewConvoyHandler creates a new convoy handler with the given fetcher.
func NewConvoyHandler(fetcher ConvoyFetcher) (*ConvoyHandler, error) {
	tmpl, err := LoadTemplates()
	if err != nil {
		return nil, err
	}

	return &ConvoyHandler{
		fetcher:  fetcher,
		template: tmpl,
	}, nil
}

// ServeHTTP handles HTTP requests and routes to appropriate handlers.
func (h *ConvoyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/village":
		h.serveVillage(w, r)
	default:
		h.serveDashboard(w, r)
	}
}

// serveDashboard renders the convoy dashboard.
func (h *ConvoyHandler) serveDashboard(w http.ResponseWriter, r *http.Request) {
	// Fetch agent statuses (primary data)
	agentStatuses, _ := h.fetcher.FetchAgentStatuses()

	// Fetch tracked wallets
	trackedWallets, _ := h.fetcher.FetchTrackedWallets()

	// Fetch whale trades
	whaleTrades, _ := h.fetcher.FetchWhaleTrades()

	// Legacy data (empty for agent-focused dashboard)
	convoys, _ := h.fetcher.FetchConvoys()
	mergeQueue, _ := h.fetcher.FetchMergeQueue()
	polecats, _ := h.fetcher.FetchPolecats()

	data := ConvoyData{
		AgentStatuses:  agentStatuses,
		TrackedWallets: trackedWallets,
		WhaleTrades:    whaleTrades,
		Convoys:        convoys,
		MergeQueue:     mergeQueue,
		Polecats:       polecats,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.template.ExecuteTemplate(w, "convoy.html", data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// serveVillage renders the whale village visualization.
func (h *ConvoyHandler) serveVillage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.template.ExecuteTemplate(w, "village.html", nil); err != nil {
		http.Error(w, "Failed to render village template", http.StatusInternalServerError)
		return
	}
}
