package thaiqrgo

import (
	"fmt"
	"strconv"
	"strings"
)

// BOTBarcode represents a BOT Barcode structure.
type BOTBarcode struct {
	// BillerID is the biller identifier (Tax ID + Suffix)
	BillerID string

	// Ref1 is the reference number 1 / Customer number
	Ref1 string

	// Ref2 is the reference number 2 (optional)
	Ref2 *string

	// Amount is the transaction amount (optional)
	Amount *float64
}

// BOTBarcodeFromString parses a BOT Barcode data string.
//
// The payload must start with '|' and contain 4 fields separated by '\r'.
// Returns an error if the format is invalid.
func BOTBarcodeFromString(payload string) (*BOTBarcode, error) {
	if !strings.HasPrefix(payload, "|") {
		return nil, fmt.Errorf("invalid barcode format: must start with '|'")
	}

	data := strings.Split(payload[1:], "\r")
	if len(data) != 4 {
		return nil, fmt.Errorf("invalid barcode format: expected 4 fields, got %d", len(data))
	}

	billerID := data[0]
	ref1 := data[1]
	ref2 := data[2]
	amountStr := data[3]

	var ref2Ptr *string
	if len(ref2) > 0 {
		ref2Ptr = &ref2
	}

	var amountPtr *float64
	if amountStr != "0" {
		amountInt, err := strconv.ParseInt(amountStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount format: %w", err)
		}
		amount := float64(amountInt) / 100.0
		// Round to 2 decimal places
		amount = float64(int64(amount*100)) / 100.0
		amountPtr = &amount
	}

	return &BOTBarcode{
		BillerID: billerID,
		Ref1:     ref1,
		Ref2:     ref2Ptr,
		Amount:   amountPtr,
	}, nil
}

// String implements the fmt.Stringer interface.
//
// Returns the barcode in the standard BOT format: |billerID\rref1\rref2\ramount
func (b *BOTBarcode) String() string {
	var amountStr string
	if b.Amount != nil {
		amountInt := int64(*b.Amount * 100)
		amountStr = strconv.FormatInt(amountInt, 10)
	} else {
		amountStr = "0"
	}

	ref2 := ""
	if b.Ref2 != nil {
		ref2 = *b.Ref2
	}

	return fmt.Sprintf("|%s\r%s\r%s\r%s", b.BillerID, b.Ref1, ref2, amountStr)
}

// ToQRTag30 converts BOT Barcode to PromptPay QR Tag 30 (Bill Payment).
//
// This method works for some billers, depending on the destination bank.
// Returns a QR code payload string, or an error if conversion fails.
//
// Note: To avoid circular dependency between packages, users should call
// generate.BOTBarcodeToQR directly instead of this method.
// This method is kept for compatibility with the TypeScript version.
func (b *BOTBarcode) ToQRTag30() (string, error) {
	// Note: Direct implementation here would require importing generate package,
	// which creates a circular dependency. Use generate.BOTBarcodeToQR instead.
	// For compatibility, we provide a stub that returns an error with guidance.
	return "", fmt.Errorf("ToQRTag30 cannot be called directly due to package structure - use generate.BOTBarcodeToQR(b.BillerID, b.Ref1, b.Ref2, b.Amount) instead")
}

