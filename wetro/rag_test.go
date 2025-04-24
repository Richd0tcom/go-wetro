package wetro

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRAGClient(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/collection/create/":
			response := CollectionCreateResponse{
				Success: true,
			}
			json.NewEncoder(w).Encode(response)
		case "/collection/get/test-collection/":
			response := GetCollectionResponse{
				Success: true,
				CollectionID: "test-collection",
			}
			json.NewEncoder(w).Encode(response)
		case "/collection/all/":
			response := ListCollectionResponse{
				Count: 2,
				Results: []CollectionItem{
					{CollectionID: "collection1"},
					{CollectionID: "collection2"},
				},
			}
			json.NewEncoder(w).Encode(response)
		case "/collection/query/":
			response := StandardResponse{
				Success: true,
				Tokens:  10,
				Response: map[string]interface{}{
					"results": []map[string]interface{}{
						{"text": "result1", "score": 0.9},
						{"text": "result2", "score": 0.8},
					},
				},
			}
			json.NewEncoder(w).Encode(response)
		case "/collection/chat/":
			response := StandardResponse{
				Success: true,
				Tokens:  15,
				Response: map[string]string{
					"response": "This is a test response",
				},
			}
			json.NewEncoder(w).Encode(response)
		case "/resource/insert/":
			response := ResourceInsertResponse{
				Success: true,
				Tokens:  6,
			}
			json.NewEncoder(w).Encode(response)
		case "/resource/remove/":
			response := ResourceDeleteResponse{
				Success: true,
			}
			json.NewEncoder(w).Encode(response)
		case "/collection/delete/":
			response := DeleteCollectionResponse{
				Success: true,
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
	ragClient := client.RAG

	ctx := context.Background()

	t.Run("CreateCollection", func(t *testing.T) {
		resp, err := ragClient.CreateCollection(ctx, "test-collection")
		if err != nil {
			t.Fatalf("CreateCollection failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("GetCollection", func(t *testing.T) {
		resp, err := ragClient.GetCollection(ctx, "test-collection")
		if err != nil {
			t.Fatalf("GetCollection failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
		if resp.CollectionID != "test-collection" {
			t.Errorf("Expected collection ID 'test-collection', got '%s'", resp.CollectionID)
		}
	})

	t.Run("ListCollections", func(t *testing.T) {
		resp, err := ragClient.ListCollections(ctx)
		if err != nil {
			t.Fatalf("ListCollections failed: %v", err)
		}
		if resp.Count < 1{
			t.Error("Expected success to be true")
		}
		if len(resp.Results) != 2 {
			t.Errorf("Expected 2 collections, got %d", len(resp.Results))
		}
	})

	t.Run("QueryCollection", func(t *testing.T) {
		req := QueryRequest{
			CollectionID: "test-collection",
			Query:        "test query",
		}
		resp, err := ragClient.QueryCollection(ctx, req)
		if err != nil {
			t.Fatalf("QueryCollection failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
		if resp.Tokens != 10 {
			t.Errorf("Expected 10 tokens, got %d", resp.Tokens)
		}
	})

	t.Run("ChatWithCollection", func(t *testing.T) {
		req := ChatRequest{
			CollectionID: "test-collection",
			Message:      "test message",
		}
		resp, err := ragClient.ChatWithCollection(ctx, req)
		if err != nil {
			t.Fatalf("ChatWithCollection failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
		if resp.Tokens != 15 {
			t.Errorf("Expected 15 tokens, got %d", resp.Tokens)
		}
	})

	t.Run("InsertResource", func(t *testing.T) {
		resp, err := ragClient.InsertResource(ctx, "test-collection", "test-resource", ResourceTypeText)
		if err != nil {
			t.Fatalf("InsertResource failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
		if resp.Tokens != 6 {
			t.Errorf("Expected 6 tokens, got %d", resp.Tokens)
		}
	})

	t.Run("RemoveResource", func(t *testing.T) {
		req := ResourceDeleteRequest{
			CollectionID: "test-collection",
			ResourceID:   "test-resource",
		}
		resp, err := ragClient.RemoveResource(ctx, req)
		if err != nil {
			t.Fatalf("RemoveResource failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
	})

	t.Run("DeleteCollection", func(t *testing.T) {
		resp, err := ragClient.DeleteCollection(ctx, "test-collection")
		if err != nil {
			t.Fatalf("DeleteCollection failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success to be true")
		}
	})
}
