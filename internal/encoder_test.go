package internal

import "testing"

func TestEncodeTag81(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "simple message",
			message: "Hi",
			want:    "00480069",
		},
		{
			name:    "empty string",
			message: "",
			want:    "",
		},
		{
			name:    "Hello World",
			message: "Hello",
			want:    "00480065006C006C006F",
		},
		{
			name:    "unicode characters",
			message: "ส",
			want:    "0E2A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EncodeTag81(tt.message)
			if got != tt.want {
				t.Errorf("EncodeTag81(%q) = %v, want %v", tt.message, got, tt.want)
			}
		})
	}
}
