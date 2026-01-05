// Package qrverify provides verified QR code generation.
//
// It wraps github.com/boombuler/barcode (encoding) and
// github.com/makiuchi-d/gozxing (decoding) to guarantee that every
// generated QR code can be successfully decoded back to the original input data.
//
// # Quick Start
//
// Generate a verified QR code with sensible defaults:
//
//	png, err := qrverify.Encode("https://example.com", nil)
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
// # Decoding
//
// Decode QR codes from images:
//
//	img, _ := png.Decode(file)
//	data, err := qrverify.Decode(img)
//
// # Verification
//
// All Encode functions verify the generated QR code decodes correctly.
// Use Verify to check existing QR code images:
//
//	err := qrverify.Verify(pngBytes, "expected data")
package qrverify
