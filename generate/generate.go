package generate

import (
	"strconv"
	"strings"

	"thai-qr-go"
	"thai-qr-go/internal"
)

// ProxyType constants for PromptPay AnyID.
const (
	ProxyTypeMSISDN   = "01" // Mobile number
	ProxyTypeNATID    = "02" // National ID or Tax ID
	ProxyTypeEWALLETID = "03" // E-Wallet ID
	ProxyTypeBANKACC  = "04" // Bank Account (Reserved)
)

// AnyIDConfig configures a PromptPay AnyID QR code.
type AnyIDConfig struct {
	// Type is the proxy type (MSISDN, NATID, EWALLETID, BANKACC)
	Type string

	// Target is the recipient number
	Target string

	// Amount is the transaction amount (optional)
	Amount *float64
}

// AnyID generates a PromptPay AnyID (Tag 29) QR code payload.
func AnyID(config AnyIDConfig) (string, error) {
	target := config.Target

	if config.Type == "MSISDN" {
		// Remove leading 0 and prepend country code 66, pad to 13 digits
		target = strings.TrimPrefix(target, "0")
		target = "66" + target
		// Pad to 13 digits from the right
		for len(target) < 13 {
			target = "0" + target
		}
		if len(target) > 13 {
			target = target[len(target)-13:]
		}
	}

	proxyTypeValue := config.Type
	switch config.Type {
	case "MSISDN":
		proxyTypeValue = ProxyTypeMSISDN
	case "NATID":
		proxyTypeValue = ProxyTypeNATID
	case "EWALLETID":
		proxyTypeValue = ProxyTypeEWALLETID
	case "BANKACC":
		proxyTypeValue = ProxyTypeBANKACC
	default:
		return "", &InvalidConfigError{Field: "Type", Value: config.Type}
	}

	tag29 := thaiqrgo.Encode([]thaiqrgo.TLVTag{
		thaiqrgo.Tag("00", "A000000677010111"),
		thaiqrgo.Tag(proxyTypeValue, target),
	})

	var payload []thaiqrgo.TLVTag
	payload = append(payload, thaiqrgo.Tag("00", "01"))
	if config.Amount != nil {
		payload = append(payload, thaiqrgo.Tag("01", "12"))
	} else {
		payload = append(payload, thaiqrgo.Tag("01", "11"))
	}
	payload = append(payload, thaiqrgo.Tag("29", tag29))
	payload = append(payload, thaiqrgo.Tag("53", "764"))
	payload = append(payload, thaiqrgo.Tag("58", "TH"))

	if config.Amount != nil {
		amountStr := strconv.FormatFloat(*config.Amount, 'f', 2, 64)
		payload = append(payload, thaiqrgo.Tag("54", amountStr))
	}

	return thaiqrgo.WithCRCTag(thaiqrgo.Encode(payload), "63", true), nil
}

// BillPaymentConfig configures a PromptPay Bill Payment QR code.
type BillPaymentConfig struct {
	// BillerID is the biller identifier (National ID or Tax ID + Suffix)
	BillerID string

	// Amount is the transaction amount (optional)
	Amount *float64

	// Ref1 is reference 1
	Ref1 string

	// Ref2 is reference 2 (optional)
	Ref2 *string

	// Ref3 is reference 3 (optional, undocumented)
	Ref3 *string
}

// BillPayment generates a PromptPay Bill Payment (Tag 30) QR code payload.
func BillPayment(config BillPaymentConfig) (string, error) {
	tag30 := []thaiqrgo.TLVTag{
		thaiqrgo.Tag("00", "A000000677010112"),
		thaiqrgo.Tag("01", config.BillerID),
		thaiqrgo.Tag("02", config.Ref1),
	}

	if config.Ref2 != nil {
		tag30 = append(tag30, thaiqrgo.Tag("03", *config.Ref2))
	}

	var payload []thaiqrgo.TLVTag
	payload = append(payload, thaiqrgo.Tag("00", "01"))
	if config.Amount != nil {
		payload = append(payload, thaiqrgo.Tag("01", "12"))
	} else {
		payload = append(payload, thaiqrgo.Tag("01", "11"))
	}
	payload = append(payload, thaiqrgo.Tag("30", thaiqrgo.Encode(tag30)))
	payload = append(payload, thaiqrgo.Tag("53", "764"))
	payload = append(payload, thaiqrgo.Tag("58", "TH"))

	if config.Amount != nil {
		amountStr := strconv.FormatFloat(*config.Amount, 'f', 2, 64)
		payload = append(payload, thaiqrgo.Tag("54", amountStr))
	}

	if config.Ref3 != nil {
		tag62 := thaiqrgo.Encode([]thaiqrgo.TLVTag{
			thaiqrgo.Tag("07", *config.Ref3),
		})
		payload = append(payload, thaiqrgo.Tag("62", tag62))
	}

	return thaiqrgo.WithCRCTag(thaiqrgo.Encode(payload), "63", true), nil
}

// TrueMoneyConfig configures a TrueMoney QR code.
type TrueMoneyConfig struct {
	// MobileNo is the mobile number
	MobileNo string

	// Amount is the transaction amount (optional)
	Amount *float64

	// Message is a personal message for Tag 81 (optional)
	Message *string
}

// TrueMoney generates a QR code for TrueMoney Wallet.
//
// This QR code can also be scanned with other apps, just like a regular e-Wallet PromptPay QR,
// but the Personal Message (Tag 81) will be ignored.
func TrueMoney(config TrueMoneyConfig) (string, error) {
	tag29 := thaiqrgo.Encode([]thaiqrgo.TLVTag{
		thaiqrgo.Tag("00", "A000000677010111"),
		thaiqrgo.Tag("03", "14000"+config.MobileNo),
	})

	var payload []thaiqrgo.TLVTag
	payload = append(payload, thaiqrgo.Tag("00", "01"))
	if config.Amount != nil {
		payload = append(payload, thaiqrgo.Tag("01", "12"))
	} else {
		payload = append(payload, thaiqrgo.Tag("01", "11"))
	}
	payload = append(payload, thaiqrgo.Tag("29", tag29))
	payload = append(payload, thaiqrgo.Tag("53", "764"))
	payload = append(payload, thaiqrgo.Tag("58", "TH"))

	if config.Amount != nil {
		amountStr := strconv.FormatFloat(*config.Amount, 'f', 2, 64)
		payload = append(payload, thaiqrgo.Tag("54", amountStr))
	}

	if config.Message != nil {
		encodedMsg := internal.EncodeTag81(*config.Message)
		payload = append(payload, thaiqrgo.Tag("81", encodedMsg))
	}

	return thaiqrgo.WithCRCTag(thaiqrgo.Encode(payload), "63", true), nil
}

// SlipVerifyConfig configures a Slip Verify QR code.
type SlipVerifyConfig struct {
	// SendingBank is the bank code
	SendingBank string

	// TransRef is the transaction reference
	TransRef string
}

// SlipVerify generates a Slip Verify QR code.
//
// This is also called "Mini-QR" that is embedded in slips used for verifying transactions.
func SlipVerify(config SlipVerifyConfig) (string, error) {
	tag00 := thaiqrgo.Encode([]thaiqrgo.TLVTag{
		thaiqrgo.Tag("00", "000001"),
		thaiqrgo.Tag("01", config.SendingBank),
		thaiqrgo.Tag("02", config.TransRef),
	})

	payload := []thaiqrgo.TLVTag{
		thaiqrgo.Tag("00", tag00),
		thaiqrgo.Tag("51", "TH"),
	}

	return thaiqrgo.WithCRCTag(thaiqrgo.Encode(payload), "91", true), nil
}

// TrueMoneySlipVerifyConfig configures a TrueMoney Slip Verify QR code.
type TrueMoneySlipVerifyConfig struct {
	// EventType is the event type (e.g., "P2P")
	EventType string

	// TransactionID is the transaction ID
	TransactionID string

	// Date is the date in DDMMYYYY format
	Date string
}

// TrueMoneySlipVerify generates a TrueMoney Slip Verify QR code.
//
// Same as a regular Slip Verify QR but with some differences:
//   - Tag 00 and 01 are set to '01'
//   - Tag 51 does not exist
//   - Additional tags that are TrueMoney-specific
//   - CRC checksum is case-sensitive (lowercase)
func TrueMoneySlipVerify(config TrueMoneySlipVerifyConfig) (string, error) {
	tag00 := thaiqrgo.Encode([]thaiqrgo.TLVTag{
		thaiqrgo.Tag("00", "01"),
		thaiqrgo.Tag("01", "01"),
		thaiqrgo.Tag("02", config.EventType),
		thaiqrgo.Tag("03", config.TransactionID),
		thaiqrgo.Tag("04", config.Date),
	})

	payload := []thaiqrgo.TLVTag{
		thaiqrgo.Tag("00", tag00),
	}

	return thaiqrgo.WithCRCTag(thaiqrgo.Encode(payload), "91", false), nil
}

// BOTBarcodeConfig configures a BOT Barcode.
type BOTBarcodeConfig struct {
	// BillerID is the biller identifier (Tax ID + Suffix)
	BillerID string

	// Ref1 is reference number 1 / Customer number
	Ref1 string

	// Ref2 is reference number 2 (optional)
	Ref2 *string

	// Amount is the transaction amount (optional)
	Amount *float64
}

// BOTBarcode generates a BOT Barcode string.
func BOTBarcode(config BOTBarcodeConfig) string {
	barcode := &thaiqrgo.BOTBarcode{
		BillerID: config.BillerID,
		Ref1:     config.Ref1,
		Ref2:     config.Ref2,
		Amount:   config.Amount,
	}
	return barcode.String()
}

// BOTBarcodeToQR converts a BOT Barcode to a PromptPay QR Tag 30 (Bill Payment) payload.
//
// This function works for some billers, depending on the destination bank.
// It takes the same parameters as BOTBarcode and returns a QR code payload.
func BOTBarcodeToQR(billerID, ref1 string, ref2 *string, amount *float64) (string, error) {
	config := BillPaymentConfig{
		BillerID: billerID,
		Ref1:     ref1,
		Ref2:     ref2,
		Amount:   amount,
	}
	return BillPayment(config)
}

// InvalidConfigError represents an error in configuration.
type InvalidConfigError struct {
	Field string
	Value string
}

func (e *InvalidConfigError) Error() string {
	return "invalid config: " + e.Field + " = " + e.Value
}


