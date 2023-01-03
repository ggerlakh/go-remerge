package hashtool

import (
	"crypto/sha256"
	"fmt"
)

// Sha256 computing hash from input string
func Sha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}
