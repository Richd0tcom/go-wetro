package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Richd0tcom/wetrocloud-sdk-go/wetrocloud"
)

func main() {
	// Get API credentials from environment
	apiKey := os.Getenv("WETROCLOUD_API_KEY")
	if apiKey == "" {
		log.Fatal("WETROCLOUD_API_KEY environment variable is required")
	}

	// Initialize the client with a custom timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	client := wetrocloud.NewClient(
		"https://api.wetrocloud.com",
		apiKey,
		wetrocloud.WithHTTPClient(httpClient),
	)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create a new knowledge base collection
	fmt.Println("Creating knowledge base collection...")
	collection, err := client.CreateCollection(ctx, "Product Documentation KB")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	collectionID := collection.CollectionID
	fmt.Printf("Created collection with ID: %s\n", collectionID)

	// Add different types of resources to the collection
	fmt.Println("\nAdding resources to the collection...")

	// Add a text resource
	textData := `
	WetroCloud is a powerful platform for building RAG applications.
	It provides APIs for document processing, retrieval, and generation.
	Developers can use it to create AI-powered applications with ease.
	`
	textMetadata := map[string]interface{}{
		"title":    "WetroCloud Overview",
		"category": "product",
		"tags":     []string{"overview", "introduction"},
	}

	_, err = client.TextResource(ctx, collectionID, textData, textMetadata)
	if err != nil {
		log.Fatalf("Failed to add text resource: %v", err)
	}
	fmt.Println("✓ Added text resource")

	// Add a web resource
	webURL := "https://example.com/wetrocloud-docs"
	webMetadata := map[string]interface{}{
		"title":    "Official Documentation",
		"category": "documentation",
		"tags":     []string{"docs", "reference"},
	}

	_, err = client.WebResource(ctx, collectionID, webURL, webMetadata)
	if err != nil {
		log.Fatalf("Failed to add web resource: %v", err)
	}
	fmt.Println("✓ Added web resource")

	// Add a JSON resource
	jsonData := map[string]interface{}{
		"features": []string{
			"Document processing",
			"Semantic search",
			"LLM integration",
			"Custom agents",
		},
		"pricing": map[string]interface{}{
			"basic":      "Free",
			"pro":        "$49/month",
			"enterprise": "Contact sales",
		},
	}
	jsonMetadata := map[string]interface{}{
		"title":    "Product Features and Pricing",
		"category": "marketing",
		"tags":     []string{"features", "pricing"},
	}

	jd, _ := wetrocloud.ToJSONSchema(jsonData)
	_, err = client.JSONResource(ctx, collectionID, jd, jsonMetadata)
	if err != nil {
		log.Fatalf("Failed to add JSON resource: %v", err)
	}
	fmt.Println("✓ Added JSON resource")

	// Query the knowledge base
	fmt.Println("\nQuerying the knowledge base...")
	queryReq := &wetrocloud.QueryRequest{
		CollectionID: collectionID,
		Query:        "What features does WetroCloud offer?",
	}

	queryResp, err := client.QueryCollection(ctx, queryReq)
	if err != nil {
		log.Fatalf("Failed to query collection: %v", err)
	}

	fmt.Printf("\nQuery response:\n%s\n", queryResp.Response)
	fmt.Printf("Tokens used: %d\n", queryResp.Tokens)

	// Categorize some data
	fmt.Println("\nCategorizing data...")
	categorizeData := "WetroCloud offers document processing capabilities with semantic search."
	categorizeSchema := map[string]interface{}{
		"categories": []string{"Product", "Technology", "Documentation", "Marketing"},
	}
	schemaStr, _:= wetrocloud.ToJSONSchema(categorizeSchema)
	categorizeReq := &wetrocloud.CategorizeRequest{
		Data:   categorizeData,
		JSONSchema: schemaStr,
	}

	categorizeResp, err := client.Categorize(ctx, categorizeReq)
	if err != nil {
		log.Fatalf("Failed to categorize data: %v", err)
	}

	fmt.Printf("\nCategorization response:\n%s\n", categorizeResp.Response)
	fmt.Printf("Tokens used: %d\n", categorizeResp.Tokens)

	// Cleanup (optional - commenting out to preserve the collection)
	/*
		fmt.Println("\nCleaning up (deleting collection)...")
		_, err = client.DeleteCollection(ctx, collectionID)
		if err != nil {
			log.Fatalf("Failed to delete collection: %v", err)
		}
		fmt.Println("✓ Collection deleted")
	*/

	fmt.Println("\nExample completed successfully!")
}
