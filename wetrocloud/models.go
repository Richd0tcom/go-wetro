package wetrocloud

// StandardResponse represents the standard API response structure
type StandardResponse struct {
	Success  bool `json:"success"`
	Tokens   int  `json:"tokens"`
	Response any  `json:"response,omitempty"`
}

// CollectionCreateResponse contains the response from creating a collection
type CollectionCreateResponse struct {
	Success      bool   `json:"success"`
	CollectionID string `json:"collection_id"`
}

// Collection represents a collection in the WetroCloud system
type Collection struct {
	CollectionID string `json:"colection_id"`
	CreatedAt    string `json:"created_at"`
}

// CollectionsResponse contains the response from listing collections
type CollectionsResponse struct {
	//Total number of collections
	Count int `json:"count"`

	//URL for the next pagination item
	Next string `json:"next"`

	//URL for the previous pagination item.
	Previous string `json:"previous"`

	//A list of all available collections.
	Results []Collection `json:"results"`
}

// ResourceInsertRequest represents a request to insert a resource
type ResourceInsertRequest struct {
	CollectionID string `json:"collection_id"`
	Type         string `json:"type"`
	Resource     string `json:"resource"`
}

// contains the response from inserting a resource
type ResourceInsertResponse struct {
	ResourceID string `json:"resource_id"`
	Success    bool   `json:"success"`
	Tokens     int    `json:"tokens"`
}

type ResourceDeleteRequest struct {
	CollectionID string `json:"collection_id"`
	ResourceID   string `json:"resource_id"`
}

// QueryRequest represents a request to query a collection
type QueryRequest struct {
	CollectionID string `json:"collection_id"`
	Query        string `json:"request_query"`
	Model        string `json:"model,omitempty"`
}

type ChatRequest struct {
	CollectionID string `json:"collection_id"`
	Message      string `json:"message"`
	ChatHistory  string `json:"chat_history"`
}

// CategorizeRequest represents a request to categorize data
type CategorizeRequest struct {
	Type       string `json:"type"`
	Resource   string `json:"resource"`
	JSONSchema string `json:"json_schema"`
	Categories string `json:"categories"`
	Prompt     string `json:"prompt"`
}

type MessageObject struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TextGenerationRequest struct {
	Messages []MessageObject `json:"messages"`
	Model    string          `json:"model,omitempty"`
}

type ImageToFreeTextRequest struct {
	ImageURL string `json:"image_url"`
	Query    string `json:"request_query"`
}

type DataExtractionRequest struct {
	WebURL string `json:"website"`
	Schema any `json:"json_schema"`
}
