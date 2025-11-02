package internal

import (
	"fmt"
)

// EncodeTag81 generates a UCS-2-like hex string for Tag 81.
//
// This method is equivalent to:
// Buffer.from(message, 'utf16le').swap16().toString('hex').toUpperCase()
//
// It encodes each character's Unicode code point as a 4-digit uppercase hex string.
// This function is exported for use within the module.
func EncodeTag81(message string) string {
	var result string
	for _, c := range message {
		result += fmt.Sprintf("%04X", c)
	}
	return result
}
