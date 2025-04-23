package wetrocloud

import "encoding/json"

// ChatModel represents the available chat models supported by the API.
// Each model has specific capabilities and performance characteristics.
type ChatModel string

const (
	ChatGPT4Latest            ChatModel = "chatgpt-4o-latest"
	Claude35Haiku20241022     ChatModel = "claude-3-5-haiku-20241022"
	Claude35Sonnet20240620    ChatModel = "claude-3-5-sonnet-20240620"
	Claude35Sonnet20241022    ChatModel = "claude-3-5-sonnet-20241022"
	Claude37Sonnet20250219    ChatModel = "claude-3-7-sonnet-20250219"
	Claude3Haiku20240307      ChatModel = "claude-3-haiku-20240307"
	Claude3Opus20240229       ChatModel = "claude-3-opus-20240229"
	Claude3Sonnet20240229     ChatModel = "claude-3-sonnet-20240229"
	DeepseekR1DistillLlama70b ChatModel = "deepseek-r1-distill-llama-70b"
	GPT35Turbo                ChatModel = "gpt-3.5-turbo"
	GPT4                      ChatModel = "gpt-4"
	GPT4Turbo                 ChatModel = "gpt-4-turbo"
	GPT4TurboPreview          ChatModel = "gpt-4-turbo-preview"
	GPT45Preview              ChatModel = "gpt-4.5-preview"
	GPT4O                     ChatModel = "gpt-4o"
	GPT4OMini                 ChatModel = "gpt-4o-mini"
	Llama318B                 ChatModel = "llama-3.1-8b"
	Llama318BInstant          ChatModel = "llama-3.1-8b-instant"
	Llama321BPreview          ChatModel = "llama-3.2-1b-preview"
	Llama323BPreview          ChatModel = "llama-3.2-3b-preview"
	Llama3211BVisionPreview   ChatModel = "llama-3.2-11b-vision-preview"
	Llama3290BVisionPreview   ChatModel = "llama-3.2-90b-vision-preview"
	Llama3370B                ChatModel = "llama-3.3-70b"
	Llama3370BSpecDec         ChatModel = "llama-3.3-70b-specdec"
	Llama3370BVersatile       ChatModel = "llama-3.3-70b-versatile"
	Llama370B8192             ChatModel = "llama3-70b-8192"
	Llama38B8192              ChatModel = "llama3-8b-8192"
	LlamaGuard38B             ChatModel = "llama-guard-3-8b"
	Mixtral8x7B32768          ChatModel = "mixtral-8x7b-32768"
	O1                        ChatModel = "o1"
	O1Mini                    ChatModel = "o1-mini"
	O1Preview                 ChatModel = "o1-preview"
	O3Mini                    ChatModel = "o3-mini"
	Qwen25_32B                ChatModel = "qwen-2.5-32b"
	Qwen25Coder32B            ChatModel = "qwen-2.5-coder-32b"
)

// ResourceType represents the type of resource that can be processed by the API.
// Different resource types may have different processing requirements and capabilities.
type ResourceType string

const (
	ResourceTypeText    ResourceType = "text"
	ResourceTypeWeb     ResourceType = "web"
	ResourceTypeFile    ResourceType = "file"
	ResourceTypeJSON    ResourceType = "json"
	ResourceTypeYouTube ResourceType = "youtube"
)

// StandardResponse represents the standard API response structure used across most endpoints.
// It includes success status, token usage, and the actual response data.
type StandardResponse struct {
	Success  bool `json:"success"`
	Tokens   int  `json:"tokens"`
	Response any  `json:"response,omitempty"`
}

// CollectionCreateResponse contains the response from creating a collection.
// It includes the success status and the ID of the created collection.
type CollectionCreateResponse struct {
	Success      bool   `json:"success"`
	CollectionID string `json:"collection_id"`
}

// CollectionItem represents a collection in the WetroCloud system.
// It contains basic information about a collection.
type CollectionItem struct {

	// The unique identifier of the collection
	CollectionID string `json:"collection_id"`

	// The timestamp when the collection was created
	CreatedAt string `json:"created_at"`
}

// GetCollectionResponse contains the response from retrieving a collection.
// It includes success status, whether the collection was found, and its ID.
type GetCollectionResponse struct {
	Success      bool   `json:"success"`
	Found        bool   `json:"found"`
	CollectionID string `json:"collection_id"`
}

// CollectionsResponse contains the response from listing collections
type ListCollectionResponse struct {
	//Total number of collections
	Count int `json:"count"`

	//URL for the next pagination item
	Next string `json:"next"`

	//URL for the previous pagination item.
	Previous string `json:"previous"`

	//A list of all available collections.
	Results []CollectionItem `json:"results"`
}

// ResourceInsertRequest represents a request to insert a resource into a collection.
// It specifies the collection, resource type, and the actual resource content.
type ResourceInsertRequest struct {
	CollectionID string       `json:"collection_id"`
	Type         ResourceType `json:"type"`
	Resource     string       `json:"resource"`
}

func (r *ResourceInsertRequest) Validate(v *validator) bool {
	v.Check(r.CollectionID != "", "collection_id", "collection_id should not be empty")
	v.Check(r.Type != "", "type", "resource type should not be empty")
	return v.Valid()
}

// ResourceInsertResponse contains the response from inserting a resource.
// It includes status, the ID of the inserted resource, and token usage.
type ResourceInsertResponse struct {
	ResourceID string `json:"resource_id"`
	Success    bool   `json:"success"`
	Tokens     int    `json:"tokens,omitempty"`
}

// ResourceDeleteRequest represents a request to remove a resource from a collection.
type ResourceDeleteRequest struct {
	CollectionID string `json:"collection_id"`
	ResourceID   string `json:"resource_id"`
}

type ResourceDeleteResponse struct {
	Success bool `json:"success"`
}

// QueryRequest represents a request to query a collection
type QueryRequest struct {
	CollectionID string `json:"collection_id"`
	Query        string `json:"request_query"`

	// (optional) The model to use for processing
	Model ChatModel `json:"model,omitempty"`

	//(Optional) JSON schema for response formatting
	JSONSchema json.RawMessage `json:"json_schema,omitempty"`

	// (Optional) rules for schema validation.
	//must be present if JSONSchema is used
	JSONSchemaRules json.RawMessage `json:"json_schema_rules,omitempty"`

	Stream bool `json:"stream"`
}

func (r *QueryRequest) Validate(v *validator) bool {
	v.Check(r.CollectionID != "", "collection_id", "collection_id should not be empty")

	if len(r.JSONSchema) > 0 {
		v.Check(len(r.JSONSchemaRules) > 0, "json_schema_rules", "must have json_schema_rules if json_schema is used")
	}
	
	return v.Valid()
}

// Message represents a single message in a chat conversation.
// It's a map of string key-value pairs for flexibility.
type Message map[string]string

// ChatRequest represents a request to chat with a collection.
// It supports conversation history and streaming options.
type ChatRequest struct {
	CollectionID string    `json:"collection_id"`
	Message      string    `json:"message"`
	ChatHistory  []Message `json:"chat_history"`
	Stream       bool      `json:"stream"`
}

type DeleteCollectionResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// CategorizeRequest represents a request to categorize data
// It supports custom schemas and categories for classification.
type CategorizeRequest struct {
	Type     ResourceType `json:"type"`
	Resource string       `json:"resource"`

	// Schema for the categorization result
	JSONSchema json.RawMessage `json:"json_schema"`

	// Available categories for classification
	Categories []string `json:"categories"`

	// Custom prompt for categorization
	Prompt string `json:"prompt"`
}

type MessageObject struct {

	// The role of the message sender (e.g., "user", "assistant")
	Role string `json:"role"`

	// The actual message content
	Content string `json:"content"`
}

// TextGenerationRequest represents a request to generate text.
type TextGenerationRequest struct {
	Messages []MessageObject `json:"messages"`
	Model    ChatModel       `json:"model,omitempty"`
}

// ImageToTextRequest represents a request to generate text from an image.
type ImageToTextRequest struct {
	ImageURL string `json:"image_url"`
	Query    string `json:"request_query"`
}

// DataExtractionRequest represents a request to extract data from a website.
// It includes the website URL and a schema for the expected data structure.
type DataExtractionRequest struct {
	WebURL string `json:"website"`
	Schema any    `json:"json_schema"`
}

// Errors
// APIError represents an error from the Wetrocloud API
type APIError struct {
	Message    string
	StatusCode int
	Payload    any
}

func (e APIError) Error() string {
	return e.Message
}
