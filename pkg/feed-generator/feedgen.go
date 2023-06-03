package feedgenerator

import (
	"context"

	appbsky "github.com/bluesky-social/indigo/api/bsky"
)

type Feed interface {
	GetPage(ctx context.Context, userDID string, limit int64, cursor string) (feedPosts []*appbsky.FeedDefs_SkeletonFeedPost, newCursor *string, err error)
	Describe(ctx context.Context) (*appbsky.FeedDescribeFeedGenerator_Feed, error)
}

type FeedGenerator struct {
	FeedActorDID          string          // DID of the Repo the Feed is published under
	ServiceEndpoint       string          // URL of the FeedGenerator service
	ServiceDID            string          // DID of the FeedGenerator service
	AcceptableURIPrefixes []string        // URIs that the FeedGenerator is allowed to generate feeds for
	Feeds                 map[string]Feed // map of FeedName to Feed
}

// NewFeedGenerator returns a new FeedGenerator
func NewFeedGenerator(
	ctx context.Context,
	feedActorDID string,
	serviceDID string,
	acceptableDIDs []string,
	serviceEndpoint string,
) (*FeedGenerator, error) {
	acceptableURIPrefixes := []string{}
	for _, did := range acceptableDIDs {
		acceptableURIPrefixes = append(acceptableURIPrefixes, "at://"+did+"/app.bsky.feed.generator/")
	}

	return &FeedGenerator{
		Feeds:                 map[string]Feed{},
		FeedActorDID:          feedActorDID,
		ServiceDID:            serviceDID,
		AcceptableURIPrefixes: acceptableURIPrefixes,
		ServiceEndpoint:       serviceEndpoint,
	}, nil
}

// AddFeed adds a feed to the FeedGenerator
func (fg *FeedGenerator) AddFeed(feedName string, feed Feed) {
	if fg.Feeds == nil {
		fg.Feeds = map[string]Feed{}
	}

	fg.Feeds[feedName] = feed
}
