package qrverify

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// decode reads a QR code from an image. Internal use only.
// Always uses TRY_HARDER hint for maximum accuracy.
func decode(img image.Image) (string, error) {
	// Convert image to BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("failed to create bitmap: %w", err)
	}

	// Set up hints with TRY_HARDER
	hints := make(map[gozxing.DecodeHintType]interface{})
	hints[gozxing.DecodeHintType_TRY_HARDER] = true

	// Decode
	reader := qrcode.NewQRCodeReader()
	result, err := reader.Decode(bmp, hints)
	if err != nil {
		return "", fmt.Errorf("failed to decode QR code: %w", err)
	}

	return result.GetText(), nil
}

// Verify checks that qrImage (PNG bytes) decodes to expectedData.
// Returns nil on success, VerificationError if mismatch, or error if decode fails.
func Verify(qrImage []byte, expectedData string) error {
	// Decode PNG
	img, err := png.Decode(bytes.NewReader(qrImage))
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %w", err)
	}

	// Decode QR
	decoded, err := decode(img)
	if err != nil {
		return err
	}

	// Strict byte-for-byte comparison
	if decoded != expectedData {
		return &VerificationError{
			Original: expectedData,
			Decoded:  decoded,
			Recovery: Medium, // Default, actual recovery unknown for external images
		}
	}

	return nil
}
