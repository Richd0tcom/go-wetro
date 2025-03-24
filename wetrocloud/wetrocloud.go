package wetrocloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
)

// Client represents a WetroCloud API client
type Client struct {
	BaseURL    string
	APIKey     string
	APIVersion string
	HTTPClient *http.Client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = client
	}
}

// WithAPIVersion sets a custom API version
func WithAPIVersion(version string) ClientOption {

	return func(c *Client) {
		c.APIVersion = version
	}
}

// NewClient creates a new WetroCloud API client
func NewClient(baseURL, apiKey string, options ...ClientOption) *Client {
	client := &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		APIVersion: "v1",
		HTTPClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// buildURL builds a complete URL for the API
func (c *Client) buildURL(endpoint string) string {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return ""
	}

	u.Path = path.Join(c.APIVersion, endpoint)
	return u.String()
}

// doRequest performs an HTTP request and decodes the response
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body any, target any) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(endpoint), bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}

	return nil
}

func (c *Client) doMultipartRequest(ctx context.Context, method, endpoint, contentType string, buf io.Reader, target any) error {

	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(endpoint), buf)

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.APIKey))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if target != nil {
		return json.NewDecoder(resp.Body).Decode(target)
	}

	return nil
}

// CreateCollection creates a new collection
func (c *Client) CreateCollection(ctx context.Context, id string) (*CollectionCreateResponse, error) {

	buf := &bytes.Buffer{}

	mpw := multipart.NewWriter(buf)

	idWriter, err := mpw.CreateFormField("collection_id")

	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}
	_, err = idWriter.Write([]byte(id))

	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	response := &CollectionCreateResponse{}
	contentType := mpw.FormDataContentType()

	err = mpw.Close()
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	err = c.doMultipartRequest(ctx, http.MethodPost, "collection/create", contentType, buf, response)
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	return response, nil
}

// GetAllCollections retrieves all collections
func (c *Client) GetAllCollections(ctx context.Context) (*CollectionsResponse, error) {
	response := &CollectionsResponse{}

	err := c.doRequest(ctx, http.MethodGet, "collection/all", nil, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// InsertResource inserts a resource into a collection
func (c *Client) insertResource(ctx context.Context, request *ResourceInsertRequest) (*ResourceInsertResponse, error) {
	response := &ResourceInsertResponse{}

	err := c.doRequest(ctx, http.MethodPost, "resource/insert", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// QueryCollection queries a collection
func (c *Client) QueryCollection(ctx context.Context, request *QueryRequest) (*StandardResponse, error) {
	response := &StandardResponse{}

	err := c.doRequest(ctx, http.MethodPost, "collection/query", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ChatCollection(ctx context.Context, request *ChatRequest) (*StandardResponse, error) {
	response := &StandardResponse{}

	buf := &bytes.Buffer{}

	mpw := multipart.NewWriter(buf)

	idWriter, _ := mpw.CreateFormField("collection_id")
	_, _ = idWriter.Write([]byte(request.CollectionID))

	msgWriter, _ := mpw.CreateFormField("message")
	_, _ = msgWriter.Write([]byte(request.Message))

	chWriter, _ := mpw.CreateFormField("chat_history")
	_, _ = chWriter.Write([]byte(request.ChatHistory))

	contentType := mpw.FormDataContentType()

	err := mpw.Close()
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	err = c.doMultipartRequest(ctx, http.MethodPost, "collection/query", contentType, buf, response)
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	return response, nil
}
func (c *Client) RemoveResource(ctx context.Context, request *ResourceDeleteRequest) (*StandardResponse, error) {
	response := &StandardResponse{}

	buf := &bytes.Buffer{}

	mpw := multipart.NewWriter(buf)

	idWriter, _ := mpw.CreateFormField("collection_id")
	_, err := idWriter.Write([]byte(request.CollectionID))
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	resIDWriter, _ := mpw.CreateFormField("resource_id")
	_, err = resIDWriter.Write([]byte(request.ResourceID))
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	contentType := mpw.FormDataContentType()

	err = mpw.Close()
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	err = c.doMultipartRequest(ctx, http.MethodPost, "/resource/remove/", contentType, buf, response)
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}
	return response, nil
}

// DeleteCollection deletes a collection
func (c *Client) DeleteCollection(ctx context.Context, collectionID string) (*StandardResponse, error) {
	payload := map[string]string{"collection_id": collectionID}
	response := &StandardResponse{}

	err := c.doRequest(ctx, http.MethodDelete, "collection/delete", payload, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Categorize categorizes data according to specific rules or schema
func (c *Client) Categorize(ctx context.Context, request *CategorizeRequest) (*StandardResponse, error) {
	response := &StandardResponse{}

	err := c.doRequest(ctx, http.MethodPost, "/categorize", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) TextGeneration(ctx context.Context, request *TextGenerationRequest) (*StandardResponse, error) {
	response := &StandardResponse{}

	buf := &bytes.Buffer{}

	mpw := multipart.NewWriter(buf)

	msgWriter, _ := mpw.CreateFormField("messages")
	msg, err := ToJSONSchema(request.Messages)
	if err != nil {
		return nil, err
	}

	_, err = msgWriter.Write([]byte(msg))
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	mdlWriter, _ := mpw.CreateFormField("model")
	_, err = mdlWriter.Write([]byte(request.Model))
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	contentType := mpw.FormDataContentType()

	err = mpw.Close()
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	err = c.doMultipartRequest(ctx, http.MethodPost, "/text-generation/", contentType, buf, response)
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	return response, nil
}

func (c *Client) Image2FreeText(ctx context.Context, request *ImageToFreeTextRequest) (*StandardResponse, error) {
	response := &StandardResponse{}

	err := c.doRequest(ctx, http.MethodPost, "/image-to-text/", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) DataExtraction(ctx context.Context, request *DataExtractionRequest) (*StandardResponse, error) {
	response := &StandardResponse{}

	buf := &bytes.Buffer{}

	mpw := multipart.NewWriter(buf)

	webWriter, _:=mpw.CreateFormField("website")
	_, err:=webWriter.Write([]byte(request.WebURL))
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	jsWriter, _:= mpw.CreateFormField("json_schema")

	js, _:= ToJSONSchema(request.Schema)

	_, err=jsWriter.Write([]byte(js))
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	contentType:= mpw.FormDataContentType()
	mpw.Close()

	err= c.doMultipartRequest(ctx, http.MethodPost, "/data-extraction/", contentType, buf, response)
	if err != nil {
		//TODO(Rich): replace error with well formatted error
		return nil, err
	}

	return response, nil
}