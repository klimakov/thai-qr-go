package thaiqrgo

import "testing"

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantLen int
		wantErr bool
	}{
		{
			name:    "simple tags",
			payload: "000411110104222202043333",
			wantLen: 3,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantLen {
				t.Errorf("Decode() len = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name    string
		tags    []TLVTag
		want    string
		wantErr bool
	}{
		{
			name: "simple tags",
			tags: []TLVTag{
				{ID: "00", Value: "1111", Length: 4},
				{ID: "01", Value: "2222", Length: 4},
			},
			want:    "0004111101042222",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.tags)
			if got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTag(t *testing.T) {
	tag := Tag("00", "test")
	if tag.ID != "00" || tag.Value != "test" || tag.Length != 4 {
		t.Errorf("Tag() = %+v, want ID=00 Value=test Length=4", tag)
	}
}

func TestGet(t *testing.T) {
	tags := []TLVTag{
		{ID: "00", Value: "1111", Length: 4},
		{ID: "01", Value: "2222", Length: 4},
		{
			ID: "29", Value: "encoded", Length: 7,
			SubTags: []TLVTag{
				{ID: "01", Value: "subvalue", Length: 8},
			},
		},
	}

	got := Get(tags, "01", "")
	if got == nil || got.Value != "2222" {
		t.Errorf("Get() = %v, want Value=2222", got)
	}

	got = Get(tags, "29", "01")
	if got == nil || got.Value != "subvalue" {
		t.Errorf("Get() subTag = %v, want Value=subvalue", got)
	}

	// Test non-existent tag
	got = Get(tags, "99", "")
	if got != nil {
		t.Errorf("Get() with non-existent tag should return nil, got %v", got)
	}

	// Test non-existent subTag
	got = Get(tags, "29", "99")
	if got != nil {
		t.Errorf("Get() with non-existent subTag should return nil, got %v", got)
	}

	// Test tag without subTags
	got = Get(tags, "00", "01")
	if got != nil {
		t.Errorf("Get() with subTagID on tag without subTags should return nil, got %v", got)
	}
}

func TestDecode_EdgeCases(t *testing.T) {
	// Test invalid length (too short)
	_, err := Decode("000")
	if err == nil {
		t.Error("Decode() should return error for incomplete tag header")
	}

	// Test invalid length value (non-numeric)
	_, err = Decode("00AB1234")
	if err == nil {
		t.Error("Decode() should return error for invalid length")
	}

	// Test incomplete value
	_, err = Decode("000412")
	if err == nil {
		t.Error("Decode() should return error for incomplete tag value")
	}

	// Test empty payload
	got, err := Decode("")
	if err != nil {
		t.Errorf("Decode() empty payload should not error, got %v", err)
	}
	if len(got) != 0 {
		t.Errorf("Decode() empty payload should return empty slice, got %v", got)
	}
}

func TestEncode_EdgeCases(t *testing.T) {
	// Test empty tags
	got := Encode([]TLVTag{})
	if got != "" {
		t.Errorf("Encode() empty tags should return empty string, got %v", got)
	}

	// Test with subTags
	tags := []TLVTag{
		{
			ID: "29",
			SubTags: []TLVTag{
				{ID: "01", Value: "value1", Length: 6},
				{ID: "02", Value: "value2", Length: 6},
			},
			Length: 20,
		},
	}
	got = Encode(tags)
	if len(got) == 0 {
		t.Error("Encode() with subTags should return non-empty string")
	}
}

func TestChecksum(t *testing.T) {
	tests := []struct {
		name      string
		payload   string
		upperCase bool
	}{
		{
			name:      "uppercase",
			payload:   "test",
			upperCase: true,
		},
		{
			name:      "lowercase",
			payload:   "test",
			upperCase: false,
		},
		{
			name:      "empty string uppercase",
			payload:   "",
			upperCase: true,
		},
		{
			name:      "empty string lowercase",
			payload:   "",
			upperCase: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Checksum(tt.payload, tt.upperCase)
			if len(got) != 4 {
				t.Errorf("Checksum() length = %v, want 4", len(got))
			}
			if tt.upperCase {
				for _, c := range got {
					if c >= 'a' && c <= 'f' {
						t.Errorf("Checksum() with upperCase=true should return uppercase, got %v", got)
						break
					}
				}
			}
		})
	}
}

func TestWithCRCTag(t *testing.T) {
	payload := "00020101021129370016A000000677010111011300668012345675802TH5303764"
	got := WithCRCTag(payload, "63", true)

	// Should append tag ID, length (04), and 4-digit CRC
	if len(got) != len(payload)+8 {
		t.Errorf("WithCRCTag() length = %v, want %v", len(got), len(payload)+8)
	}

	// Should end with CRC (4 hex digits)
	if len(got) >= 4 {
		crcPart := got[len(got)-4:]
		for _, c := range crcPart {
			if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F')) {
				t.Errorf("WithCRCTag() CRC part should be hex uppercase, got %v", crcPart)
				break
			}
		}
	}

	// Test lowercase
	gotLower := WithCRCTag(payload, "63", false)
	// Note: They may be equal if CRC contains only digits (0-9)
	// So we just verify both are valid strings
	if len(gotLower) != len(got) {
		t.Errorf("WithCRCTag() with upperCase=false length = %v, want %v", len(gotLower), len(got))
	}
}
