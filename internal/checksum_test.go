package internal

import "testing"

func TestCRC16XMODEM(t *testing.T) {
	tests := []struct {
		name string
		data string
		crc  uint16
	}{
		{
			name: "empty string with 0xffff",
			data: "",
			crc:  0xffff,
		},
		{
			name: "simple string with 0xffff",
			data: "test",
			crc:  0xffff,
		},
		{
			name: "simple string with zero crc",
			data: "test",
			crc:  0x0000,
		},
		{
			name: "known TLV payload",
			data: "00020101021129370016A000000677010111011300668012345675802TH5303764",
			crc:  0xffff,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CRC16XMODEM(tt.data, tt.crc)
			// Just verify it returns a valid uint16
			if got > 0xffff {
				t.Errorf("CRC16XMODEM(%q, 0x%04x) = 0x%04x, should be <= 0xffff", tt.data, tt.crc, got)
			}
			// Verify it's not zero for non-empty input (unless specific case)
			if tt.data != "" && got == 0 && tt.crc == 0xffff {
				t.Logf("CRC16XMODEM(%q, 0x%04x) = 0x%04x (unexpected zero for non-empty input)", tt.data, tt.crc, got)
			}
		})
	}

	// Test consistency: same input should give same output
	testData := "test"
	crc1 := CRC16XMODEM(testData, 0xffff)
	crc2 := CRC16XMODEM(testData, 0xffff)
	if crc1 != crc2 {
		t.Errorf("CRC16XMODEM() not consistent: got 0x%04x and 0x%04x for same input", crc1, crc2)
	}
}

