package jsontool

import (
	"bytes"
	"encoding/json"
)

// ExtendedMarshal func for HTML chars escaping in JSON marshaling
func ExtendedMarshal(data interface{}, prefix string, indent string, EscapeHTML bool) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(EscapeHTML)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(data)
	return buffer.Bytes(), err
}
