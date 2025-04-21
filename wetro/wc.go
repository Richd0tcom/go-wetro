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
	"strings"
)

type RAGClient struct {
	client *APIClient
}
type ToolsClient struct {
	client *APIClient
}

type APIClient struct {
	baseURL    string
	apiKey     string
	apiVersion string
	httpClient *http.Client
}

type Client struct {
	RAG   *RAGClient
	Tools *ToolsClient
}

// ClientOption is a function that configures a Client
type ClientOption func(*APIClient)

func NewClient(apiKey string, options ...ClientOption) *Client {
	apiClient := &APIClient{
		baseURL:    "https://api.wetrocloud.com",
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

	// Check for errors
	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error   string      `json:"error"`
			Detail  string      `json:"detail"`
			Payload any `json:"payload,omitempty"`
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
			Error   string      `json:"error"`
			Detail  string      `json:"detail"`
			Payload interface{} `json:"payload"`
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

	id, err := generateUUID()
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

func (c *RAGClient) CreateCollection(ctx context.Context, id string) (*CollectionCreateResponse, error) {

	requestData := map[string]string{
		"collection_id": id,
	}
	response := &CollectionCreateResponse{}

	err := c.client.doRequest(ctx, http.MethodPost, "collection/create", map[string]string{}, requestData, response)
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	return response, nil
}

// GetCollection retrieves a collection
func (c *APIClient) GetCollection(ctx context.Context, collectionID string) (*GetCollectionResponse, error) {
	var response GetCollectionResponse
	err := c.doRequest(ctx, "GET", fmt.Sprintf("/v1/get/%s/", collectionID), nil, nil, &response)
	return &response, err
}

// ListCollections lists all collections
func (c *APIClient) ListCollections(ctx context.Context) (*ListCollectionResponse, error) {
	var response ListCollectionResponse
	err := c.doRequest(ctx, http.MethodGet, "/v1/list/", nil, nil, &response)
	return &response, err
}

// QueryCollection queries a collection
func (c *APIClient) QueryCollection(ctx context.Context, request QueryRequest) (*StandardResponse, error) {
	var response StandardResponse

	//validate query request

	err := c.doRequest(ctx, "POST", fmt.Sprintf("/v1/query/%s/", request.CollectionID), nil, request, &response)
	return &response, err
}

// ChatWithCollection chats with a collection
func (c *APIClient) ChatWithCollection(ctx context.Context, request ChatRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.doRequest(ctx, "POST", fmt.Sprintf("/v1/chat/%s/", request.CollectionID), nil, request, &response)
	return &response, err
}

// InsertResource inserts a resource into a collection
func (c *APIClient) InsertResource(ctx context.Context, collectionID string, resource any, resourceType ResourceType) (*ResourceInsertResponse, error) {
	// Handle file upload if resource is a file path
	var resourceURL string
	var response ResourceInsertResponse
	if path, ok := resource.(string); ok && !strings.HasPrefix(path, "http") {
		url, err := c.uploadFile(ctx, collectionID, path)
		if err != nil {
			return nil, err
		}
		resourceURL = url
	} else if _, ok := resource.(io.Reader); ok {
		url, err := c.uploadBytes(ctx, collectionID, resource)
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

	err := c.doRequest(ctx, "POST", "/v1/insert/", nil, payload, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// RemoveResource removes a resource from a collection
func (c *APIClient) RemoveResource(ctx context.Context, request ResourceDeleteRequest) (*ResourceDeleteResponse, error) {
	var response ResourceDeleteResponse
	err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/v1/remove/"), nil, request, &response)
	return &response, err
}

// DeleteCollection deletes a collection
func (c *APIClient) DeleteCollection(ctx context.Context, collectionID string) (*DeleteCollectionResponse, error) {
	var response DeleteCollectionResponse
	err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/v1/delete/%s/", collectionID), nil, map[string]any{
		"collection_id": collectionID,
	}, &response)
	return &response, err
}

// CategorizeData categorizes data
func (c *APIClient) CategorizeData(ctx context.Context, payload CategorizeRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.doRequest(ctx, "POST", "/v1/categorize/", nil, payload, &response)
	return &response, err
}

// GenerateText generates text
func (c *APIClient) GenerateText(ctx context.Context, payload TextGenerationRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.doRequest(ctx, "POST", "/v1/generate/", nil, payload, &response)
	return &response, err
}

// ImageToText generates text from an image
func (c *APIClient) ImageToText(ctx context.Context, payload ImageToTextRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.doRequest(ctx, "POST", "/v1/image-to-text/", nil, payload, &response)
	return &response, err
}

// ExtractData extracts data from a website
func (c *APIClient) ExtractData(ctx context.Context, payload DataExtractionRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.doRequest(ctx, "POST", "/v1/extract/", nil, payload, &response)
	return &response, err
}
