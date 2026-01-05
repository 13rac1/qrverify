# qrverify

[![Go Reference](https://pkg.go.dev/badge/github.com/13rac1/qrverify.svg)](https://pkg.go.dev/github.com/13rac1/qrverify)
[![CI](https://github.com/13rac1/qrverify/actions/workflows/ci.yml/badge.svg)](https://github.com/13rac1/qrverify/actions/workflows/ci.yml)

Verified QR code generation for Go. Every generated QR code is guaranteed to decode back to the original data.

## Installation

```bash
go get github.com/13rac1/qrverify
```

## Quick Start

```go
png, err := qrverify.Encode("https://example.com", nil)
if err != nil {
    log.Fatal(err)
}
os.WriteFile("qr.png", png, 0644)
```

## Features

- **Verified output** - All generated QR codes are decoded and verified before returning
- **Size validation** - Validates data fits within QR capacity limits before encoding
- **Simple API** - `Encode()` with nil options for defaults, or custom options for control

## Implementation

The encoder and decoder libraries were selected based on performance and accuracy benchmarks from [qr-benchmarks](https://13rac1.github.io/qr-benchmarks/).

## API

| Function | Description |
|----------|-------------|
| `Encode(data, opts)` | Generate QR code (opts=nil for defaults: 256px, Medium recovery) |
| `EncodeToFile(data, filename, opts)` | Generate and write to file |
| `EncodeDetailed(data, opts)` | Generate with metadata result |
| `Decode(img)` | Decode QR code from image |
| `Verify(png, expected)` | Verify existing QR code |

## CLI

```bash
go install github.com/13rac1/qrverify/cmd/qrverify@latest

qrverify encode "https://example.com" -o qr.png
qrverify verify qr.png "https://example.com"
qrverify demo
```

## License

MIT License - see [LICENSE](LICENSE)
