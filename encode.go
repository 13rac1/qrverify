package qrverify

import (
	"bytes"
	"fmt"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// Maximum data capacity in bytes for QR Version 40 (largest standard QR code)
// at each error correction level. Based on binary/byte mode encoding.
// See Recovery type for error correction percentages.
const (
	MaxBytesLow     = 2953
	MaxBytesMedium  = 2331
	MaxBytesHigh    = 1663
	MaxBytesHighest = 1273
)

// maxBytes returns the maximum data capacity for a recovery level.
func maxBytes(r Recovery) int {
	switch r {
	case Low:
		return MaxBytesLow
	case High:
		return MaxBytesHigh
	case Highest:
		return MaxBytesHighest
	default:
		return MaxBytesMedium
	}
}

// recoveryLevel maps Recovery to barcode qr.ErrorCorrectionLevel.
func recoveryLevel(r Recovery) qr.ErrorCorrectionLevel {
	switch r {
	case Low:
		return qr.L
	case High:
		return qr.Q
	case Highest:
		return qr.H
	default:
		return qr.M
	}
}

// encodeAndVerify generates a QR code and verifies it decodes correctly.
// Returns PNG bytes and error.
func encodeAndVerify(data string, recovery Recovery, size int) ([]byte, error) {
	level := recoveryLevel(recovery)

	// Create QR code with 8-bit greyscale
	bc, err := qr.EncodeWithColor(data, level, qr.Auto, barcode.ColorScheme8)
	if err != nil {
		return nil, fmt.Errorf("failed to create QR code: %w", err)
	}

	// Scale to size
	bc, err = barcode.Scale(bc, size, size)
	if err != nil {
		return nil, fmt.Errorf("failed to scale QR code: %w", err)
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, bc); err != nil {
		return nil, fmt.Errorf("failed to generate PNG: %w", err)
	}

	// Verify by decoding
	if err := Verify(buf.Bytes(), data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Encode generates a verified QR code PNG image.
// Returns error if the generated code cannot be decoded back to data.
//
// If opts is nil or opts.Recovery is zero, uses Medium recovery (15%).
func Encode(data string, opts *EncodeOptions) ([]byte, error) {
	result, err := EncodeDetailed(data, opts)
	if err != nil {
		return nil, err
	}
	return result.Image, nil
}

// EncodeToFile generates a verified QR code and writes it to filename.
func EncodeToFile(data string, filename string, opts *EncodeOptions) error {
	png, err := Encode(data, opts)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, png, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// EncodeDetailed returns the verified QR code with metadata.
func EncodeDetailed(data string, opts *EncodeOptions) (*Result, error) {
	size := 256
	recovery := Medium
	if opts != nil {
		if opts.Size > 0 {
			size = opts.Size
		}
		if opts.Recovery != 0 {
			recovery = opts.Recovery
		}
	}

	if len(data) > maxBytes(recovery) {
		return nil, fmt.Errorf("data too large: %d bytes exceeds %d byte limit for %v recovery",
			len(data), maxBytes(recovery), recovery)
	}

	png, err := encodeAndVerify(data, recovery, size)
	if err != nil {
		return nil, err
	}

	return &Result{
		Image:    png,
		Data:     data,
		Recovery: recovery,
		Size:     size,
	}, nil
}
