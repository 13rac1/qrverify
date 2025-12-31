package qrverify

import "image/color"

// Recovery specifies QR code error correction level.
type Recovery int

const (
	Low     Recovery = iota // 7% error correction
	Medium                  // 15% error correction (default)
	High                    // 25% error correction
	Highest                 // 30% error correction
)

// String returns the recovery level name.
func (r Recovery) String() string {
	switch r {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Highest:
		return "Highest"
	default:
		return "Recovery(unknown)"
	}
}

// EncodeOptions configures QR code generation.
// Zero values provide sensible defaults.
type EncodeOptions struct {
	// Recovery specifies error correction level.
	// Zero value uses Medium, then auto-escalates on verification failure.
	Recovery Recovery

	// Size is the image dimension in pixels.
	// Zero value uses 256.
	Size int

	// ForegroundColor sets the dark modules color.
	// Zero value uses black.
	ForegroundColor color.Color

	// BackgroundColor sets the light modules color.
	// Zero value uses white.
	BackgroundColor color.Color
}

// Result contains a verified QR code with metadata.
type Result struct {
	Image    []byte   // PNG image bytes
	Data     string   // Verified input data
	Version  int      // QR code version (1-40)
	Recovery Recovery // Final recovery level used
	Size     int      // Image dimensions in pixels
}
