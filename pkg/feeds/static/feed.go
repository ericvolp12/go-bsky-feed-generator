package static

import (
	"context"
	"fmt"
	"strconv"

	appbsky "github.com/bluesky-social/indigo/api/bsky"
)

type StaticFeed struct {
	FeedActorDID   string
	FeedName       string
	StaticPostURIs []string
}

// NewStaticFeed returns a new StaticFeed, a list of aliases for the feed, and an error
// StaticFeed is a trivial implementation of the Feed interface, so its aliases are just the input feedName
func NewStaticFeed(ctx context.Context, feedActorDID string, feedName string, staticPostURIs []string) (*StaticFeed, []string, error) {
	return &StaticFeed{
		FeedActorDID:   feedActorDID,
		FeedName:       feedName,
		StaticPostURIs: staticPostURIs,
	}, []string{feedName}, nil
}

// GetPage returns a list of FeedDefs_SkeletonFeedPost, a new cursor, and an error
// It takes a feed name, a user DID, a limit, and a cursor
// The feed name can be used to produce different feeds from the same feed generator
func (sf *StaticFeed) GetPage(ctx context.Context, feed string, userDID string, limit int64, cursor string) ([]*appbsky.FeedDefs_SkeletonFeedPost, *string, error) {
	cursorAsInt := int64(0)
	var err error

	if cursor != "" {
		cursorAsInt, err = strconv.ParseInt(cursor, 10, 64)
		if err != nil {
			return nil, nil, fmt.Errorf("cursor is not an integer: %w", err)
		}
	}

	posts := []*appbsky.FeedDefs_SkeletonFeedPost{}

	for i, postURI := range sf.StaticPostURIs {
		if int64(i) < cursorAsInt {
			continue
		}

		if int64(len(posts)) >= limit {
			break
		}

		posts = append(posts, &appbsky.FeedDefs_SkeletonFeedPost{
			Post: postURI,
		})
	}

	cursorAsInt += int64(len(posts))

	var newCursor *string

	if cursorAsInt < int64(len(sf.StaticPostURIs)) {
		newCursor = new(string)
		*newCursor = strconv.FormatInt(cursorAsInt, 10)
	}

	return posts, newCursor, nil
}

// Describe returns a list of FeedDescribeFeedGenerator_Feed, and an error
// StaticFeed is a trivial implementation of the Feed interface, so it returns a single FeedDescribeFeedGenerator_Feed
// For a more complicated feed, this function would return a list of FeedDescribeFeedGenerator_Feed with the URIs of aliases
// supported by the feed
func (sf *StaticFeed) Describe(ctx context.Context) ([]appbsky.FeedDescribeFeedGenerator_Feed, error) {
	return []appbsky.FeedDescribeFeedGenerator_Feed{
		{
			Uri: "at://" + sf.FeedActorDID + "/app.bsky.feed.generator/" + sf.FeedName,
		},
	}, nil
}
