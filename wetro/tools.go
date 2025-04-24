// Copyright 2025 Richd0tcom. All rights reserved.
// Use of this source code is governed by an MIT style
// license that can be found in the LICENSE file.


package wetro

import (
	"context"
	"net/http"
)

// ToolsClient provides methods for working with various AI tools.
// It handles text generation, categorization, and other AI operations.
type ToolsClient struct {
	client *APIClient
}

// CategorizeData categorizes data
func (c *ToolsClient) CategorizeData(ctx context.Context, payload CategorizeRequest) (StandardResponse, error) {
	var response StandardResponse

	err := c.client.doRequest(ctx, http.MethodPost, "/categorize/", nil, payload, &response)
	if err != nil {
		return StandardResponse{}, err
	}
	return response, nil
}

// GenerateText generates text
func (c *ToolsClient) GenerateText(ctx context.Context, payload TextGenerationRequest) (StandardResponse, error) {
	var response StandardResponse

	err := c.client.doRequest(ctx, http.MethodPost, "/text-generation/", nil, payload, &response)

	if err != nil {
		return StandardResponse{}, err
	}
	return response, nil
}

// ImageToText generates text from an image
func (c *ToolsClient) ImageToText(ctx context.Context, payload ImageToTextRequest) (StandardResponse, error) {
	var response StandardResponse

	err := c.client.doRequest(ctx, http.MethodPost, "/image-to-text/", nil, payload, &response)

	if err != nil {
		return StandardResponse{}, err
	}
	return response, nil
}

// ExtractData extracts data from a website
func (c *ToolsClient) ExtractData(ctx context.Context, payload DataExtractionRequest) (StandardResponse, error) {
	var response StandardResponse

	err := c.client.doRequest(ctx, http.MethodPost, "/data-extraction/", nil, payload, &response)

	if err != nil {
		return StandardResponse{}, err
	}
	return response, nil
}
