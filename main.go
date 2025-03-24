package main


import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Richd0tcom/wetrocloud-sdk-go/wetrocloud"
)

func main() {
	// Replace with your actual API key and base URL
	apiKey := os.Getenv("WETROCLOUD_API_KEY")
	if apiKey == "" {
		log.Fatal("WETROCLOUD_API_KEY environment variable is required")
	}

	baseURL := os.Getenv("WETROCLOUD_API_URL")
	if baseURL == "" {
		baseURL = "https://api.wetrocloud.com" // Default API URL
	}

	// Create a new client
	client := wetrocloud.NewClient(baseURL, apiKey)

	// Use context with timeout or cancellation as needed for your application
	ctx := context.Background()

	// Example: Create a new collection
	collectionResp, err := client.CreateCollection(ctx, "My Test Collection")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	fmt.Printf("Created collection with ID: %s\n", collectionResp.CollectionID)

	// Example: Add a text resource to the collection
	textResp, err := client.TextResource(ctx, collectionResp.CollectionID, "This is a sample text for RAG processing", nil)
	if err != nil {
		log.Fatalf("Failed to add text resource: %v", err)
	}

	fmt.Printf("Added text resource. Success: %v, Tokens used: %d\n", textResp.Success, textResp.Tokens)

	// Example: Query the collection
	queryReq := &wetrocloud.QueryRequest{
		CollectionID: collectionResp.CollectionID,
		Query:        "What is this text about?",
	}

	queryResp, err := client.QueryCollection(ctx, queryReq)
	if err != nil {
		log.Fatalf("Failed to query collection: %v", err)
	}

	fmt.Printf("Query response: %s\n", queryResp.Response)
}
