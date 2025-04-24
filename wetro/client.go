// Copyright 2025 Richd0tcom. All rights reserved.
// Use of this source code is governed by an MIT style
// license that can be found in the LICENSE file.

/*
Package wetrocloud provides a Go client for interacting with the WetroCloud API.
It supports both RAG (Retrieval-Augmented Generation) and various AI tools functionality.
*/
package wetrocloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// APIClient represents the main client for interacting with the WetroCloud API.
// It handles authentication, request formatting, and response processing.
type APIClient struct {
	baseURL    string
	apiKey     string
	apiVersion string
	httpClient *http.Client
}

// Client represents the main entry point for the WetroCloud SDK.
// It provides access to both RAG and Tools functionality.
type Client struct {
	RAG   *RAGClient
	Tools *ToolsClient
}

// ClientOption represents a function that can modify the APIClient configuration.
// It's used for setting various client options during initialization.
type ClientOption func(*APIClient)

func NewClient(apiKey string, options ...ClientOption) *Client {
	apiClient := &APIClient{
		baseURL:    "https://api.wetrocloud.com/",
		apiKey:     apiKey,
		apiVersion: "v1",
		httpClient: &http.Client{},
	}

	for _, opt := range options {
		opt(apiClient)
	}

	return &Client{
		RAG:   NewRAGClient(apiClient),
		Tools: NewToolsClient(apiClient),
	}
}

func NewRAGClient(api *APIClient) *RAGClient {
	return &RAGClient{client: api}
}

func NewToolsClient(api *APIClient) *ToolsClient {
	return &ToolsClient{client: api}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *APIClient) {
		c.httpClient = client
	}
}

// WithAPIVersion sets a custom API version
func WithAPIVersion(version string) ClientOption {
	return func(c *APIClient) {
		c.apiVersion = version
	}
}

func (c *APIClient) doRequest(ctx context.Context, method, endpoint string, params map[string]string, data interface{}, response interface{}) error {
	url := fmt.Sprintf("%s%s%s", c.baseURL, c.apiVersion, endpoint)

	// Add referrer parameter
	if params == nil {
		params = make(map[string]string)
	}
	params["referrer"] = "GO_SDK"

	// Create request
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	// Add query parameters
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// bodyBytes, err := io.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)

	// Check for errors
	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error   string `json:"error"`
			Detail  string `json:"detail"`
			Payload any    `json:"payload,omitempty"`
		}
		fmt.Println("kjhflkjdhoj", resp.StatusCode)
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return &APIError{
				Message:    "Failed to parse error response",
				StatusCode: resp.StatusCode,
			}
		}

		return &APIError{
			Message:    parseError(resp),
			StatusCode: resp.StatusCode,
			Payload:    errorResp.Payload,
		}
	}

	// Parse response
	if response != nil {
		return json.NewDecoder(resp.Body).Decode(response)
	}
	return nil
}

func (c *APIClient) doMultipartRequest(ctx context.Context, method, endpoint string, data map[string]interface{}, response interface{}) error {
	url := fmt.Sprintf("%s%s%s", c.baseURL, c.apiVersion, endpoint)

	// Create multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add form fields
	for k, v := range data {
		if str, ok := v.(string); ok {
			writer.WriteField(k, str)
		} else {
			jsonData, err := json.Marshal(v)
			if err != nil {
				return err
			}
			writer.WriteField(k, string(jsonData))
		}
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, &buf)
	if err != nil {
		return err
	}

	// Add headers
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.apiKey))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error   string `json:"error"`
			Detail  string `json:"detail"`
			Payload any    `json:"payload"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return &APIError{
				Message:    "Failed to parse error response",
				StatusCode: resp.StatusCode,
			}
		}

		return &APIError{
			Message:    parseError(resp),
			StatusCode: resp.StatusCode,
			Payload:    errorResp.Payload,
		}
	}

	// Parse response
	if response != nil {
		return json.NewDecoder(resp.Body).Decode(response)
	}
	return nil
}

func (c *APIClient) uploadBytes(ctx context.Context, collectionID string, resource any) (string, error) {

	if _, isReadable := resource.(io.Reader); !isReadable {
		return "", fmt.Errorf("Invalid Resource")
	}
	var reader io.Reader = resource.(io.Reader)

	id, err := GenerateID()
	if err != nil {
		return "", err
	}
	return c.upload(ctx, reader, collectionID, id)
}

// Helper method for file upload
func (c *APIClient) uploadFile(ctx context.Context, collectionID string, filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("File %s does not exist", filePath)
		}
		return "", err
	}
	defer file.Close()

	return c.upload(ctx, file, collectionID, filepath.Base(filePath))
}

func (c *APIClient) upload(ctx context.Context, reader io.Reader, collectionID, filename string) (string, error) {
	uploadURL := "https://file-upload-service-python.vercel.app/upload/"

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add collection_id
	writer.WriteField("collection_id", collectionID)

	// Add file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, reader)
	if err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", errors.New("file upload failed")
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	url, ok := result["url"]
	if !ok {
		return "", errors.New("no URL in response")
	}

	return url, nil
}
