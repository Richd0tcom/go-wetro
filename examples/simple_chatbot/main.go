package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Richd0tcom/wetrocloud-sdk-go/wetrocloud"
)

func main() {
	// Get API credentials from environment
	apiKey := os.Getenv("WETROCLOUD_API_KEY")
	if apiKey == "" {
		log.Fatal("WETROCLOUD_API_KEY environment variable is required")
	}

	// Initialize the client
	client := wetrocloud.NewClient(
		"https://api.wetrocloud.com",
		apiKey,
	)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create a new collection for the chatbot knowledge base
	fmt.Println("Creating chatbot knowledge base...")
	collection, err := client.CreateCollection(ctx, "Simple Chatbot Knowledge Base")
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	collectionID := collection.CollectionID
	fmt.Printf("Created collection with ID: %s\n", collectionID)

	// Add knowledge to the chatbot
	addKnowledgeToChatbot(ctx, client, collectionID)

	// Start the chat loop
	fmt.Println("\n=== WetroCloud Chatbot ===")
	fmt.Println("Type 'exit' to quit.")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()
		if strings.ToLower(userInput) == "exit" {
			break
		}

		// Query the chatbot
		response, err := queryChatbot(ctx, client, collectionID, userInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Chatbot: %s\n", response)
	}

	// Clean up resources
	fmt.Println("\nCleaning up resources...")
	_, err = client.DeleteCollection(ctx, collectionID)
	if err != nil {
		log.Printf("Warning: Failed to delete collection: %v", err)
	} else {
		fmt.Println("Successfully deleted collection.")
	}
}

// addKnowledgeToChatbot adds knowledge to the chatbot collection
func addKnowledgeToChatbot(ctx context.Context, client *wetrocloud.Client, collectionID string) {
	// Add company information
	companyInfo := `
	WetroCloud is a cloud platform that provides RAG (Retrieval Augmented Generation) and LLM agent APIs.
	The company was founded in 2023 and is headquartered in San Francisco, California.
	WetroCloud's mission is to make AI accessible to developers of all skill levels.
	The company offers a range of pricing plans from free tier to enterprise solutions.
	`

	fmt.Println("Adding company information...")
	_, err := client.TextResource(ctx, collectionID, companyInfo, map[string]interface{}{
		"category": "company",
		"tags":     []string{"about", "information"},
	})
	if err != nil {
		log.Fatalf("Failed to add company information: %v", err)
	}

	// Add product information
	productInfo := `
	WetroCloud's main product is a comprehensive API for building RAG applications.
	Features include:
	1. Document processing and ingestion
	2. Semantic search and retrieval
	3. LLM integration
	4. Custom agents
	5. Categorization tools
	
	The platform supports multiple resource types including text, files, web content, 
	JSON data, YouTube videos, and audio files.
	`

	fmt.Println("Adding product information...")
	_, err = client.TextResource(ctx, collectionID, productInfo, map[string]interface{}{
		"category": "product",
		"tags":     []string{"features", "capabilities"},
	})
	if err != nil {
		log.Fatalf("Failed to add product information: %v", err)
	}

	// Add FAQ information
	faqInfo := `
	Frequently Asked Questions:
	
	Q: How do I get started with WetroCloud?
	A: Sign up for a free account at wetrocloud.com, get your API key, and start using our SDKs.
	
	Q: What programming languages do you support?
	A: We have official SDKs for JavaScript, Python, Go, and Java.
	
	Q: Is there a free tier?
	A: Yes, we offer a generous free tier with up to 1000 API calls per month.
	
	Q: How is billing calculated?
	A: Billing is based on the number of API calls and tokens processed.
	
	Q: Do you provide custom solutions?
	A: Yes, our enterprise tier includes custom solution development and dedicated support.
	`

	fmt.Println("Adding FAQ information...")
	_, err = client.TextResource(ctx, collectionID, faqInfo, map[string]interface{}{
		"category": "support",
		"tags":     []string{"faq", "help"},
	})
	if err != nil {
		log.Fatalf("Failed to add FAQ information: %v", err)
	}
}

// queryChatbot queries the chatbot with user input and returns the response
func queryChatbot(ctx context.Context, client *wetrocloud.Client, collectionID, query string) (string, error) {
	queryReq := &wetrocloud.QueryRequest{
		CollectionID: collectionID,
		Query:        query,
		MaxResults:   wetrocloud.IntPtr(3),
	}

	resp, err := client.QueryCollection(ctx, queryReq)
	if err != nil {
		return "", err
	}

	return resp.Response, nil
}
