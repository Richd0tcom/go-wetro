package wetrocloud

import (
	"context"
)

// Resource types supported by the API
const (
	ResourceTypeFile    = "file"
	ResourceTypeText    = "text"
	ResourceTypeJSON    = "json"
	ResourceTypeWeb     = "web"
	ResourceTypeYoutube = "youtube"
	ResourceTypeAudio   = "audio"
)

// TextResource creates a resource insert request for text content
func (c *Client) TextResource(ctx context.Context, collectionID, text string, metadata interface{}) (*ResourceInsertResponse, error) {
	request := &ResourceInsertRequest{
		CollectionID: collectionID,
		Type:         ResourceTypeText,
		Resource:     text,
	}

	return c.insertResource(ctx, request)
}

// FileResource creates a resource insert request for a file URL
func (c *Client) FileResource(ctx context.Context, collectionID, fileURL string, metadata interface{}) (*ResourceInsertResponse, error) {
	request := &ResourceInsertRequest{
		CollectionID: collectionID,
		Type:         ResourceTypeFile,
		Resource:     fileURL,
	}

	return c.insertResource(ctx, request)
}

// JSONResource creates a resource insert request for JSON data
func (c *Client) JSONResource(ctx context.Context, collectionID, jsonData string, metadata interface{}) (*ResourceInsertResponse, error) {
	request := &ResourceInsertRequest{
		CollectionID: collectionID,
		Type:         ResourceTypeJSON,
		Resource:     jsonData,
	}

	return c.insertResource(ctx, request)
}

// WebResource creates a resource insert request for a web URL
func (c *Client) WebResource(ctx context.Context, collectionID, webURL string, metadata interface{}) (*ResourceInsertResponse, error) {
	request := &ResourceInsertRequest{
		CollectionID: collectionID,
		Type:         ResourceTypeWeb,
		Resource:     webURL,
	}

	return c.insertResource(ctx, request)
}

// YoutubeResource creates a resource insert request for a YouTube URL
func (c *Client) YoutubeResource(ctx context.Context, collectionID, youtubeURL string, metadata interface{}) (*ResourceInsertResponse, error) {
	request := &ResourceInsertRequest{
		CollectionID: collectionID,
		Type:         ResourceTypeYoutube,
		Resource:     youtubeURL,
	}

	return c.insertResource(ctx, request)
}

// AudioResource creates a resource insert request for an audio file URL
func (c *Client) AudioResource(ctx context.Context, collectionID, audioURL string, metadata interface{}) (*ResourceInsertResponse, error) {
	request := &ResourceInsertRequest{
		CollectionID: collectionID,
		Type:         ResourceTypeAudio,
		Resource:     audioURL,
	}

	return c.insertResource(ctx, request)
}
