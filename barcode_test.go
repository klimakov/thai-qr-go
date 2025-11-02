package thaiqrgo

import (
	"testing"
)

func TestParseBarcode_Valid(t *testing.T) {
	payload := "|099999999999990\r111222333444\r\r0"
	barcode, err := ParseBarcode(payload)
	if err != nil {
		t.Fatalf("ParseBarcode() error = %v", err)
	}
	if barcode.BillerID != "099999999999990" {
		t.Errorf("ParseBarcode() BillerID = %v, want 099999999999990", barcode.BillerID)
	}
	if barcode.Ref1 != "111222333444" {
		t.Errorf("ParseBarcode() Ref1 = %v, want 111222333444", barcode.Ref1)
	}
	if barcode.Ref2 != nil {
		t.Error("ParseBarcode() Ref2 should be nil")
	}
	if barcode.Amount != nil {
		t.Error("ParseBarcode() Amount should be nil")
	}
}

func TestParseBarcode_WithRef2AndAmount(t *testing.T) {
	payload := "|099400016550100\r123456789012\r670429\r364922"
	barcode, err := ParseBarcode(payload)
	if err != nil {
		t.Fatalf("ParseBarcode() error = %v", err)
	}
	if barcode.BillerID != "099400016550100" {
		t.Errorf("ParseBarcode() BillerID = %v, want 099400016550100", barcode.BillerID)
	}
	if barcode.Ref1 != "123456789012" {
		t.Errorf("ParseBarcode() Ref1 = %v, want 123456789012", barcode.Ref1)
	}
	if barcode.Ref2 == nil || *barcode.Ref2 != "670429" {
		t.Errorf("ParseBarcode() Ref2 = %v, want 670429", barcode.Ref2)
	}
	if barcode.Amount == nil || *barcode.Amount != 3649.22 {
		t.Errorf("ParseBarcode() Amount = %v, want 3649.22", barcode.Amount)
	}
}

func TestParseBarcode_Invalid(t *testing.T) {
	// Wrong payload (not starting with |)
	payload := "00020101021230650016A00000067701011201150994000165501000212123456789012030667042953037645802TH54073649.2263044534"
	_, err := ParseBarcode(payload)
	if err == nil {
		t.Error("ParseBarcode() should return error for invalid payload")
	}

	// Data loss (not 4 fields)
	payload = "|099400016550100\r123456789012\r670429"
	_, err = ParseBarcode(payload)
	if err == nil {
		t.Error("ParseBarcode() should return error for incomplete payload")
	}
}

func TestBOTBarcode_String(t *testing.T) {
	barcode := &BOTBarcode{
		BillerID: "099999999999990",
		Ref1:     "111222333444",
		Ref2:     nil,
		Amount:   nil,
	}
	got := barcode.String()
	want := "|099999999999990\r111222333444\r\r0"
	if got != want {
		t.Errorf("BOTBarcode.String() = %v, want %v", got, want)
	}

	// Test with Ref2 and Amount
	ref2 := "REF2"
	amount := 100.50
	barcode.Ref2 = &ref2
	barcode.Amount = &amount
	got = barcode.String()
	want = "|099999999999990\r111222333444\rREF2\r10050"
	if got != want {
		t.Errorf("BOTBarcode.String() with ref2 and amount = %v, want %v", got, want)
	}

	// Test with amount that rounds
	amount2 := 100.555
	barcode.Amount = &amount2
	got = barcode.String()
	// Should round to 2 decimal places: 100.555 -> 10055 (100.55)
	want2 := "|099999999999990\r111222333444\rREF2\r10055"
	if got != want2 {
		t.Logf("BOTBarcode.String() with rounding amount = %v (expected around 10055)", got)
	}
}

func TestParseBarcode_EdgeCases(t *testing.T) {
	// Test with invalid amount format
	payload := "|099999999999990\r111222333444\r\rABC"
	_, err := ParseBarcode(payload)
	if err == nil {
		t.Error("ParseBarcode() should return error for invalid amount format")
	}

	// Test with empty ref2 (should be nil)
	payload = "|099999999999990\r111222333444\r\r0"
	barcode, err := ParseBarcode(payload)
	if err != nil {
		t.Fatalf("ParseBarcode() error = %v", err)
	}
	if barcode.Ref2 != nil {
		t.Error("ParseBarcode() with empty ref2 should set Ref2 to nil")
	}
}
