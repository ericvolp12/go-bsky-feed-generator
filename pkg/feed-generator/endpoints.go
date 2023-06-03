package feedgenerator

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	appbsky "github.com/bluesky-social/indigo/api/bsky"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (fg *FeedGenerator) GetWellKnownDID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"@context": []string{"https://www.w3.org/ns/did/v1"},
		"id":       fg.ServiceDID,
		"service": []gin.H{
			{
				"id":              "#bsky_fg",
				"type":            "BskyFeedGenerator",
				"serviceEndpoint": fg.ServiceEndpoint,
			},
		},
	})
}

func (fg *FeedGenerator) DescribeFeedGenerator(c *gin.Context) {
	tracer := otel.Tracer("feedgenerator")
	ctx, span := tracer.Start(c.Request.Context(), "FeedGenerator:DescribeFeedGenerator")
	defer span.End()

	feedDescriptions := []*appbsky.FeedDescribeFeedGenerator_Feed{}

	for _, feed := range fg.Feeds {
		feedDescription, err := feed.Describe(ctx)
		if err != nil {
			span.RecordError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		feedDescriptions = append(feedDescriptions, feedDescription)
	}

	span.SetAttributes(attribute.Int("feeds.length", len(feedDescriptions)))

	feedGeneratorDescription := appbsky.FeedDescribeFeedGenerator_Output{
		Did:   fg.FeedActorDID,
		Feeds: feedDescriptions,
	}

	c.JSON(http.StatusOK, feedGeneratorDescription)
}

func (fg *FeedGenerator) GetFeedSkeleton(c *gin.Context) {
	// Incoming requests should have a query parameter "feed" that looks like:
	// 		at://did:web:feedsky.jazco.io/app.bsky.feed.generator/feed-name
	// Also a query parameter "limit" that looks like: 50
	// Also a query parameter "cursor" that is either the empty string
	// or the cursor returned from a previous request
	tracer := otel.Tracer("feed-generator")
	ctx, span := tracer.Start(c.Request.Context(), "FeedGenerator:GetFeedSkeleton")
	defer span.End()

	feedQuery := c.Query("feed")
	if feedQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "feed query parameter is required"})
		return
	}

	c.Set("feedQuery", feedQuery)
	span.SetAttributes(attribute.String("feed.query", feedQuery))

	feedPrefix := ""
	for _, acceptablePrefix := range fg.AcceptableURIPrefixes {
		if strings.HasPrefix(feedQuery, acceptablePrefix) {
			feedPrefix = acceptablePrefix
			break
		}
	}

	if feedPrefix == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "this feed generator does not serve feeds for the given DID"})
		return
	}

	// Get the feed name from the query
	feedName := strings.TrimPrefix(feedQuery, feedPrefix)
	if feedName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "feed name is required"})
		return
	}

	span.SetAttributes(attribute.String("feed.name", feedName))
	c.Set("feedName", feedName)

	// Get the limit from the query, default to 50, maximum of 250
	limit := int64(50)
	limitQuery := c.Query("limit")
	span.SetAttributes(attribute.String("feed.limit.raw", limitQuery))
	if limitQuery != "" {
		parsedLimit, err := strconv.ParseInt(limitQuery, 10, 64)
		if err != nil {
			span.SetAttributes(attribute.Bool("feed.limit.failed_to_parse", true))
			limit = 50
		} else {
			limit = parsedLimit
			if limit > 250 {
				span.SetAttributes(attribute.Bool("feed.limit.clamped", true))
				limit = 250
			}
		}
	}

	span.SetAttributes(attribute.Int64("feed.limit.parsed", limit))

	// Get the cursor from the query
	cursor := c.Query("cursor")
	c.Set("cursor", cursor)

	if fg.Feeds == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "feed generator has no feeds configured"})
		return
	}

	feed, ok := fg.Feeds[feedName]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "feed not found"})
		return
	}

	// Get the feed items
	feedItems, newCursor, err := feed.GetPage(ctx, limit, cursor)
	if err != nil {
		span.RecordError(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get feed items: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, appbsky.FeedGetFeedSkeleton_Output{
		Feed:   feedItems,
		Cursor: newCursor,
	})
}
