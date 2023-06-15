# go-bsky-feed-generator
A minimal implementation of a BlueSky Feed Generator in Go


## Requirements

To run this feed generator, all you need is `docker` with `docker-compose`.

## Running

Start up the feed generator by running: `make up`

This will build the feed generator service binary inside a docker container and stand up the service on your machine at port `9032`.

To view a sample static feed (with only one post) go to:

- [`http://localhost:9032/xrpc/app.bsky.feed.getFeedSkeleton?feed=at://did:plc:replace-me-with-your-did/app.bsky.feed.generator/static`](http://localhost:9032/xrpc/app.bsky.feed.getFeedSkeleton?feed=at://did:plc:replace-me-with-your-did/app.bsky.feed.generator/static)

Update the variables in `.env` when you actually want to deploy the service somewhere, at which point `did:plc:replace-me-with-your-did` should be replaced with the value of `FEED_ACTOR_DID`.

## Accessing

This service exposes the following routes:

- `/.well-known/did.json`
  - This route is used by ATProto to verify ownership of the DID the service is claiming, it's a static JSON document.
  - You can see how this is generated in `pkg/gin/endpoints.go:GetWellKnownDID()`
- `/xrpc/app.bsky.feed.getFeedSkeleton`
  - This route is what clients call to generate a feed page, it includes three query parameters for feed generation: `feed`, `cursor`, and `limit`
  - You can see how those are parsed and handled in `pkg/gin/endpoints.go:GetFeedSkeleton()`
- `/xrpc/app.bsky.feed.describeFeedGenerator`
  - This route is how the service advertises which feeds it supports to clients.
  - You can see how those are parsed and handled in `pkg/gin/endpoints.go:DescribeFeeds()`

## Publishing

Once you've got your feed generator up and running and have it exposed to the internet, you can publish the feed using the script from the official BSky repo [here](https://github.com/bluesky-social/feed-generator/blob/main/scripts/publishFeedGen.ts).

Your feed will be published under _your_ DID and should show up in your profile under the `feeds` tab.

## Architecture

This repo is structured to abstract away a `Feed` interface that allows for you to add all sorts of feeds to the router.

These feeds can be simple static feeds like the `pkg/feeds/static/feed.go` implementation, or they can be much more complex feeds that draw on different data sources and filter them in cool ways to produce pages of feed items.

The `Feed` interface is defined by any struct implementing two functions:

``` go
type Feed interface {
	GetPage(ctx context.Context, feed string, userDID string, limit int64, cursor string) (feedPosts []*appbsky.FeedDefs_SkeletonFeedPost, newCursor *string, err error)
	Describe(ctx context.Context) ([]appbsky.FeedDescribeFeedGenerator_Feed, error)
}
```

`GetPage` gets a page of a feed for a given user with the limit and cursor provided, this is the main function that serves posts to a user.

`Describe` is used by the router to advertise what feeds are available, for foward compatibility, `Feed`s should be self describing in case this endpoint allows more details about feeds to be provided.

You can configure external resources and requirements in your Feed implementation before `Adding` the feed to the `FeedRouter` with `feedRouter.AddFeed([]string{"{feed_name}"}, feedInstance)`

This `Feed` interface is somewhat flexible right now but it could be better. I'm not sure if it will change in the future so keep that in mind when using this template.

- This has since been updated to allow a Feed to take in a feed name when generating a page and register multiple aliases for feeds that are supported.
