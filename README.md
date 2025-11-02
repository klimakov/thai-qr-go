# Thai QR Go

All-in-one Go library for PromptPay & EMVCo QR Codes

## Features

- **Parse** — PromptPay & EMVCo QR Code data strings into object
- **Generate** — QR Code data from pre-made templates (PromptPay AnyID, PromptPay Bill Payment, TrueMoney, etc.)
- **Manipulate** — any values from parsed QR Code data (transfer amount, account number) and encodes back into QR Code data
- **Validate** — checksum and data structure for known QR Code formats (Slip Verify API Mini QR)

## Installation

### Library

```bash
go get thai-qr-go
```

### CLI Tool

Build the CLI tool:

```bash
git clone <repository-url>
cd thai-qr-go
go build -o bin/thai-qr-cli ./cmd/thai-qr-cli
```

Or install globally:

```bash
go install thai-qr-go/cmd/thai-qr-cli@latest
```

## Usage

### CLI Tool

Parse QR code from command line:

```bash
# Parse QR code (JSON output)
thai-qr-cli "00020101021129370016A000000677010111011300668012345675802TH530376463046197"

# Parse with text output
thai-qr-cli -format text "00020101021129370016..."

# Parse with strict CRC validation
thai-qr-cli -strict "00020101021129370016..."

# Parse BOT Barcode (use \r for carriage return)
thai-qr-cli '|099999999999990\r111222333444\r\r0'

# Show help
thai-qr-cli --help
```

**Output formats:**
- `json` (default) - Structured JSON output with all parsed data
- `text` - Human-readable text format with extracted information

### Library Usage

### Parsing data and get value from tag

```go
package main

import (
    "fmt"
    "thai-qr-go"
)

func main() {
    // Example data
    ppqr, err := thaiqrgo.Parse("000201010211...", false, true)
    if err != nil {
        panic(err)
    }
    
    // Get Value of Tag ID '00'
    value := ppqr.GetTagValue("00")
    fmt.Println(value) // Returns '01'
}
```

### Build QR data and append CRC tag

```go
package main

import (
    "fmt"
    "thai-qr-go"
)

func main() {
    // Example data
    data := []thaiqrgo.TLVTag{
        thaiqrgo.Tag("00", "01"),
        thaiqrgo.Tag("01", "11"),
        // ...
    }
    
    // Set CRC Tag ID '63'
    payload := thaiqrgo.WithCRCTag(thaiqrgo.Encode(data), "63")
    fmt.Println(payload) // Returns '000201010211...'
}
```

### Generate PromptPay Bill Payment QR

```go
package main

import (
    "fmt"
    "thai-qr-go/generate"
)

func main() {
    payload, err := generate.BillPayment(generate.BillPaymentConfig{
        BillerID: "1xxxxxxxxxxxx",
        Amount:   300.0,
        Ref1:     "INV12345",
    })
    if err != nil {
        panic(err)
    }
    
    // TODO: Create QR Code from payload
    fmt.Println(payload)
}
```

### Validate & extract data from Slip Verify QR

```go
package main

import (
    "fmt"
    "thai-qr-go/validate"
)

func main() {
    data, err := validate.SlipVerify("00550006000001...")
    if err != nil {
        fmt.Printf("Invalid Payload: %v\n", err)
        return
    }
    
    fmt.Printf("Sending Bank: %s, Trans Ref: %s\n", data.SendingBank, data.TransRef)
}
```

## References

- [EMV QR Code](https://www.emvco.com/emv-technologies/qrcodes/)
- [Thai QR Payment Standard](https://www.bot.or.th/content/dam/bot/fipcs/documents/FPG/2562/ThaiPDF/25620084.pdf)
- [Slip Verify API Mini QR Data](https://developer.scb/assets/documents/documentation/qr-payment/extracting-data-from-mini-qr.pdf)
- [BOT Barcode Standard](https://www.bot.or.th/content/dam/bot/documents/th/our-roles/payment-systems/about-payment-systems/Std_Barcode.pdf)

## License

This project is MIT licensed.

Copyright (c) 2024-2025 Klimakov Oleg

See [LICENSE](LICENSE) file for details.

