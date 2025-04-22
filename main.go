package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	wetro "github.com/Richd0tcom/go-wetro/wetro"
)

func main() {
	// Initialize the client with your API key
	client := wetro.NewClient("WETRO_API_KEY")

	// Create a context
	ctx := context.Background()

	// RAG Client Examples
	ragExamples(client, ctx)

	// Tools Client Examples
	toolsExamples(client, ctx)
}

func ragExamples(client *wetro.Client, ctx context.Context) {
	fmt.Println("\n=== RAG Client Examples ===")

	// Example 1: Create a collection
	createResp, err := client.RAG.CreateCollection(ctx, "my-docs")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}
	fmt.Printf("Created collection: %+v\n", createResp)

	// Example 2: Get collection details
	getResp, err := client.RAG.GetCollection(ctx, "my-docs")
	if err != nil {
		log.Fatalf("Failed to get collection: %v", err)
	}
	fmt.Printf("Collection details: %+v\n", getResp)

	// Example 3: List all collections
	listResp, err := client.RAG.ListCollections(ctx)
	if err != nil {
		log.Fatalf("Failed to list collections: %v", err)
	}
	fmt.Printf("All collections: %+v\n", listResp)

	// Example 4: Insert a document into the collection
	insertResp, err := client.RAG.InsertResource(ctx, "my-docs", "This is a sample document about machine learning.", wetro.ResourceTypeText)
	if err != nil {
		log.Fatalf("Failed to insert resource: %v", err)
	}
	fmt.Printf("Inserted resource: %+v\n", insertResp)

	// Example 5: Query the collection
	queryReq := wetro.QueryRequest{
		CollectionID: "my-docs",
		Query:        "What is machine learning?",
	}
	queryResp, err := client.RAG.QueryCollection(ctx, queryReq)
	if err != nil {
		log.Fatalf("Failed to query collection: %v", err)
	}
	fmt.Printf("Query results: %+v\n", queryResp)

	// Example 6: Chat with the collection
	chatReq := wetro.ChatRequest{
		CollectionID: "my-docs",
		Message:      "Can you explain machine learning in simple terms?",
	}
	chatResp, err := client.RAG.ChatWithCollection(ctx, chatReq)
	if err != nil {
		log.Fatalf("Failed to chat with collection: %v", err)
	}
	fmt.Printf("Chat response: %+v\n", chatResp)

	// Example 7: Remove a resource
	removeReq := wetro.ResourceDeleteRequest{
		CollectionID: "my-docs",
		ResourceID:   "resource-id",
	}
	removeResp, err := client.RAG.RemoveResource(ctx, removeReq)
	if err != nil {
		log.Fatalf("Failed to remove resource: %v", err)
	}
	fmt.Printf("Removed resource: %+v\n", removeResp)

	// Example 8: Delete the collection
	deleteResp, err := client.RAG.DeleteCollection(ctx, "my-docs")
	if err != nil {
		log.Fatalf("Failed to delete collection: %v", err)
	}
	fmt.Printf("Deleted collection: %+v\n", deleteResp)
}

func toolsExamples(client *wetro.Client, ctx context.Context) {
	fmt.Println("\n=== Tools Client Examples ===")

	// Example 1: Categorize text
	categorizeReq := wetro.CategorizeRequest{
		Type:     wetro.ResourceTypeText,
		Resource: "This is a technical document about artificial intelligence and machine learning.",
		JSONSchema: json.RawMessage(`{
			"type": "object",
			"properties": {
				"category": {"type": "string"},
				"confidence": {"type": "number"}
			}
		}`),
		Categories: []string{"Technology", "Science", "Education"},
		Prompt:     "Categorize this text based on its content",
	}
	categorizeResp, err := client.Tools.CategorizeData(ctx, categorizeReq)
	if err != nil {
		log.Fatalf("Failed to categorize data: %v", err)
	}
	fmt.Printf("Categorization results: %+v\n", categorizeResp)

	// Example 2: Generate text
	generateReq := wetro.TextGenerationRequest{
		Messages: []wetro.MessageObject{
			{
				Role:    "user",
				Content: "Write a short paragraph about the benefits of artificial intelligence",
			},
		},
		Model: wetro.GPT35Turbo,
	}
	generateResp, err := client.Tools.GenerateText(ctx, generateReq)
	if err != nil {
		log.Fatalf("Failed to generate text: %v", err)
	}
	fmt.Printf("Generated text: %+v\n", generateResp)

	// Example 3: Image to text
	imageReq := wetro.ImageToTextRequest{
		ImageURL: "https://example.com/ai-image.jpg",
		Query:    "Describe what you see in this image",
	}
	imageResp, err := client.Tools.ImageToText(ctx, imageReq)
	if err != nil {
		log.Fatalf("Failed to convert image to text: %v", err)
	}
	fmt.Printf("Image description: %+v\n", imageResp)

	// Example 4: Extract data from a webpage
	extractReq := wetro.DataExtractionRequest{
		WebURL: "https://example.com/article",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type": "string",
				},
				"author": map[string]interface{}{
					"type": "string",
				},
				"content": map[string]interface{}{
					"type": "string",
				},
				"publishDate": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}
	extractResp, err := client.Tools.ExtractData(ctx, extractReq)
	if err != nil {
		log.Fatalf("Failed to extract data: %v", err)
	}
	fmt.Printf("Extracted data: %+v\n", extractResp)
}
