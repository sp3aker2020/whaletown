package web

// EmptyConvoyFetcher returns empty data when not in a Whale Town workspace.
// This allows the dashboard to run on Render and show the empty state.
type EmptyConvoyFetcher struct{}

// NewEmptyConvoyFetcher creates a fetcher that returns empty data.
func NewEmptyConvoyFetcher() *EmptyConvoyFetcher {
	return &EmptyConvoyFetcher{}
}

// FetchConvoys returns an empty convoy list.
func (f *EmptyConvoyFetcher) FetchConvoys() ([]ConvoyRow, error) {
	return nil, nil
}

// FetchMergeQueue returns an empty merge queue.
func (f *EmptyConvoyFetcher) FetchMergeQueue() ([]MergeQueueRow, error) {
	return nil, nil
}

// FetchPolecats returns an empty polecat list.
func (f *EmptyConvoyFetcher) FetchPolecats() ([]PolecatRow, error) {
	return nil, nil
}
