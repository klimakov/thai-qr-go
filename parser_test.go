package thaiqrgo

import "testing"

func TestParse_InvalidString(t *testing.T) {
	_, err := Parse("AAAA0000", false, true)
	if err == nil {
		t.Error("Parse() should return error for invalid string")
	}
}

func TestParse_TLV(t *testing.T) {
	qr, err := Parse("000411110104222202043333", false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(qr.GetTags()) != 3 {
		t.Errorf("Parse() tag count = %v, want 3", len(qr.GetTags()))
	}
}

func TestParse_GetTag(t *testing.T) {
	qr, err := Parse("000411110104222202043333", false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	tag := qr.GetTag("01", "")
	if tag == nil || tag.Value != "2222" {
		t.Errorf("Parse() GetTag(01) = %v, want Value=2222", tag)
	}
}

func TestParse_StrictInvalidChecksum(t *testing.T) {
	payload := "00020101021229370016A0000006770101110113006680111111153037645802TH540520.156304FFFF"
	_, err := Parse(payload, true, true)
	if err == nil {
		t.Error("Parse() with strict=true should return error for invalid checksum")
	}
}

func TestParse_StrictValidChecksum(t *testing.T) {
	payload := "00020101021229370016A0000006770101110113006680111111153037645802TH540520.15630442BE"
	qr, err := Parse(payload, true, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	value := qr.GetTagValue("29", "01")
	if value != "0066801111111" {
		t.Errorf("Parse() GetTagValue(29, 01) = %v, want 0066801111111", value)
	}
}

// Test compatibility with promptpay-qr library
func TestParse_PromptPayQR_PhoneNumber(t *testing.T) {
	// From promptpay-qr test: '0801234567' without amount
	payload := "00020101021129370016A000000677010111011300668012345675802TH530376463046197"
	qr, err := Parse(payload, false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check Tag 29 exists
	tag29 := qr.GetTag("29", "")
	if tag29 == nil {
		t.Error("Parse() should have Tag 29")
	}

	// Check Tag 29.01 (phone number)
	phoneTag := qr.GetTag("29", "01")
	if phoneTag == nil {
		t.Fatal("Parse() should have Tag 29.01")
		return
	}
	if phoneTag.Value != "0066801234567" {
		t.Errorf("Parse() Tag 29.01 = %v, want 0066801234567", phoneTag.Value)
	}
}

func TestParse_PromptPayQR_PhoneNumberWithAmount(t *testing.T) {
	// From promptpay-qr test: '000-000-0000' with amount 4.22
	payload := "00020101021229370016A000000677010111011300660000000005802TH530376454044.226304E469"
	qr, err := Parse(payload, false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check amount tag
	amount := qr.GetTagValue("54", "")
	if amount != "4.22" {
		t.Errorf("Parse() Tag 54 = %v, want 4.22", amount)
	}
}

func TestParse_PromptPayQR_NationalID(t *testing.T) {
	// From promptpay-qr test: '1111111111111'
	payload := "00020101021129370016A000000677010111021311111111111115802TH530376463047B5A"
	qr, err := Parse(payload, false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check Tag 29.02 (National ID)
	nidTag := qr.GetTag("29", "02")
	if nidTag == nil {
		t.Fatal("Parse() should have Tag 29.02")
		return
	}
	if nidTag.Value != "1111111111111" {
		t.Errorf("Parse() Tag 29.02 = %v, want 1111111111111", nidTag.Value)
	}
}

func TestParse_PromptPayQR_TaxID(t *testing.T) {
	// From promptpay-qr test: '0123456789012'
	payload := "00020101021129370016A000000677010111021301234567890125802TH530376463040CBD"
	qr, err := Parse(payload, false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check Tag 29.02 (Tax ID)
	taxTag := qr.GetTag("29", "02")
	if taxTag == nil {
		t.Fatal("Parse() should have Tag 29.02")
		return
	}
	if taxTag.Value != "0123456789012" {
		t.Errorf("Parse() Tag 29.02 = %v, want 0123456789012", taxTag.Value)
	}
}

func TestParse_PromptPayQR_EWalletID(t *testing.T) {
	// From promptpay-qr test: '012345678901234'
	payload := "00020101021129390016A00000067701011103150123456789012345802TH530376463049781"
	qr, err := Parse(payload, false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check Tag 29.03 (E-Wallet ID)
	ewalletTag := qr.GetTag("29", "03")
	if ewalletTag == nil {
		t.Fatal("Parse() should have Tag 29.03")
		return
	}
	if ewalletTag.Value != "012345678901234" {
		t.Errorf("Parse() Tag 29.03 = %v, want 012345678901234", ewalletTag.Value)
	}
}

func TestValidateChecksum(t *testing.T) {
	payload := "00020101021229370016A0000006770101110113006680111111153037645802TH540520.15630442BE"
	qr, err := Parse(payload, false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if !qr.Validate("63") {
		t.Error("Validate() should return true for valid checksum")
	}

	// Test with invalid CRC tag ID
	if qr.Validate("99") {
		t.Error("Validate() with non-existent CRC tag should recalculate and return false")
	}
}

func TestParse_EdgeCases(t *testing.T) {
	// Test with subTags=false
	qr, err := Parse("000411110104222202043333", false, false)
	if err != nil {
		t.Fatalf("Parse() with subTags=false error = %v", err)
	}
	if len(qr.GetTags()) != 3 {
		t.Errorf("Parse() with subTags=false tag count = %v, want 3", len(qr.GetTags()))
	}

	// Test with empty payload
	_, err = Parse("", false, true)
	if err == nil {
		t.Error("Parse() should return error for empty payload")
	}

	// Test with payload too short for strict mode
	_, err = Parse("12", true, true)
	if err == nil {
		t.Error("Parse() with strict=true should return error for payload too short")
	}

	// Test with invalid TLV format
	_, err = Parse("0001", false, true)
	if err == nil {
		t.Error("Parse() should return error for invalid TLV format")
	}
}

func TestParse_SubTags(t *testing.T) {
	// Test parsing with nested subTags
	payload := "00020101021229370016A0000006770101110113006680111111153037645802TH540520.15630442BE"
	qr, err := Parse(payload, false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check that Tag 29 has subTags
	tag29 := qr.GetTag("29", "")
	if tag29 == nil {
		t.Fatal("Parse() Tag 29 should exist")
		return
	}
	if len(tag29.SubTags) == 0 {
		t.Error("Parse() Tag 29 should have subTags")
	}

	// Verify subTag exists
	subTag := qr.GetTag("29", "01")
	if subTag == nil {
		t.Error("Parse() Tag 29.01 should exist")
	}
}

func TestEMVCoQR_GetTagValue(t *testing.T) {
	qr, err := Parse("000411110104222202043333", false, true)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Test GetTagValue for existing tag
	value := qr.GetTagValue("01", "")
	if value != "2222" {
		t.Errorf("GetTagValue() = %v, want 2222", value)
	}

	// Test GetTagValue for non-existent tag
	value = qr.GetTagValue("99", "")
	if value != "" {
		t.Errorf("GetTagValue() for non-existent tag should return empty string, got %v", value)
	}

	// Test GetPayload
	payload := qr.GetPayload()
	if payload != "000411110104222202043333" {
		t.Errorf("GetPayload() = %v, want 000411110104222202043333", payload)
	}
}
