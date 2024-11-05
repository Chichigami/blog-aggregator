package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/chichigami/blog-aggregator/internal/database"
	"github.com/chichigami/blog-aggregator/internal/rss"
	"github.com/google/uuid"
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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		fmt.Println("need a title and url")
		os.Exit(1)
	}
	name := cmd.args[0]
	url := cmd.args[1]
	user_id, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user_id.ID,
	}
	_, fetchErr := s.db.CreateFeed(context.Background(), newFeed)
	if fetchErr != nil {
		return err
	}
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, v := range feeds {
		fmt.Printf("Name of Feed: %s\nURL: %s\nSubmitted by: %s\n", v.FeedName, v.Url, v.UserName)
	}
	return nil
}
