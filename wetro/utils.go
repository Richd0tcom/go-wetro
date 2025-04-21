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
		if msgList, ok := messages.([]interface{}); ok {
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
