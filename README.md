# go-bsky-feed-generator
A minimal implementation of a BlueSky Feed Generator in Go


## Requirements

To run this feed generator, all you need is `docker` with `docker-compose`.

## Running

Start up the feed generator by running: `make up`

This will build the feed generator service binary inside a docker container and stand up the service on your machine at port `9032`.

## Accessing

This service exposes the following routes:

- `/.well-known/did.json`
  - This route is used by ATProto to verify ownership of the DID the service is claiming, it's a static JSON document.
  - You can see how this is generated in `pkg/feed-generator/endpoints.go:GetWellKnownDID()`
- `/xrpc/app.bsky.feed.getFeedSkeleton`
  - This route is what clients call to generate a feed page, it includes three query parameters for feed generation: `feed`, `cursor`, and `limit`
  - You can see how those are parsed and handled in `pkg/feed-generator/endpoints.go:GetFeedSkeleton()`
- `/xrpc/app.bsky.feed.describeFeedGenerator`
  - This route is how the service advertises which feeds it supports to clients.
  - You can see how those are parsed and handled in `pkg/feed-generator/endpoints.go:DescribeFeedGenerator()`


## Architecture

This repo is structured to abstract away a `Feed` interface that allows for you to add all sorts of feeds to the generator.

These feeds can be simple static feeds like the `pkg/static-feed/feed.go` implementation, or they can be much more complex feeds that draw on different data sources and filter them in cool ways to produce pages of feed items.

The `Feed` interface is defined by any struct implementing two functions:

``` go
type Feed interface {
	GetPage(ctx context.Context, limit int64, cursor string) (feedPosts []*appbsky.FeedDefs_SkeletonFeedPost, newCursor *string, err error)
	Describe(ctx context.Context) (*appbsky.FeedDescribeFeedGenerator_Feed, error)
}
```
