package wetro

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestToolsClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/categorize/":
			response := StandardResponse{
				Success: true,
				Tokens:  5,
				Response: map[string]interface{}{
					"category": "test-category",
					"score":    0.95,
				},
			}
			json.NewEncoder(w).Encode(response)
		case "/text-generation/":
			response := StandardResponse{
				Success: true,
				Tokens:  10,
				Response: map[string]string{
					"text": "Generated text response",
				},
			}
			json.NewEncoder(w).Encode(response)
		case "/image-to-text/":
			response := StandardResponse{
				Success: true,
				Tokens:  8,
				Response: map[string]string{
					"text": "Image description",
				},
			}
			json.NewEncoder(w).Encode(response)
		case "/data-extraction/":
			response := StandardResponse{
				Success: true,
				Tokens:  12,
				Response: map[string]interface{}{
					"title":    "Test Title",
					"content":  "Test Content",
					"metadata": map[string]string{"key": "value"},
				},
			}
			json.NewEncoder(w).Encode(response)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Create a test client
	client := NewClient("test-api-key", func(c *apiClient) {
		c.baseURL = server.URL + "/"
	})
	toolsClient := client.Tools

	ctx := context.Background()

	t.Run("CategorizeData", func(t *testing.T) {
		req := CategorizeRequest{
			Type:       ResourceTypeText,
			Resource:   "test resource",
			JSONSchema: json.RawMessage(`{"type": "object"}`),
			Categories: []string{"category1", "category2"},
			Prompt:     "test prompt",
		}
		resp, err := toolsClient.CategorizeData(ctx, req)
		if err != nil {
			t.Fatalf("CategorizeData failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
		if resp.Tokens != 5 {
			t.Errorf("Expected 5 tokens, got %d", resp.Tokens)
		}
	})

	t.Run("GenerateText", func(t *testing.T) {
		req := TextGenerationRequest{
			Messages: []MessageObject{
				{
					Role:    "user",
					Content: "test message",
				},
			},
			Model: GPT35Turbo,
		}
		resp, err := toolsClient.GenerateText(ctx, req)
		if err != nil {
			t.Fatalf("GenerateText failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
		if resp.Tokens != 10 {
			t.Errorf("Expected 10 tokens, got %d", resp.Tokens)
		}
	})

	t.Run("ImageToText", func(t *testing.T) {
		req := ImageToTextRequest{
			ImageURL: "https://example.com/image.jpg",
			Query:    "What is in this image?",
		}
		resp, err := toolsClient.ImageToText(ctx, req)
		if err != nil {
			t.Fatalf("ImageToText failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
		if resp.Tokens != 8 {
			t.Errorf("Expected 8 tokens, got %d", resp.Tokens)
		}
	})

	t.Run("ExtractData", func(t *testing.T) {
		req := DataExtractionRequest{
			WebURL: "https://example.com",
			Schema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title":    map[string]interface{}{"type": "string"},
					"content":  map[string]interface{}{"type": "string"},
					"metadata": map[string]interface{}{"type": "object"},
				},
			},
		}
		resp, err := toolsClient.ExtractData(ctx, req)
		if err != nil {
			t.Fatalf("ExtractData failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
		if resp.Tokens != 12 {
			t.Errorf("Expected 12 tokens, got %d", resp.Tokens)
		}
	})
}
