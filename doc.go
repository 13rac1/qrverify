// Package qrverify provides verified QR code generation.
//
// It wraps go-qrcode (encoding) and gozxing (decoding) to guarantee
// that every generated QR code can be successfully decoded back to
// the original input data.
//
// # Quick Start
//
// Generate a verified QR code with sensible defaults:
//
//	png, err := qrverify.Quick("https://example.com")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	os.WriteFile("qr.png", png, 0644)
//
// # Custom Options
//
// Use EncodeOptions for fine-grained control:
//
//	opts := &qrverify.EncodeOptions{
//	    Recovery: qrverify.High,  // 25% error correction
//	    Size:     512,            // 512x512 pixels
//	}
//	png, err := qrverify.Encode("data", opts)
//
// # Verification
//
// All Encode functions verify the generated QR code decodes correctly.
// Use Verify to check existing QR code images:
//
//	err := qrverify.Verify(pngBytes, "expected data")
package qrverify
