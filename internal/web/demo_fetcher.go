package web

import (
	"github.com/speaker20/whaletown/internal/activity"
)

// DemoConvoyFetcher returns sample whale-themed data for demo/showcase purposes.
// This allows the dashboard to show sample data when not in a Whale Town workspace.
type DemoConvoyFetcher struct{}

// NewDemoConvoyFetcher creates a demo fetcher with sample data.
func NewDemoConvoyFetcher() *DemoConvoyFetcher {
	return &DemoConvoyFetcher{}
}

// FetchConvoys returns sample whale-themed convoy data.
func (f *DemoConvoyFetcher) FetchConvoys() ([]ConvoyRow, error) {
	return []ConvoyRow{
		{
			ID:         "hq-cv-whale1",
			Title:      "🐋 Great Migration to v2.0",
			Status:     "open",
			Progress:   "3/5",
			Completed:  3,
			Total:      5,
			WorkStatus: "active",
			LastActivity: activity.Info{
				FormattedAge: "2m ago",
				ColorClass:   activity.ColorGreen,
			},
			TrackedIssues: []TrackedIssue{
				{ID: "wt-update-deps", Title: "Update dependencies", Status: "closed", Assignee: "frontend-lab/polecats/Nemo"},
				{ID: "wt-fix-tests", Title: "Fix failing tests", Status: "closed", Assignee: "backend-lab/polecats/Orca"},
				{ID: "wt-update-docs", Title: "Update documentation", Status: "closed", Assignee: "docs-lab/polecats/Pearl"},
				{ID: "wt-security-audit", Title: "Security audit", Status: "hooked", Assignee: "backend-lab/polecats/Fin"},
				{ID: "wt-e2e-tests", Title: "End-to-end testing", Status: "open", Assignee: ""},
			},
		},
		{
			ID:         "hq-cv-whale2",
			Title:      "🫧 Bubble Net Authentication",
			Status:     "open",
			Progress:   "1/3",
			Completed:  1,
			Total:      3,
			WorkStatus: "active",
			LastActivity: activity.Info{
				FormattedAge: "5m ago",
				ColorClass:   activity.ColorGreen,
			},
			TrackedIssues: []TrackedIssue{
				{ID: "wt-jwt-tokens", Title: "Implement JWT tokens", Status: "closed", Assignee: "auth-lab/polecats/Moby"},
				{ID: "wt-oauth-flow", Title: "OAuth 2.0 flow", Status: "hooked", Assignee: "auth-lab/polecats/Splash"},
				{ID: "wt-session-mgmt", Title: "Session management", Status: "open", Assignee: ""},
			},
		},
		{
			ID:         "hq-cv-whale3",
			Title:      "🌊 Deep Dive Performance",
			Status:     "open",
			Progress:   "0/4",
			Completed:  0,
			Total:      4,
			WorkStatus: "waiting",
			LastActivity: activity.Info{
				FormattedAge: "unassigned",
				ColorClass:   activity.ColorUnknown,
			},
			TrackedIssues: []TrackedIssue{
				{ID: "wt-profiling", Title: "Performance profiling", Status: "open", Assignee: ""},
				{ID: "wt-db-queries", Title: "Optimize database queries", Status: "open", Assignee: ""},
				{ID: "wt-caching", Title: "Implement caching layer", Status: "open", Assignee: ""},
				{ID: "wt-load-tests", Title: "Load testing", Status: "open", Assignee: ""},
			},
		},
	}, nil
}

// FetchMergeQueue returns sample merge queue data.
func (f *DemoConvoyFetcher) FetchMergeQueue() ([]MergeQueueRow, error) {
	return []MergeQueueRow{
		{
			Number:     42,
			Repo:       "whaletown",
			Title:      "feat: Add bubble net tracking",
			URL:        "https://github.com/sp3aker2020/whaletown/pull/42",
			CIStatus:   "pass",
			Mergeable:  "ready",
			ColorClass: "mq-green",
		},
		{
			Number:     41,
			Repo:       "whaletown",
			Title:      "fix: Pod member session cleanup",
			URL:        "https://github.com/sp3aker2020/whaletown/pull/41",
			CIStatus:   "pending",
			Mergeable:  "ready",
			ColorClass: "mq-yellow",
		},
		{
			Number:     40,
			Repo:       "whaletown",
			Title:      "docs: Update whale lore",
			URL:        "https://github.com/sp3aker2020/whaletown/pull/40",
			CIStatus:   "pass",
			Mergeable:  "ready",
			ColorClass: "mq-green",
		},
	}, nil
}

// FetchPolecats returns sample polecat/worker data.
func (f *DemoConvoyFetcher) FetchPolecats() ([]PolecatRow, error) {
	return []PolecatRow{
		{
			Name:      "Nemo",
			Rig:       "frontend-lab",
			SessionID: "wt-frontend-lab-Nemo",
			LastActivity: activity.Info{
				FormattedAge: "30s ago",
				ColorClass:   activity.ColorGreen,
			},
			StatusHint: "Building components...",
		},
		{
			Name:      "Orca",
			Rig:       "backend-lab",
			SessionID: "wt-backend-lab-Orca",
			LastActivity: activity.Info{
				FormattedAge: "1m ago",
				ColorClass:   activity.ColorGreen,
			},
			StatusHint: "Running tests...",
		},
		{
			Name:      "Moby",
			Rig:       "auth-lab",
			SessionID: "wt-auth-lab-Moby",
			LastActivity: activity.Info{
				FormattedAge: "3m ago",
				ColorClass:   activity.ColorGreen,
			},
			StatusHint: "Implementing OAuth flow",
		},
		{
			Name:      "refinery",
			Rig:       "whaletown",
			SessionID: "wt-whaletown-refinery",
			LastActivity: activity.Info{
				FormattedAge: "10s ago",
				ColorClass:   activity.ColorGreen,
			},
			StatusHint: "Processing 3 PRs",
		},
	}, nil
}
