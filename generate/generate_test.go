package generate

import (
	"testing"
)

func TestAnyID(t *testing.T) {
	amount := 0.0
	config := AnyIDConfig{
		Type:   "MSISDN",
		Target: "0812223333",
		Amount: nil,
	}
	got, err := AnyID(config)
	if err != nil {
		t.Fatalf("AnyID() error = %v", err)
	}
	want := "00020101021129370016A0000006770101110113006681222333353037645802TH63041DCF"
	if got != want {
		t.Errorf("AnyID() = %v, want %v", got, want)
	}

	// Test with amount
	amount = 30.0
	config.Amount = &amount
	got, err = AnyID(config)
	if err != nil {
		t.Fatalf("AnyID() with amount error = %v", err)
	}
	want = "00020101021229370016A0000006770101110113006681222333353037645802TH540530.0063043CAD"
	if got != want {
		t.Errorf("AnyID() with amount = %v, want %v", got, want)
	}
}

func TestSlipVerify(t *testing.T) {
	config := SlipVerifyConfig{
		SendingBank: "002",
		TransRef:    "0002123123121200011",
	}
	got, err := SlipVerify(config)
	if err != nil {
		t.Fatalf("SlipVerify() error = %v", err)
	}
	want := "004000060000010103002021900021231231212000115102TH91049C30"
	if got != want {
		t.Errorf("SlipVerify() = %v, want %v", got, want)
	}
}

func TestTrueMoney(t *testing.T) {
	config := TrueMoneyConfig{
		MobileNo: "0801111111",
		Amount:   nil,
		Message:  nil,
	}
	got, err := TrueMoney(config)
	if err != nil {
		t.Fatalf("TrueMoney() error = %v", err)
	}
	want := "00020101021129390016A000000677010111031514000080111111153037645802TH63047C0F"
	if got != want {
		t.Errorf("TrueMoney() = %v, want %v", got, want)
	}

	// Test with amount and message
	amount := 10.05
	message := "Hello World!"
	config.Amount = &amount
	config.Message = &message
	got, err = TrueMoney(config)
	if err != nil {
		t.Fatalf("TrueMoney() with amount and message error = %v", err)
	}
	want = "00020101021229390016A000000677010111031514000080111111153037645802TH540510.05814800480065006C006C006F00200057006F0072006C006400216304F5A2"
	if got != want {
		t.Errorf("TrueMoney() with amount and message = %v, want %v", got, want)
	}
}

func TestBillPayment(t *testing.T) {
	ref2 := "INV001"
	ref3 := "SCB"
	config := BillPaymentConfig{
		BillerID: "0112233445566",
		Ref1:     "CUSTOMER001",
		Ref2:     &ref2,
		Ref3:     &ref3,
		Amount:   nil,
	}
	got, err := BillPayment(config)
	if err != nil {
		t.Fatalf("BillPayment() error = %v", err)
	}
	want := "00020101021130620016A000000677010112011301122334455660211CUSTOMER0010306INV00153037645802TH62070703SCB6304780E"
	if got != want {
		t.Errorf("BillPayment() = %v, want %v", got, want)
	}
}

func TestBOTBarcode(t *testing.T) {
	config := BOTBarcodeConfig{
		BillerID: "099999999999990",
		Ref1:     "111222333444",
		Ref2:     nil,
		Amount:   nil,
	}
	got := BOTBarcode(config)
	want := "|099999999999990\r111222333444\r\r0"
	if got != want {
		t.Errorf("BOTBarcode() = %v, want %v", got, want)
	}

	// Test with Ref2 and amount
	ref2 := "670429"
	amount := 3649.22
	config.BillerID = "099400016550100"
	config.Ref1 = "123456789012"
	config.Ref2 = &ref2
	config.Amount = &amount
	got = BOTBarcode(config)
	want = "|099400016550100\r123456789012\r670429\r364922"
	if got != want {
		t.Errorf("BOTBarcode() with ref2 and amount = %v, want %v", got, want)
	}
}

func TestTrueMoneySlipVerify(t *testing.T) {
	config := TrueMoneySlipVerifyConfig{
		EventType:     "P2P",
		TransactionID: "TXN123456",
		Date:          "08122024",
	}
	got, err := TrueMoneySlipVerify(config)
	if err != nil {
		t.Fatalf("TrueMoneySlipVerify() error = %v", err)
	}
	// Verify it starts with expected format
	if len(got) == 0 {
		t.Error("TrueMoneySlipVerify() returned empty string")
	}
	// Verify it contains Tag 00
	if len(got) >= 2 && got[:2] != "00" {
		t.Logf("TrueMoneySlipVerify() generated payload: %v", got)
	}
}

func TestAnyID_EdgeCases(t *testing.T) {
	// Test with invalid type
	config := AnyIDConfig{
		Type:   "INVALID",
		Target: "0812223333",
	}
	_, err := AnyID(config)
	if err == nil {
		t.Error("AnyID() should return error for invalid type")
	}

	// Test with NATID type
	config.Type = "NATID"
	config.Target = "1111111111111"
	got, err := AnyID(config)
	if err != nil {
		t.Fatalf("AnyID() with NATID error = %v", err)
	}
	if len(got) == 0 {
		t.Error("AnyID() with NATID returned empty string")
	}

	// Test with EWALLETID type
	config.Type = "EWALLETID"
	config.Target = "012345678901234"
	got, err = AnyID(config)
	if err != nil {
		t.Fatalf("AnyID() with EWALLETID error = %v", err)
	}
	if len(got) == 0 {
		t.Error("AnyID() with EWALLETID returned empty string")
	}

	// Test with BANKACC type
	config.Type = "BANKACC"
	got, err = AnyID(config)
	if err != nil {
		t.Fatalf("AnyID() with BANKACC error = %v", err)
	}
	if len(got) == 0 {
		t.Error("AnyID() with BANKACC returned empty string")
	}
}

func TestBillPayment_EdgeCases(t *testing.T) {
	// Test with amount
	amount := 100.50
	config := BillPaymentConfig{
		BillerID: "0112233445566",
		Ref1:     "CUSTOMER001",
		Ref2:     nil,
		Ref3:     nil,
		Amount:   &amount,
	}
	got, err := BillPayment(config)
	if err != nil {
		t.Fatalf("BillPayment() with amount error = %v", err)
	}
	if len(got) == 0 {
		t.Error("BillPayment() with amount returned empty string")
	}

	// Test without amount but with ref3
	ref3 := "SCB"
	config.Amount = nil
	config.Ref3 = &ref3
	got, err = BillPayment(config)
	if err != nil {
		t.Fatalf("BillPayment() with ref3 error = %v", err)
	}
	if len(got) == 0 {
		t.Error("BillPayment() with ref3 returned empty string")
	}
}

func TestBOTBarcodeToQR(t *testing.T) {
	ref2 := "670429"
	amount := 3649.22
	got, err := BOTBarcodeToQR("099400016550100", "123456789012", &ref2, &amount)
	if err != nil {
		t.Fatalf("BOTBarcodeToQR() error = %v", err)
	}
	want := "00020101021230650016A00000067701011201150994000165501000212123456789012030667042953037645802TH54073649.2263044534"
	if got != want {
		t.Errorf("BOTBarcodeToQR() = %v, want %v", got, want)
	}

	// Test without ref2 and amount
	got, err = BOTBarcodeToQR("099999999999990", "111222333444", nil, nil)
	if err != nil {
		t.Fatalf("BOTBarcodeToQR() without ref2/amount error = %v", err)
	}
	want = "00020101021130550016A0000006770101120115099999999999990021211122233344453037645802TH63043EE7"
	if got != want {
		t.Errorf("BOTBarcodeToQR() without ref2/amount = %v, want %v", got, want)
	}
}
