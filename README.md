# WetroCloud SDK for Go

A Go SDK for interacting with the WetroCloud API, which allows you to build, scale, and deploy with Retrieval Augmented Generation (RAG) and LLM agents.

## Installation

```bash
go get github.com/Richd0tcom/wetrocloud-sdk-go
```

## Getting Started

To use the SDK, import it in your Go code:

```go
import "github.com/Richd0tcom/wetrocloud-sdk-go/wetrocloud"
```

### Initialize the Client

```go
client := wetrocloud.NewClient(
    "https://api.wetrocloud.com",
    "your-api-key",
)
```

You can also customize the client with options:

```go
// Custom HTTP client with timeout
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}

client := wetrocloud.NewClient(
    "https://api.wetrocloud.com",
    "your-api-key",
    wetrocloud.WithHTTPClient(httpClient),
    wetrocloud.WithAPIVersion("v2"), // If you need a different API version
)
```

## Usage Examples

### Create a Collection

```go
ctx := context.Background()
collection, err := client.CreateCollection(ctx, "My Collection")
if err != nil {
    log.Fatalf("Failed to create collection: %v", err)
}

fmt.Printf("Collection created with ID: %s\n", collection.CollectionID)
```

### Get All Collections

```go
ctx := context.Background()
collections, err := client.GetAllCollections(ctx)
if err != nil {
    log.Fatalf("Failed to get collections: %v", err)
}

for _, collection := range collections.Collections {
    fmt.Printf("Collection: %s (ID: %s)\n", collection.Name, collection.ID)
}
```

### Add Resources to a Collection

#### Text Resource

```go
ctx := context.Background()
resp, err := client.TextResource(ctx, collectionID, "This is a sample text", nil)
if err != nil {
    log.Fatalf("Failed to add text resource: %v", err)
}
```

#### File Resource

```go
ctx := context.Background()
fileURL := "https://example.com/document.pdf"
resp, err := client.FileResource(ctx, collectionID, fileURL, nil)
if err != nil {
    log.Fatalf("Failed to add file resource: %v", err)
}
```

#### Web Resource

```go
ctx := context.Background()
webURL := "https://example.com/article"
resp, err := client.WebResource(ctx, collectionID, webURL, nil)
if err != nil {
    log.Fatalf("Failed to add web resource: %v", err)
}
```

#### JSON Resource

```go
ctx := context.Background()
jsonData := map[string]interface{}{
    "title": "Example Document",
    "content": "This is the content of the document",
}
resp, err := client.JSONResource(ctx, collectionID, jsonData, nil)
if err != nil {
    log.Fatalf("Failed to add JSON resource: %v", err)
}
```

#### YouTube Resource

```go
ctx := context.Background()
youtubeURL := "https://www.youtube.com/watch?v=example"
resp, err := client.YoutubeResource(ctx, collectionID, youtubeURL, nil)
if err != nil {
    log.Fatalf("Failed to add YouTube resource: %v", err)
}
```

#### Audio Resource

```go
ctx := context.Background()
audioURL := "https://example.com/audio.mp3"
resp, err := client.AudioResource(ctx, collectionID, audioURL, nil)
if err != nil {
    log.Fatalf("Failed to add audio resource: %v", err)
}
```

### Query a Collection

```go
ctx := context.Background()
queryReq := &wetrocloud.QueryRequest{
    CollectionID: collectionID,
    Query:        "What is this document about?",
    MaxResults:   wetrocloud.IntPtr(5),   // Optional
    Confidence:   wetrocloud.FloatPtr(0.7), // Optional
}

resp, err := client.QueryCollection(ctx, queryReq)
if err != nil {
    log.Fatalf("Failed to query collection: %v", err)
}

fmt.Printf("Response: %s\n", resp.Response)
```

### Delete a Collection

```go
ctx := context.Background()
resp, err := client.DeleteCollection(ctx, collectionID)
if err != nil {
    log.Fatalf("Failed to delete collection: %v", err)
}
```

### Categorize Data

```go
ctx := context.Background()
data := "This is a sample text that needs to be categorized"
schema := map[string]interface{}{
    "categories": []string{"Business", "Technology", "Science"},
}

req := &wetrocloud.CategorizeRequest{
    Data:   data,
    Schema: schema,
}

resp, err := client.Categorize(ctx, req)
if err != nil {
    log.Fatalf("Failed to categorize data: %v", err)
}

fmt.Printf("Categorization response: %s\n", resp.Response)
```

## Supported Resource Types

The WetroCloud API supports various resource types:

- **File**: Various file types including .csv, .docx, .epub, .hwp, .ipynb, .jpeg, .jpg, .mbox, .md, .mp3, .mp4, .pdf, .png, .ppt, .pptm, .pptx.
- **Text**: Plain text content.
- **JSON**: Structured data in JSON format.
- **Web**: Web-based resources, such as websites.
- **YouTube**: YouTube videos with YouTube URLs.
- **Audio**: Various audio file types including .3ga, .8svx, .aac, .ac3, .aif, .aiff, .alac, .amr, .ape, .au, .dss, .flac, .flv, .m4a, .m4b, .m4p, .m4r, .mp3, .mpga, .ogg, .oga, .mogg, .opus, .qcp, .tta, .voc, .wav, .wma, .wv.

## License

This SDK is distributed under the terms of the MIT license.