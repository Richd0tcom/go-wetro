# WetroCloud API Wrapper for Go

A Go wrapper for interacting with the WetroCloud API, providing access to Retrieval Augmented Generation (RAG) and AI tools functionality.

## Installation

```bash
go get github.com/Richd0tcom/go-wetro
```

## Getting Started

To use the wrapper, import it in your Go code:

```go
import "github.com/Richd0tcom/go-wetro/wetro"
```

### Initialize the Client

```go
client := wetro.NewClient("your-api-key")
```

You can customize the client with options:

```go

// Custom HTTP client with timeout
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}


client := wetro.NewClient("your-api-key", 
    wetrocloud.WithHTTPClient(httpClient),
    wetrocloud.WithAPIVersion("v2"), // If you need a different API version
)
```

## Features

### RAG (Retrieval-Augmented Generation)

The RAG client provides functionality for managing collections and performing semantic search:

```go
ctx := context.Background()


// Create a collection
resp, err := client.RAG.CreateCollection(ctx, "my-docs")
if err != nil {
    log.Fatal(err)
}

// Get collection details
collection, err := client.RAG.GetCollection(ctx, "my-docs")
if err != nil {
    log.Fatal(err)
}

// List all collections
collections, err := client.RAG.ListCollections(ctx)
if err != nil {
    log.Fatal(err)
}

// Insert a resource
insertResp, err := client.RAG.InsertResource(ctx, "my-docs", "document text", wetro.ResourceTypeText)
if err != nil {
    log.Fatal(err)
}

// Query a collection
queryResp, err := client.RAG.QueryCollection(ctx, wetro.QueryRequest{
    CollectionID: "my-docs",
    Query:        "What is this about?",
})
if err != nil {
    log.Fatal(err)
}

// Chat with a collection
chatResp, err := client.RAG.ChatWithCollection(ctx, wetro.ChatRequest{
    CollectionID: "my-docs",
    Message:      "Explain this to me",
})
if err != nil {
    log.Fatal(err)
}

// Remove a resource
removeResp, err := client.RAG.RemoveResource(ctx, wetro.ResourceDeleteRequest{
    CollectionID: "my-docs",
    ResourceID:   "resource-id",
})
if err != nil {
    log.Fatal(err)
}

// Delete a collection
deleteResp, err := client.RAG.DeleteCollection(ctx, "my-docs")
if err != nil {
    log.Fatal(err)
}
```

### AI Tools

The Tools client provides various AI-powered utilities:

```go
ctx := context.Background()


// Categorize text
categorizeResp, err := client.Tools.CategorizeData(ctx, wetro.CategorizeRequest{
    Type:     wetro.ResourceTypeText,
    Resource: "This is a technical document",
    JSONSchema: json.RawMessage(`{
        "type": "object",
        "properties": {
            "category": {"type": "string"},
            "confidence": {"type": "number"}
        }
    }`),
    Categories: []string{"Technology", "Science"},
    Prompt:     "Categorize this text",
})
if err != nil {
    log.Fatal(err)
}

// Generate text
generateResp, err := client.Tools.GenerateText(ctx, wetro.TextGenerationRequest{
    Messages: []wetro.MessageObject{
        {
            Role:    "user",
            Content: "Write a short paragraph",
        },
    },
    Model: wetro.GPT35Turbo,
})
if err != nil {
    log.Fatal(err)
}

// Convert image to text
imageResp, err := client.Tools.ImageToText(ctx, wetro.ImageToTextRequest{
    ImageURL: "https://example.com/image.jpg",
    Query:    "Describe this image",
})
if err != nil {
    log.Fatal(err)
}

// Extract data from a webpage
extractResp, err := client.Tools.ExtractData(ctx, wetro.DataExtractionRequest{
    WebURL: "https://example.com",
    Schema: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "title":    map[string]interface{}{"type": "string"},
            "content":  map[string]interface{}{"type": "string"},
        },
    },
})
if err != nil {
    log.Fatal(err)
}
```

## Supported Resource Types

The SDK supports the following resource types:

```go
const (
    ResourceTypeText    ResourceType = "text"
    ResourceTypeWeb     ResourceType = "web"
    ResourceTypeFile    ResourceType = "file"
    ResourceTypeJSON    ResourceType = "json"
    ResourceTypeYouTube ResourceType = "youtube"
)
```
The WetroCloud API supports various resource types:

- **File**: Various file types including .csv, .docx, .epub, .hwp, .ipynb, .jpeg, .jpg, .mbox, .md, .mp3, .mp4, .pdf, .png, .ppt, .pptm, .pptx.
- **Text**: Plain text content.
- **JSON**: Structured data in JSON format.
- **Web**: Web-based resources, such as websites.
- **YouTube**: YouTube videos with YouTube URLs.
- **Audio**: Various audio file types including .3ga, .8svx, .aac, .ac3, .aif, .aiff, .alac, .amr, .ape, .au, .dss, .flac, .flv, .m4a, .m4b, .m4p, .m4r, .mp3, .mpga, .ogg, .oga, .mogg, .opus, .qcp, .tta, .voc, .wav, .wma, .wv.

## Available Chat Models

The SDK supports various chat models for text generation:

```go
const (
    ChatGPT4Latest            ChatModel = "chatgpt-4o-latest"
    Claude35Haiku20241022     ChatModel = "claude-3-5-haiku-20241022"
    Claude35Sonnet20240620    ChatModel = "claude-3-5-sonnet-20240620"
    Claude35Sonnet20241022    ChatModel = "claude-3-5-sonnet-20241022"
    GPT4TurboPreview          ChatModel = "gpt-4-turbo-preview"
    GPT45Preview              ChatModel = "gpt-4.5-preview"
    GPT4O                     ChatModel = "gpt-4o"
    GPT4OMini                 ChatModel = "gpt-4o-mini"
    // ... and more
)
```

## Error Handling

The SDK uses a custom error type for API errors:

```go
type APIError struct {
    Message    string
    StatusCode int 
    Payload    any
}
```

## Documentation
For more details, check out the official API documentation: [Wetrocloud Docs](https://docs.wetrocloud.com/introduction)

## License

This SDK is distributed under the terms of the MIT license.