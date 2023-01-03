package jsontool

import (
	"bytes"
	"encoding/json"
)

// ExtenedMarshal func for HTML chars escaping in JSON marshaling
func ExtenedMarshal(data interface{}, prefix string, indent string, EscapeHTML bool) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(EscapeHTML)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(data)
	return buffer.Bytes(), err
}
