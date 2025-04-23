package wetrocloud

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func ToJSONSchema(schema any) (string, error) {
	b, err := json.Marshal(schema)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid, nil
}

func parseError(resp *http.Response) string {
	var errorData map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&errorData); err != nil {
		return "Unknown error"
	}

	if error, ok := errorData["error"].(string); ok {
		return error
	}
	if detail, ok := errorData["detail"].(string); ok {
		return detail
	}

	var errors []string
	for field, messages := range errorData {
		if msgList, ok := messages.([]any); ok {
			var msgStrings []string
			for _, msg := range msgList {
				msgStrings = append(msgStrings, fmt.Sprintf("%v", msg))
			}
			errors = append(errors, fmt.Sprintf("%s: %s", field, strings.Join(msgStrings, ", ")))
		} else {
			errors = append(errors, fmt.Sprintf("%s: %v", field, messages))
		}
	}
	return strings.Join(errors, "; ")
}

type errorFields map[string]string

// ValidationError Custom error type for better error handling.
type ValidationError struct {
	Fields  map[string]string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

type validator struct {
	errors errorFields
}

// New returns a new Validator instance.
func NewValidator() *validator {
	return &validator{errors: make(errorFields)}
}

func NewValidationError(message string, fields errorFields) *ValidationError {
	return &ValidationError{
		Fields:  fields,
		Message: message,
	}
}

// NewWithStore returns a new Validator instance with a store.
// func NewWithStore(ctx context.Context, store db.Store) *Validator {
//	return &Validator{
//		Errors: make(ErrorFields),
//		store:  store,
//		ctx:    ctx,
//	}
//}

// Valid returns true if the errors map doesn't contain any entries.
func (v *validator) Valid() bool {
	return len(v.errors) == 0
}

// AddError adds an error message to the map (so long as no entry already exists for
// the given key).
func (v *validator) AddError(key, message string) {
	if _, exists := v.errors[key]; !exists {
		v.errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
