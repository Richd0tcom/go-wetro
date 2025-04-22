package wetrocloud

import "context"

type ToolsClient struct {
	client *APIClient
}


// CategorizeData categorizes data
func (c *ToolsClient) CategorizeData(ctx context.Context, payload CategorizeRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.client.doRequest(ctx, "POST", "/categorize/", nil, payload, &response)
	return &response, err
}

// GenerateText generates text
func (c *ToolsClient) GenerateText(ctx context.Context, payload TextGenerationRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.client.doRequest(ctx, "POST", "/text-generation/", nil, payload, &response)
	return &response, err
}

// ImageToText generates text from an image
func (c *ToolsClient) ImageToText(ctx context.Context, payload ImageToTextRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.client.doRequest(ctx, "POST", "/image-to-text/", nil, payload, &response)
	return &response, err
}

// ExtractData extracts data from a website
func (c *ToolsClient) ExtractData(ctx context.Context, payload DataExtractionRequest) (*StandardResponse, error) {
	var response StandardResponse

	err := c.client.doRequest(ctx, "POST", "/data-extraction/", nil, payload, &response)
	return &response, err
}