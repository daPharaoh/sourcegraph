package backend

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/zoekt"
	"github.com/google/zoekt/query"
	"github.com/google/zoekt/rpc"
	zoektstream "github.com/google/zoekt/stream"

	"github.com/sourcegraph/sourcegraph/internal/trace/ot"
)

var zoektHTTPClient = &http.Client{
	Transport: &ot.Transport{
		RoundTripper: http.DefaultTransport,
	},
}

// ZoektStreamFunc is a convenience function to create a stream receiver from a
// function.
type ZoektStreamFunc func(*zoekt.SearchResult)

func (f ZoektStreamFunc) Send(event *zoekt.SearchResult) {
	f(event)
}

// StreamSearchEvent has fields optionally set representing events that happen
// during a search.
//
// This is a Sourcegraph extension.
type StreamSearchEvent struct {
	// SearchResult is non-nil if this event is a search result. These should be
	// combined with previous and later SearchResults.
	SearchResult *zoekt.SearchResult
}

// StreamSearchAdapter adapts a zoekt.Searcher to conform to the StreamSearch
// interface by calling zoekt.Searcher.Search.
type StreamSearchAdapter struct {
	zoekt.Searcher
}

func (s *StreamSearchAdapter) StreamSearch(ctx context.Context, q query.Q, opts *zoekt.SearchOptions, c zoekt.Sender) error {
	sr, err := s.Search(ctx, q, opts)
	if err != nil {
		return err
	}
	c.Send(sr)
	return nil
}

func (s *StreamSearchAdapter) String() string {
	return "streamSearchAdapter{" + s.Searcher.String() + "}"
}

func NewZoektStream(address string) zoekt.Streamer {
	addressWithScheme := address
	if !strings.HasPrefix(addressWithScheme, "http://") {
		addressWithScheme = "http://" + addressWithScheme
	}
	return &zoektStream{
		rpc.Client(address),
		zoektstream.NewClient(addressWithScheme, zoektHTTPClient),
	}
}

type zoektStream struct {
	zoekt.Searcher
	*zoektstream.Client
}
