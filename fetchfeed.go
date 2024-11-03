package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"

	"github.com/chichigami/blog-aggregator/internal/rss"
)

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	req, reqErr := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	req.Header.Add("User-Agent", "gator")
	client := &http.Client{}
	res, clientErr := client.Do(req)
	if clientErr != nil {
		return nil, clientErr
	}
	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	var feed rss.RSSFeed
	unmarshErr := xml.Unmarshal(body, &feed)
	if unmarshErr != nil {
		return nil, unmarshErr
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	return &feed, nil
}
