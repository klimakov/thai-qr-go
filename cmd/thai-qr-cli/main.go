package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	thaiqrgo "thai-qr-go"
	"thai-qr-go/validate"
)

const version = "1.0.0"

func main() {
	var (
		payloadFlag = flag.String("payload", "", "QR code payload string to parse")
		formatFlag  = flag.String("format", "json", "Output format: json, text (default: json)")
		strictFlag  = flag.Bool("strict", false, "Validate CRC checksum (default: false)")
		showVersion = flag.Bool("version", false, "Show version and exit")
		helpFlag    = flag.Bool("help", false, "Show help message")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [PAYLOAD]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parse Thai QR code (PromptPay/EMVCo) payload and display structured data.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s \"00020101021129370016A000000677010111011300668012345675802TH530376463046197\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -format text -strict \"00020101021129370016...\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -payload \"|099999999999990\\r111222333444\\r\\r0\"\n", os.Args[0])
	}

	flag.Parse()

	if *showVersion {
		fmt.Printf("thai-qr-cli version %s\n", version)
		os.Exit(0)
	}

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Get payload from flag or positional argument
	payload := strings.TrimSpace(*payloadFlag)
	if payload == "" {
		args := flag.Args()
		if len(args) > 0 {
			payload = strings.TrimSpace(args[0])
		}
	}

	// Convert literal escape sequences to actual characters
	payload = convertEscapes(payload)

	if payload == "" {
		fmt.Fprintf(os.Stderr, "Error: QR code payload is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Try to parse as BOT Barcode first (if starts with |)
	if strings.HasPrefix(payload, "|") {
		if err := parseBarcode(payload, *formatFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing BOT Barcode: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Try to parse as EMVCo QR code
	if err := parseQR(payload, *formatFlag, *strictFlag); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing QR code: %v\n", err)
		os.Exit(1)
	}
}

func parseQR(payload, format string, strict bool) error {
	qr, err := thaiqrgo.Parse(payload, strict, true)
	if err != nil {
		return fmt.Errorf("failed to parse QR code: %w", err)
	}

	// Try to identify QR code type and extract structured data
	result := parseQRStructured(qr)

	switch strings.ToLower(format) {
	case "text":
		printQRText(qr, result)
	case "json":
		printQRJSON(qr, result)
	default:
		return fmt.Errorf("unknown format: %s (supported: json, text)", format)
	}

	return nil
}

func convertEscapes(s string) string {
	// Convert common escape sequences
	s = strings.ReplaceAll(s, "\\r", "\r")
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\t", "\t")
	return s
}

func parseBarcode(payload, format string) error {
	barcode, err := thaiqrgo.ParseBarcode(payload)
	if err != nil {
		return fmt.Errorf("failed to parse BOT Barcode: %w", err)
	}

	switch strings.ToLower(format) {
	case "text":
		printBarcodeText(barcode)
	case "json":
		printBarcodeJSON(barcode)
	default:
		return fmt.Errorf("unknown format: %s (supported: json, text)", format)
	}

	return nil
}

// QRCodeInfo contains extracted information from QR code
type QRCodeInfo struct {
	Type        string                            `json:"type,omitempty"`
	PhoneNumber string                            `json:"phone_number,omitempty"`
	NationalID  string                            `json:"national_id,omitempty"`
	TaxID       string                            `json:"tax_id,omitempty"`
	EWalletID   string                            `json:"ewallet_id,omitempty"`
	Amount      *float64                          `json:"amount,omitempty"`
	Currency    string                            `json:"currency,omitempty"`
	Country     string                            `json:"country,omitempty"`
	BillerID    string                            `json:"biller_id,omitempty"`
	Ref1        string                            `json:"ref1,omitempty"`
	Ref2        string                            `json:"ref2,omitempty"`
	Ref3        string                            `json:"ref3,omitempty"`
	Message     string                            `json:"message,omitempty"`
	SlipVerify  *validate.SlipVerifyData          `json:"slip_verify,omitempty"`
	TrueMoney   *validate.TrueMoneySlipVerifyData `json:"truemoney_slip_verify,omitempty"`
	Tags        []thaiqrgo.TLVTag                 `json:"tags,omitempty"`
	Valid       *bool                             `json:"valid,omitempty"`
	CRCValid    bool                              `json:"crc_valid"`
}

func parseQRStructured(qr *thaiqrgo.EMVCoQR) QRCodeInfo {
	info := QRCodeInfo{}

	// Get basic tags
	currency := qr.GetTagValue("53", "")
	country := qr.GetTagValue("58", "")
	amountStr := qr.GetTagValue("54", "")

	if currency != "" {
		info.Currency = currency
	}
	if country != "" {
		info.Country = country
	}
	if amountStr != "" {
		// Try to parse amount
		var amount float64
		if _, err := fmt.Sscanf(amountStr, "%f", &amount); err == nil {
			info.Amount = &amount
		}
	}

	// Try to identify QR type and extract data
	tag29 := qr.GetTag("29", "")
	tag30 := qr.GetTag("30", "")
	tag81 := qr.GetTagValue("81", "")

	// Check for Slip Verify
	slipVerify, err := validate.SlipVerify(qr.GetPayload())
	if err == nil {
		info.Type = "SlipVerify"
		info.SlipVerify = slipVerify
		valid := qr.Validate("91")
		info.CRCValid = valid
		return info
	}

	// Check for TrueMoney Slip Verify
	trueMoneySlip, err := validate.TrueMoneySlipVerify(qr.GetPayload())
	if err == nil {
		info.Type = "TrueMoneySlipVerify"
		info.TrueMoney = trueMoneySlip
		valid := qr.Validate("91")
		info.CRCValid = valid
		return info
	}

	// PromptPay AnyID (Tag 29)
	if tag29 != nil {
		info.Type = "PromptPayAnyID"

		// Check sub-tags
		phoneTag := qr.GetTag("29", "01")
		if phoneTag != nil {
			phone := phoneTag.Value
			// Remove leading 66 and zeros, restore leading 0
			if strings.HasPrefix(phone, "0066") && len(phone) == 13 {
				phone = "0" + phone[4:]
			}
			info.PhoneNumber = phone
		}

		nidTag := qr.GetTag("29", "02")
		if nidTag != nil {
			info.NationalID = nidTag.Value
			info.TaxID = nidTag.Value
		}

		ewalletTag := qr.GetTag("29", "03")
		if ewalletTag != nil {
			info.EWalletID = ewalletTag.Value
		}

		// Check for TrueMoney format (14000 prefix)
		if phoneTag != nil && strings.HasPrefix(phoneTag.Value, "14000") {
			info.Type = "TrueMoney"
			info.PhoneNumber = strings.TrimPrefix(phoneTag.Value, "14000")
		}
	}

	// PromptPay Bill Payment (Tag 30)
	if tag30 != nil {
		info.Type = "PromptPayBillPayment"
		billerID := qr.GetTagValue("30", "01")
		ref1 := qr.GetTagValue("30", "02")
		ref2 := qr.GetTagValue("30", "03")

		if billerID != "" {
			info.BillerID = billerID
		}
		if ref1 != "" {
			info.Ref1 = ref1
		}
		if ref2 != "" {
			info.Ref2 = ref2
		}

		// Check for Ref3 in Tag 62
		ref3Tag := qr.GetTag("62", "07")
		if ref3Tag != nil {
			info.Ref3 = ref3Tag.Value
		}
	}

	// Personal message (Tag 81) - decode from hex
	if tag81 != "" {
		info.Message = decodeTag81(tag81)
	}

	// Validate CRC if Tag 63 or 91 exists
	crcTag := qr.GetTag("63", "")
	if crcTag == nil {
		crcTag = qr.GetTag("91", "")
	}
	if crcTag != nil {
		crcTagID := "63"
		if crcTag.ID == "91" {
			crcTagID = "91"
		}
		valid := qr.Validate(crcTagID)
		info.CRCValid = valid
	}

	// Include all tags if format is json
	info.Tags = qr.GetTags()

	return info
}

func decodeTag81(hexStr string) string {
	var result strings.Builder
	for i := 0; i < len(hexStr); i += 4 {
		if i+4 > len(hexStr) {
			break
		}
		var codePoint int
		if _, err := fmt.Sscanf(hexStr[i:i+4], "%x", &codePoint); err != nil {
			// If parsing fails, skip this character
			continue
		}
		result.WriteRune(rune(codePoint))
	}
	return result.String()
}

func printQRJSON(qr *thaiqrgo.EMVCoQR, info QRCodeInfo) {
	output := map[string]interface{}{
		"payload": qr.GetPayload(),
		"info":    info,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func printQRText(qr *thaiqrgo.EMVCoQR, info QRCodeInfo) {
	fmt.Println("QR Code Information:")
	fmt.Println("===================")
	fmt.Printf("Payload: %s\n\n", qr.GetPayload())

	if info.Type != "" {
		fmt.Printf("Type: %s\n", info.Type)
	}

	if info.PhoneNumber != "" {
		fmt.Printf("Phone Number: %s\n", info.PhoneNumber)
	}
	if info.NationalID != "" {
		fmt.Printf("National ID: %s\n", info.NationalID)
	}
	if info.TaxID != "" && info.TaxID != info.NationalID {
		fmt.Printf("Tax ID: %s\n", info.TaxID)
	}
	if info.EWalletID != "" {
		fmt.Printf("E-Wallet ID: %s\n", info.EWalletID)
	}
	if info.BillerID != "" {
		fmt.Printf("Biller ID: %s\n", info.BillerID)
	}
	if info.Ref1 != "" {
		fmt.Printf("Reference 1: %s\n", info.Ref1)
	}
	if info.Ref2 != "" {
		fmt.Printf("Reference 2: %s\n", info.Ref2)
	}
	if info.Ref3 != "" {
		fmt.Printf("Reference 3: %s\n", info.Ref3)
	}
	if info.Amount != nil {
		fmt.Printf("Amount: %.2f\n", *info.Amount)
	}
	if info.Currency != "" {
		fmt.Printf("Currency: %s\n", info.Currency)
	}
	if info.Country != "" {
		fmt.Printf("Country: %s\n", info.Country)
	}
	if info.Message != "" {
		fmt.Printf("Message: %s\n", info.Message)
	}

	if info.SlipVerify != nil {
		fmt.Println("\nSlip Verify Data:")
		fmt.Printf("  Sending Bank: %s\n", info.SlipVerify.SendingBank)
		fmt.Printf("  Transaction Reference: %s\n", info.SlipVerify.TransRef)
	}

	if info.TrueMoney != nil {
		fmt.Println("\nTrueMoney Slip Verify Data:")
		fmt.Printf("  Event Type: %s\n", info.TrueMoney.EventType)
		fmt.Printf("  Transaction ID: %s\n", info.TrueMoney.TransactionID)
		fmt.Printf("  Date: %s\n", info.TrueMoney.Date)
	}

	fmt.Printf("\nCRC Valid: %v\n", info.CRCValid)

	fmt.Println("\nTLV Tags:")
	printTagsTree(qr.GetTags(), 0)
}

func printTagsTree(tags []thaiqrgo.TLVTag, indent int) {
	prefix := strings.Repeat("  ", indent)
	for _, tag := range tags {
		fmt.Printf("%sTag %s (length: %d): %s\n", prefix, tag.ID, tag.Length, tag.Value)
		if len(tag.SubTags) > 0 {
			printTagsTree(tag.SubTags, indent+1)
		}
	}
}

func printBarcodeJSON(barcode *thaiqrgo.BOTBarcode) {
	output := map[string]interface{}{
		"type":      "BOTBarcode",
		"biller_id": barcode.BillerID,
		"ref1":      barcode.Ref1,
	}

	if barcode.Ref2 != nil {
		output["ref2"] = *barcode.Ref2
	}

	if barcode.Amount != nil {
		output["amount"] = *barcode.Amount
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func printBarcodeText(barcode *thaiqrgo.BOTBarcode) {
	fmt.Println("BOT Barcode Information:")
	fmt.Println("========================")
	fmt.Printf("Biller ID: %s\n", barcode.BillerID)
	fmt.Printf("Reference 1: %s\n", barcode.Ref1)

	if barcode.Ref2 != nil {
		fmt.Printf("Reference 2: %s\n", *barcode.Ref2)
	}

	if barcode.Amount != nil {
		fmt.Printf("Amount: %.2f\n", *barcode.Amount)
	}
}
