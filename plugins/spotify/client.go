package spotify

import (
	"context"
	"fmt"

	spotifylib "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	api *spotifylib.Client
	ctx context.Context
}

func NewClient(clientID, clientSecret string) (*Client, error) {
	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	httpClient := config.Client(ctx)
	api := spotifylib.New(httpClient)

	return &Client{api: api, ctx: ctx}, nil
}

func (c *Client) API() *spotifylib.Client {
	return c.api
}

func (c *Client) Context() context.Context {
	return c.ctx
}

func (c *Client) SearchArtists(name string, limit int) (*spotifylib.SearchResult, error) {
	result, err := c.api.Search(c.ctx, name, spotifylib.SearchTypeArtist, spotifylib.Limit(limit))
	if err != nil {
		return nil, fmt.Errorf("searching for artist %q: %w", name, err)
	}
	return result, nil
}
