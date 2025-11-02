package validate

import (
	"errors"
	"thai-qr-go"
)

// SlipVerifyData contains extracted data from a Slip Verify QR code.
type SlipVerifyData struct {
	// SendingBank is the bank code
	SendingBank string

	// TransRef is the transaction reference
	TransRef string
}

// SlipVerify validates and extracts data from a Slip Verify QR code.
//
// This function is used with Bank Open API to extract bank code and transaction reference.
// Returns an error if the payload is invalid or doesn't match the Slip Verify format.
func SlipVerify(payload string) (*SlipVerifyData, error) {
	ppqr, err := thaiqrgo.Parse(payload, true, true)
	if err != nil {
		return nil, err
	}

	apiType := ppqr.GetTagValue("00", "00")
	sendingBank := ppqr.GetTagValue("00", "01")
	transRef := ppqr.GetTagValue("00", "02")

	if apiType != "000001" || sendingBank == "" || transRef == "" {
		return nil, errors.New("invalid Slip Verify format: missing required fields")
	}

	return &SlipVerifyData{
		SendingBank: sendingBank,
		TransRef:    transRef,
	}, nil
}

// TrueMoneySlipVerifyData contains extracted data from a TrueMoney Slip Verify QR code.
type TrueMoneySlipVerifyData struct {
	// EventType is the event type (e.g., "P2P")
	EventType string

	// TransactionID is the transaction ID
	TransactionID string

	// Date is the date in DDMMYYYY format
	Date string
}

// TrueMoneySlipVerify validates and extracts data from a TrueMoney Slip Verify QR code.
//
// Returns an error if the payload is invalid or doesn't match the TrueMoney Slip Verify format.
func TrueMoneySlipVerify(payload string) (*TrueMoneySlipVerifyData, error) {
	ppqr, err := thaiqrgo.Parse(payload, true, true)
	if err != nil {
		return nil, err
	}

	apiType00 := ppqr.GetTagValue("00", "00")
	apiType01 := ppqr.GetTagValue("00", "01")
	eventType := ppqr.GetTagValue("00", "02")
	transactionID := ppqr.GetTagValue("00", "03")
	date := ppqr.GetTagValue("00", "04")

	if apiType00 != "01" || apiType01 != "01" {
		return nil, errors.New("invalid TrueMoney Slip Verify format: incorrect API type")
	}

	if eventType == "" || transactionID == "" || date == "" {
		return nil, errors.New("invalid TrueMoney Slip Verify format: missing required fields")
	}

	return &TrueMoneySlipVerifyData{
		EventType:     eventType,
		TransactionID: transactionID,
		Date:          date,
	}, nil
}
