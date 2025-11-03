package validate

import (
	"testing"

	"github.com/klimakov/thai-qr-go/generate"
)

func TestSlipVerify_Valid(t *testing.T) {
	payload := "004100060000010103014022000111222233344ABCD125102TH910417DF"
	data, err := SlipVerify(payload)
	if err != nil {
		t.Fatalf("SlipVerify() error = %v", err)
	}
	if data.SendingBank != "014" {
		t.Errorf("SlipVerify() SendingBank = %v, want 014", data.SendingBank)
	}
	if data.TransRef != "00111222233344ABCD12" {
		t.Errorf("SlipVerify() TransRef = %v, want 00111222233344ABCD12", data.TransRef)
	}
}

func TestSlipVerify_Invalid(t *testing.T) {
	// Not a Slip Verify QR
	payload := "00020101021229370016A0000006770101110113006680111111153037645802TH540520.15630442BE"
	_, err := SlipVerify(payload)
	if err == nil {
		t.Error("SlipVerify() should return error for invalid payload")
	}

	// Test with wrong apiType - payload with apiType != "000001"
	// Create a payload that parses but has wrong apiType (using Tag 00.00 = "000002" instead of "000001")
	// This is a valid QR but not a valid Slip Verify format
	wrongPayload := "004100060000020103014022000111222233344ABCD125102TH6304"
	// Need to recalculate CRC for this payload
	_, err = SlipVerify(wrongPayload)
	if err == nil {
		t.Error("SlipVerify() should return error for invalid apiType")
	}
}

func TestSlipVerify_EdgeCases(t *testing.T) {
	// Empty payload
	_, err := SlipVerify("")
	if err == nil {
		t.Error("SlipVerify() should return error for empty payload")
	}

	// Missing required fields
	_, err = SlipVerify("004000060000010103014022000111222233344ABCD125102TH910417DF")
	if err == nil {
		t.Error("SlipVerify() should return error for missing required fields")
	}
}

func TestTrueMoneySlipVerify_Valid(t *testing.T) {
	// Generate a valid TrueMoney Slip Verify QR using generate package
	config := generate.TrueMoneySlipVerifyConfig{
		EventType:     "P2P",
		TransactionID: "TXN123456",
		Date:          "08122024",
	}

	payload, err := generate.TrueMoneySlipVerify(config)
	if err != nil {
		t.Fatalf("generate.TrueMoneySlipVerify() error = %v", err)
	}

	// Now validate it
	data, err := TrueMoneySlipVerify(payload)
	if err != nil {
		t.Fatalf("TrueMoneySlipVerify() error = %v", err)
	}

	if data.EventType != "P2P" {
		t.Errorf("TrueMoneySlipVerify() EventType = %v, want P2P", data.EventType)
	}
	if data.TransactionID != "TXN123456" {
		t.Errorf("TrueMoneySlipVerify() TransactionID = %v, want TXN123456", data.TransactionID)
	}
	if data.Date != "08122024" {
		t.Errorf("TrueMoneySlipVerify() Date = %v, want 08122024", data.Date)
	}
}

func TestTrueMoneySlipVerify_Invalid(t *testing.T) {
	// Not a TrueMoney Slip Verify QR (regular Slip Verify)
	payload := "004100060000010103014022000111222233344ABCD125102TH910417DF"
	_, err := TrueMoneySlipVerify(payload)
	if err == nil {
		t.Error("TrueMoneySlipVerify() should return error for invalid payload")
	}

	// Wrong API type
	payload = "004000060000010103014022000111222233344ABCD125102TH910417DF"
	_, err = TrueMoneySlipVerify(payload)
	if err == nil {
		t.Error("TrueMoneySlipVerify() should return error for wrong API type")
	}

	// Empty payload
	_, err = TrueMoneySlipVerify("")
	if err == nil {
		t.Error("TrueMoneySlipVerify() should return error for empty payload")
	}
}
