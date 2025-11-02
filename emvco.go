package thaiqrgo

// EMVCoQR represents a parsed EMVCo QR code with its TLV tags.
type EMVCoQR struct {
	payload string
	tags    []TLVTag
}

// GetTag retrieves a tag or sub-tag by ID.
//
// If subTagID is provided, it searches for a sub-tag within the found tag.
// Returns nil if the tag is not found.
func (q *EMVCoQR) GetTag(tagID, subTagID string) *TLVTag {
	return Get(q.tags, tagID, subTagID)
}

// GetTagValue retrieves the value of a tag or sub-tag by ID.
//
// Returns an empty string if the tag is not found.
func (q *EMVCoQR) GetTagValue(tagID, subTagID string) string {
	tag := q.GetTag(tagID, subTagID)
	if tag == nil {
		return ""
	}
	return tag.Value
}

// GetTags returns all TLV tags in the QR code.
func (q *EMVCoQR) GetTags() []TLVTag {
	return q.tags
}

// GetPayload returns the original payload string.
func (q *EMVCoQR) GetPayload() string {
	return q.payload
}

// Validate validates the QR code payload by recalculating the CRC checksum.
//
// It filters out the CRC tag, recalculates the checksum, and compares it with the original payload.
func (q *EMVCoQR) Validate(crcTagID string) bool {
	// Filter out the CRC tag
	var tagsWithoutCRC []TLVTag
	for _, tag := range q.tags {
		if tag.ID != crcTagID {
			tagsWithoutCRC = append(tagsWithoutCRC, tag)
		}
	}

	expected := WithCRCTag(Encode(tagsWithoutCRC), crcTagID, true)
	return q.payload == expected
}

