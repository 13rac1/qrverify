package qrverify

import (
	"fmt"
	"os"

	"github.com/skip2/go-qrcode"
)

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
// If opts is nil or opts.Recovery is zero, starts with Medium and
// auto-escalates (Medium -> High -> Highest) on verification failure.
// If opts.Recovery is explicitly set, uses that level without retry.
func Encode(data string, opts *EncodeOptions) ([]byte, error) {
	size := 256
	if opts != nil && opts.Size > 0 {
		size = opts.Size
	}

	// If explicit recovery requested, use it without retry
	if opts != nil && opts.Recovery != 0 {
		png, _, err := encodeAndVerify(data, opts.Recovery, size)
		return png, err
	}

	// Auto-escalate: Medium -> High -> Highest
	var lastErr error
	for _, r := range []Recovery{Medium, High, Highest} {
		png, _, err := encodeAndVerify(data, r, size)
		if err == nil {
			return png, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("failed to encode verifiable QR code: %w", lastErr)
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

	// If explicit recovery requested, use it without retry
	if opts != nil && opts.Recovery != 0 {
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

	// Auto-escalate: Medium -> High -> Highest
	var lastErr error
	for _, r := range []Recovery{Medium, High, Highest} {
		png, version, err := encodeAndVerify(data, r, size)
		if err == nil {
			return &Result{
				Image:    png,
				Data:     data,
				Version:  version,
				Recovery: r,
				Size:     size,
			}, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("failed to encode verifiable QR code: %w", lastErr)
}

// Quick generates a verified QR code with recommended defaults:
// - Medium recovery level (15% error correction)
// - 256x256 pixels
// - Auto-retry enabled
func Quick(data string) ([]byte, error) {
	return Encode(data, nil)
}
