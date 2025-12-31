package qrverify

import (
	"image/color"
	"testing"
)

func TestRecoveryString(t *testing.T) {
	tests := []struct {
		name     string
		recovery Recovery
		want     string
	}{
		{
			name:     "Low",
			recovery: Low,
			want:     "Low",
		},
		{
			name:     "Medium",
			recovery: Medium,
			want:     "Medium",
		},
		{
			name:     "High",
			recovery: High,
			want:     "High",
		},
		{
			name:     "Highest",
			recovery: Highest,
			want:     "Highest",
		},
		{
			name:     "Invalid negative",
			recovery: Recovery(-1),
			want:     "Recovery(unknown)",
		},
		{
			name:     "Invalid positive",
			recovery: Recovery(99),
			want:     "Recovery(unknown)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.recovery.String()
			if got != tt.want {
				t.Errorf("Recovery.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEncodeOptionsZeroValue(t *testing.T) {
	// Test that zero value is valid and usable
	var opts EncodeOptions

	// Zero value should have:
	// - Recovery: Low (0)
	// - Size: 0
	// - ForegroundColor: nil
	// - BackgroundColor: nil

	if opts.Recovery != Low {
		t.Errorf("zero value Recovery = %v, want %v", opts.Recovery, Low)
	}

	if opts.Size != 0 {
		t.Errorf("zero value Size = %d, want 0", opts.Size)
	}

	if opts.ForegroundColor != nil {
		t.Errorf("zero value ForegroundColor = %v, want nil", opts.ForegroundColor)
	}

	if opts.BackgroundColor != nil {
		t.Errorf("zero value BackgroundColor = %v, want nil", opts.BackgroundColor)
	}
}

func TestEncodeOptionsNonZeroValue(t *testing.T) {
	// Test that non-zero values are preserved
	opts := EncodeOptions{
		Recovery:        Highest,
		Size:            512,
		ForegroundColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		BackgroundColor: color.RGBA{R: 0, G: 255, B: 0, A: 255},
	}

	if opts.Recovery != Highest {
		t.Errorf("Recovery = %v, want %v", opts.Recovery, Highest)
	}

	if opts.Size != 512 {
		t.Errorf("Size = %d, want 512", opts.Size)
	}

	if opts.ForegroundColor == nil {
		t.Error("ForegroundColor should not be nil")
	}

	if opts.BackgroundColor == nil {
		t.Error("BackgroundColor should not be nil")
	}
}

func TestResultStruct(t *testing.T) {
	// Test that Result can be constructed and fields are accessible
	result := Result{
		Image:    []byte{1, 2, 3},
		Data:     "test data",
		Version:  5,
		Recovery: High,
		Size:     256,
	}

	if len(result.Image) != 3 {
		t.Errorf("Image length = %d, want 3", len(result.Image))
	}

	if result.Data != "test data" {
		t.Errorf("Data = %q, want %q", result.Data, "test data")
	}

	if result.Version != 5 {
		t.Errorf("Version = %d, want 5", result.Version)
	}

	if result.Recovery != High {
		t.Errorf("Recovery = %v, want %v", result.Recovery, High)
	}

	if result.Size != 256 {
		t.Errorf("Size = %d, want 256", result.Size)
	}
}
