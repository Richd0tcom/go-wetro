package wetrocloud

// StandardResponse represents the standard API response structure
type StandardResponse struct {
	Success  bool   `json:"success"`
	Tokens   int    `json:"tokens"`
	Response string `json:"response,omitempty"`
}

// CollectionCreateResponse contains the response from creating a collection
type CollectionCreateResponse struct {
	StandardResponse
	CollectionID string `json:"collection_id"`
}

// Collection represents a collection in the WetroCloud system
type Collection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CollectionsResponse contains the response from listing collections
type CollectionsResponse struct {
	StandardResponse
	Collections []Collection `json:"collections"`
}

// Resource types supported by the API
const (
	ResourceTypeFile    = "file"
	ResourceTypeText    = "text"
	ResourceTypeJSON    = "json"
	ResourceTypeWeb     = "web"
	ResourceTypeYoutube = "youtube"
	ResourceTypeAudio   = "audio"
)

// ResourceInsertRequest represents a request to insert a resource
type ResourceInsertRequest struct {
	CollectionID string      `json:"collection_id"`
	Type         string      `json:"type"`
	Source       interface{} `json:"source"`
	Metadata     interface{} `json:"metadata,omitempty"`
}

// QueryRequest represents a request to query a collection
type QueryRequest struct {
	CollectionID string   `json:"collection_id"`
	Query        string   `json:"query"`
	MaxResults   *int     `json:"max_results,omitempty"`
	Confidence   *float64 `json:"confidence,omitempty"`
}

type ChatRequest struct {
	CollectionID string   `json:"collection_id"`
}

// CategorizeRequest represents a request to categorize data
type CategorizeRequest struct {
	Data   interface{} `json:"data"`
	Schema interface{} `json:"schema"`
}
