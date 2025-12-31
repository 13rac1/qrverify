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
png, err := qrverify.Quick("https://example.com")
if err != nil {
    log.Fatal(err)
}
os.WriteFile("qr.png", png, 0644)
```

## Features

- **Verified output** - All generated QR codes are decoded and verified before returning
- **Auto-retry** - Automatically escalates error correction level on verification failure
- **Simple API** - `Quick()` for defaults, `Encode()` for options

## API

| Function | Description |
|----------|-------------|
| `Quick(data)` | Generate with defaults (256px, Medium recovery) |
| `Encode(data, opts)` | Generate with custom options |
| `EncodeToFile(data, filename, opts)` | Generate and write to file |
| `EncodeDetailed(data, opts)` | Generate with metadata result |
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
