package wetrocloud

import (
	"encoding/json"
)

// IntPtr returns a pointer to the int value passed in
func IntPtr(v int) *int {
	return &v
}

// FloatPtr returns a pointer to the float64 value passed in
func FloatPtr(v float64) *float64 {
	return &v
}

// StringPtr returns a pointer to the string value passed in
func StringPtr(v string) *string {
	return &v
}

// BoolPtr returns a pointer to the bool value passed in
func BoolPtr(v bool) *bool {
	return &v
}


func ToJSONSchema(schema any) (string, error) {
	b, err:= json.Marshal(schema)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
