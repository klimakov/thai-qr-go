package thaiqrgo

import (
	"fmt"
	"strconv"
	"strings"

	"thai-qr-go/internal"
)

// TLVTag represents a Tag-Length-Value structure.
type TLVTag struct {
	// ID is the tag identifier (2-digit hex string)
	ID string

	// Value is the tag value
	Value string

	// SubTags are nested TLV tags (if any)
	SubTags []TLVTag

	// Length is the length of the value (or encoded subtags)
	Length int
}

// Decode decodes a TLV string into an array of TLV tags.
//
// Returns an error if the payload format is invalid (e.g., incomplete tags).
func Decode(payload string) ([]TLVTag, error) {
	var tags []TLVTag

	idx := 0
	for idx < len(payload) {
		if idx+4 > len(payload) {
			return nil, fmt.Errorf("invalid TLV format: incomplete tag header at position %d", idx)
		}

		id := payload[idx : idx+2]
		lengthStr := payload[idx+2 : idx+4]
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid TLV format: invalid length at position %d: %w", idx+2, err)
		}

		if idx+4+length > len(payload) {
			return nil, fmt.Errorf("invalid TLV format: incomplete tag value at position %d (expected %d bytes, got %d)", idx+4, length, len(payload)-idx-4)
		}

		value := payload[idx+4 : idx+4+length]

		tags = append(tags, TLVTag{
			ID:     id,
			Length: length,
			Value:  value,
		})

		idx += 4 + length
	}

	return tags, nil
}

// Encode encodes an array of TLV tags into a TLV string.
func Encode(tags []TLVTag) string {
	var payload strings.Builder

	for _, tag := range tags {
		payload.WriteString(tag.ID)
		payload.WriteString(fmt.Sprintf("%02d", tag.Length))

		if len(tag.SubTags) > 0 {
			encodedSubTags := Encode(tag.SubTags)
			payload.WriteString(encodedSubTags)
			continue
		}

		payload.WriteString(tag.Value)
	}

	return payload.String()
}

// Checksum generates a CRC16 XMODEM checksum for the provided string.
//
// The checksum is returned as a 4-digit uppercase hexadecimal string by default.
func Checksum(payload string, upperCase bool) string {
	crc := internal.CRC16XMODEM(payload, 0xffff)
	sum := fmt.Sprintf("%04x", crc)
	if upperCase {
		sum = strings.ToUpper(sum)
	} else {
		sum = strings.ToLower(sum)
	}
	return sum
}

// WithCRCTag appends a CRC tag to the TLV payload.
//
// The function adds the CRC tag ID, length (always "04"), and calculated checksum.
func WithCRCTag(payload, crcTagID string, upperCase bool) string {
	payload += fmt.Sprintf("%02s", crcTagID)
	payload += "04"
	payload += Checksum(payload, upperCase)
	return payload
}

// Get finds a tag or sub-tag by tag ID in an array of TLV tags.
//
// If subTagID is provided, it searches for a sub-tag within the found tag.
// Returns nil if the tag is not found.
func Get(tags []TLVTag, tagID, subTagID string) *TLVTag {
	var tag *TLVTag
	for i := range tags {
		if tags[i].ID == tagID {
			tag = &tags[i]
			break
		}
	}

	if tag == nil {
		return nil
	}

	if subTagID != "" {
		for i := range tag.SubTags {
			if tag.SubTags[i].ID == subTagID {
				return &tag.SubTags[i]
			}
		}
		return nil
	}

	return tag
}

// Tag creates a new TLV tag with the specified ID and value.
//
// The length is automatically calculated from the value length.
func Tag(tagID, value string) TLVTag {
	return TLVTag{
		ID:     tagID,
		Value:  value,
		Length: len(value),
	}
}
