package wetrocloud

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RAGClient struct {
	client *APIClient
}


func (c *RAGClient) CreateCollection(ctx context.Context, id string) (*CollectionCreateResponse, error) {

	requestData := map[string]string{
		"collection_id": id,
	}
	response := &CollectionCreateResponse{}

	err := c.client.doRequest(ctx, http.MethodPost, "/collection/create/", nil, requestData, response)
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	return response, nil
}

// GetCollection retrieves a collection
func (c *RAGClient) GetCollection(ctx context.Context, collectionID string) (*GetCollectionResponse, error) {
	var response GetCollectionResponse
	err := c.client.doRequest(ctx, "GET", fmt.Sprintf("/collection/get/%s/", collectionID), nil, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// ListCollections lists all collections
func (c *RAGClient) ListCollections(ctx context.Context) (*ListCollectionResponse, error) {
	var response ListCollectionResponse
	err := c.client.doRequest(ctx, http.MethodGet, "/collection/all/", nil, nil, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// QueryCollection queries a collection
func (c *RAGClient) QueryCollection(ctx context.Context, request QueryRequest) (*StandardResponse, error) {
	var response StandardResponse

	//TODO: validate query request

	err := c.client.doRequest(ctx, "POST", "/collection/query/", nil, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// ChatWithCollection chats with a collection
func (c *RAGClient) ChatWithCollection(ctx context.Context, request ChatRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.client.doRequest(ctx, "POST", "/collection/chat/", nil, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// InsertResource inserts a resource into a collection
func (c *RAGClient) InsertResource(ctx context.Context, collectionID string, resource any, resourceType ResourceType) (*ResourceInsertResponse, error) {
	// Handle file upload if resource is a file path
	var resourceURL string
	var response ResourceInsertResponse
	if path, ok := resource.(string); ok && !strings.HasPrefix(path, "http") {
		url, err := c.client.uploadFile(ctx, collectionID, path)
		if err != nil {
			return nil, err
		}
		resourceURL = url
	} else if _, ok := resource.(io.Reader); ok {
		fmt.Println("oass vibe check")
		url, err := c.client.uploadBytes(ctx, collectionID, resource)
		if err != nil {
			return nil, err
		}
		resourceURL = url
	} else {
		resourceURL = fmt.Sprintf("%v", resource)
	}
	payload := ResourceInsertRequest{
		CollectionID: collectionID,
		Resource:     resourceURL,
		Type:         resourceType,
	}

	err := c.client.doRequest(ctx, "POST", "/resource/insert/", nil, payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// RemoveResource removes a resource from a collection
func (c *RAGClient) RemoveResource(ctx context.Context, request ResourceDeleteRequest) (*ResourceDeleteResponse, error) {
	var response ResourceDeleteResponse
	err := c.client.doRequest(ctx, "DELETE", "/resource/remove/", nil, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// DeleteCollection deletes a collection
func (c *RAGClient) DeleteCollection(ctx context.Context, collectionID string) (*DeleteCollectionResponse, error) {
	var response DeleteCollectionResponse
	err := c.client.doRequest(ctx, "DELETE", "/collection/delete/", nil, map[string]any{
		"collection_id": collectionID,
	}, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}