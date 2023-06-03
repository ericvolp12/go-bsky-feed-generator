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

func NewStaticFeed(ctx context.Context, feedActorDID string, feedName string, staticPostURIs []string) *StaticFeed {
	return &StaticFeed{
		FeedActorDID:   feedActorDID,
		FeedName:       feedName,
		StaticPostURIs: staticPostURIs,
	}
}

func (sf *StaticFeed) GetPage(ctx context.Context, userDID string, limit int64, cursor string) ([]*appbsky.FeedDefs_SkeletonFeedPost, *string, error) {
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

func (sf *StaticFeed) Describe(ctx context.Context) (*appbsky.FeedDescribeFeedGenerator_Feed, error) {
	return &appbsky.FeedDescribeFeedGenerator_Feed{
		Uri: "at://" + sf.FeedActorDID + "/app.bsky.feed.generator/" + sf.FeedName,
	}, nil
}
