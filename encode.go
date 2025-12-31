package qrverify

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

// Maximum data capacity in bytes for QR Version 40 at each recovery level.
const (
	MaxBytesLow     = 2953 // 7% error correction
	MaxBytesMedium  = 2331 // 15% error correction
	MaxBytesHigh    = 1663 // 25% error correction
	MaxBytesHighest = 1273 // 30% error correction
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

// recoveryLevel maps Recovery to go-qrcode RecoveryLevel.
func recoveryLevel(r Recovery) qrcode.RecoveryLevel {
	switch r {
	case Low:
		return qrcode.Low
	case High:
		return qrcode.High
	case Highest:
		return qrcode.Highest
	default:
		return qrcode.Medium
	}
}

// encodeAndVerify generates a QR code and verifies it decodes correctly.
// Returns PNG bytes, QR version, and error.
func encodeAndVerify(data string, recovery Recovery, size int) ([]byte, int, error) {
	// Create QR code
	qr, err := qrcode.New(data, recoveryLevel(recovery))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create QR code: %w", err)
	}

	// Generate PNG
	png, err := qr.PNG(size)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to generate PNG: %w", err)
	}

	// Verify by decoding
	if err := Verify(png, data); err != nil {
		return nil, 0, err
	}

	return png, qr.VersionNumber, nil
}

// Encode generates a verified QR code PNG image.
// Returns error if the generated code cannot be decoded back to data.
//
// If opts is nil or opts.Recovery is zero, uses Medium recovery (15%).
func Encode(data string, opts *EncodeOptions) ([]byte, error) {
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

	png, _, err := encodeAndVerify(data, recovery, size)
	return png, err
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

	png, version, err := encodeAndVerify(data, recovery, size)
	if err != nil {
		return nil, err
	}

	return &Result{
		Image:    png,
		Data:     data,
		Version:  version,
		Recovery: recovery,
		Size:     size,
	}, nil
}

// Quick generates a verified QR code with recommended defaults:
// - Medium recovery level (15% error correction)
// - 256x256 pixels
func Quick(data string) ([]byte, error) {
	return Encode(data, nil)
}
