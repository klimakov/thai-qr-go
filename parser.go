package thaiqrgo

import (
	"fmt"
	"regexp"
	"strings"
)

var tlvPattern = regexp.MustCompile(`^\d{4}.+`)

// Parse parses an EMVCo-compatible QR code data string.
//
// Parameters:
//   - payload: QR code data string from the scanner
//   - strict: If true, validates CRC checksum before parsing
//   - subTags: If true, parses TLV sub-tags recursively
//
// Returns an EMVCoQR instance with TLV tags, or an error if parsing fails.
func Parse(payload string, strict, subTags bool) (*EMVCoQR, error) {
	if !tlvPattern.MatchString(payload) {
		return nil, fmt.Errorf("invalid QR code format: payload must start with 4 digits")
	}

	if strict {
		if len(payload) < 4 {
			return nil, fmt.Errorf("invalid QR code format: payload too short for checksum validation")
		}
		expected := strings.ToUpper(payload[len(payload)-4:])
		calculated := Checksum(payload[:len(payload)-4], true)
		if expected != calculated {
			return nil, fmt.Errorf("invalid CRC checksum: expected %s, got %s", expected, calculated)
		}
	}

	tags, err := Decode(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode TLV: %w", err)
	}

	if len(tags) == 0 {
		return nil, fmt.Errorf("no tags found in payload")
	}

	if subTags {
		for i := range tags {
			if len(tags[i].Value) >= 4 && tlvPattern.MatchString(tags[i].Value) {
				sub, err := Decode(tags[i].Value)
				if err != nil {
					continue
				}

				// Validate that all decoded sub-tags are valid (length matches value length)
				valid := true
				for _, subTag := range sub {
					if subTag.Length != len(subTag.Value) {
						valid = false
						break
					}
				}

				if valid && len(sub) > 0 {
					tags[i].SubTags = sub
				}
			}
		}
	}

	return &EMVCoQR{
		payload: payload,
		tags:    tags,
	}, nil
}

// ParseBarcode parses a BOT Barcode data string.
//
// Returns a BOTBarcode instance, or an error if parsing fails.
func ParseBarcode(payload string) (*BOTBarcode, error) {
	return BOTBarcodeFromString(payload)
}

