package qrverify

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skip2/go-qrcode"
)

func TestEncode(t *testing.T) {
	data := "https://example.com"
	png, err := Encode(data, nil)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	if len(png) == 0 {
		t.Fatal("Expected non-empty PNG data")
	}

	// Verify the generated QR code
	if err := Verify(png, data); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}

func TestEncodeWithOptions(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		opts     *EncodeOptions
		wantSize int
	}{
		{
			name:     "custom size",
			data:     "test data",
			opts:     &EncodeOptions{Size: 512},
			wantSize: 512,
		},
		{
			name:     "low recovery",
			data:     "test data",
			opts:     &EncodeOptions{Recovery: Low, Size: 256},
			wantSize: 256,
		},
		{
			name:     "high recovery",
			data:     "test data",
			opts:     &EncodeOptions{Recovery: High, Size: 256},
			wantSize: 256,
		},
		{
			name:     "highest recovery",
			data:     "test data",
			opts:     &EncodeOptions{Recovery: Highest, Size: 256},
			wantSize: 256,
		},
		{
			name:     "default size with recovery",
			data:     "test data",
			opts:     &EncodeOptions{Recovery: Medium},
			wantSize: 256,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			png, err := Encode(tt.data, tt.opts)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			if len(png) == 0 {
				t.Fatal("Expected non-empty PNG data")
			}

			// Verify the generated QR code
			if err := Verify(png, tt.data); err != nil {
				t.Errorf("Verification failed: %v", err)
			}
		})
	}
}

func TestEncodeDefaultRecovery(t *testing.T) {
	// Use simple data that should work with default Medium recovery
	data := "test"
	png, err := Encode(data, nil)
	if err != nil {
		t.Fatalf("Encode with default recovery failed: %v", err)
	}

	if len(png) == 0 {
		t.Fatal("Expected non-empty PNG data")
	}

	// Verify it works
	if err := Verify(png, data); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}

func TestEncodeToFile(t *testing.T) {
	data := "https://example.com/test"
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.png")

	err := EncodeToFile(data, filename, nil)
	if err != nil {
		t.Fatalf("EncodeToFile failed: %v", err)
	}

	// Read the file back
	pngData, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if len(pngData) == 0 {
		t.Fatal("Expected non-empty file")
	}

	// Verify the QR code from file
	if err := Verify(pngData, data); err != nil {
		t.Errorf("Verification of file contents failed: %v", err)
	}
}

func TestEncodeToFileWithOptions(t *testing.T) {
	data := "custom options test"
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "custom.png")

	opts := &EncodeOptions{
		Size:     512,
		Recovery: High,
	}

	err := EncodeToFile(data, filename, opts)
	if err != nil {
		t.Fatalf("EncodeToFile with options failed: %v", err)
	}

	// Verify file exists and is valid
	pngData, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if err := Verify(pngData, data); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}

func TestEncodeToFileInvalidPath(t *testing.T) {
	data := "test"
	// Use invalid path (directory that doesn't exist)
	err := EncodeToFile(data, "/nonexistent/directory/file.png", nil)
	if err == nil {
		t.Fatal("Expected error for invalid path, got nil")
	}
}

func TestEncodeDetailed(t *testing.T) {
	data := "detailed test data"
	opts := &EncodeOptions{
		Size:     512,
		Recovery: High,
	}

	result, err := EncodeDetailed(data, opts)
	if err != nil {
		t.Fatalf("EncodeDetailed failed: %v", err)
	}

	// Check all Result fields
	if len(result.Image) == 0 {
		t.Error("Expected non-empty Image")
	}

	if result.Data != data {
		t.Errorf("Expected Data=%q, got %q", data, result.Data)
	}

	if result.Version < 1 || result.Version > 40 {
		t.Errorf("Expected Version between 1-40, got %d", result.Version)
	}

	if result.Recovery != High {
		t.Errorf("Expected Recovery=High, got %v", result.Recovery)
	}

	if result.Size != 512 {
		t.Errorf("Expected Size=512, got %d", result.Size)
	}

	// Verify the image
	if err := Verify(result.Image, data); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}

func TestEncodeDetailedWithDefaults(t *testing.T) {
	data := "default detailed test"

	result, err := EncodeDetailed(data, nil)
	if err != nil {
		t.Fatalf("EncodeDetailed with defaults failed: %v", err)
	}

	if result.Size != 256 {
		t.Errorf("Expected default Size=256, got %d", result.Size)
	}

	// Recovery should be default Medium
	if result.Recovery != Medium {
		t.Errorf("Expected Recovery to be Medium, got %v", result.Recovery)
	}
}

func TestQuick(t *testing.T) {
	data := "quick test"
	png, err := Quick(data)
	if err != nil {
		t.Fatalf("Quick failed: %v", err)
	}

	if len(png) == 0 {
		t.Fatal("Expected non-empty PNG data")
	}

	// Verify the generated QR code
	if err := Verify(png, data); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}

func TestDataTypes(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{
			name: "numeric",
			data: "1234567890",
		},
		{
			name: "alphanumeric",
			data: "ABCDEFGHIJKLMNOPQRSTUVWXYZ 0123456789",
		},
		{
			name: "lowercase alphanumeric",
			data: "abcdefghijklmnopqrstuvwxyz",
		},
		{
			name: "mixed case",
			data: "Hello World 123",
		},
		{
			name: "UTF-8",
			data: "Hello 世界 Мир",
		},
		{
			name: "emoji",
			data: "Hello 👋 🌍",
		},
		{
			name: "URL",
			data: "https://example.com/path?query=value&foo=bar",
		},
		{
			name: "JSON-like",
			data: `{"key":"value","number":123}`,
		},
		{
			name: "special characters",
			data: "!@#$%^&*()_+-=[]{}|;':\",./<>?",
		},
		{
			name: "newlines",
			data: "line1\nline2\nline3",
		},
		{
			name: "tabs",
			data: "col1\tcol2\tcol3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			png, err := Encode(tt.data, nil)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			if len(png) == 0 {
				t.Fatal("Expected non-empty PNG data")
			}

			// Verify the generated QR code
			if err := Verify(png, tt.data); err != nil {
				t.Errorf("Verification failed: %v", err)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		wantError bool
	}{
		{
			name:      "empty string",
			data:      "",
			wantError: true, // QR codes require at least 1 character
		},
		{
			name:      "single character",
			data:      "X",
			wantError: false,
		},
		{
			name:      "whitespace only",
			data:      "   ",
			wantError: false,
		},
		{
			name:      "single space",
			data:      " ",
			wantError: false,
		},
		{
			name:      "tab character",
			data:      "\t",
			wantError: false,
		},
		{
			name:      "newline character",
			data:      "\n",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			png, err := Encode(tt.data, nil)
			if tt.wantError {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			if len(png) == 0 {
				t.Fatal("Expected non-empty PNG data")
			}

			// Verify the generated QR code
			if err := Verify(png, tt.data); err != nil {
				t.Errorf("Verification failed: %v", err)
			}
		})
	}
}

func TestEncodeTooLarge(t *testing.T) {
	// QR code version 40 with Low recovery can hold ~2953 bytes
	// Create data larger than maximum capacity
	largeData := strings.Repeat("A", 3000)

	_, err := Encode(largeData, &EncodeOptions{Recovery: Low})
	if err == nil {
		t.Fatal("Expected error for data too large, got nil")
	}

	// Also test with default recovery (Medium)
	_, err = Encode(largeData, nil)
	if err == nil {
		t.Fatal("Expected error for data too large with default recovery, got nil")
	}
}

func TestRecoveryLevel(t *testing.T) {
	tests := []struct {
		name     string
		recovery Recovery
		want     qrcode.RecoveryLevel
	}{
		{
			name:     "Low",
			recovery: Low,
			want:     qrcode.Low,
		},
		{
			name:     "Medium",
			recovery: Medium,
			want:     qrcode.Medium,
		},
		{
			name:     "High",
			recovery: High,
			want:     qrcode.High,
		},
		{
			name:     "Highest",
			recovery: Highest,
			want:     qrcode.Highest,
		},
		{
			name:     "invalid defaults to Medium",
			recovery: Recovery(99),
			want:     qrcode.Medium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := recoveryLevel(tt.recovery)
			if got != tt.want {
				t.Errorf("recoveryLevel(%v) = %v, want %v", tt.recovery, got, tt.want)
			}
		})
	}
}

func TestEncodeAndVerify(t *testing.T) {
	data := "test data for encodeAndVerify"
	size := 256

	png, version, err := encodeAndVerify(data, Medium, size)
	if err != nil {
		t.Fatalf("encodeAndVerify failed: %v", err)
	}

	if len(png) == 0 {
		t.Fatal("Expected non-empty PNG data")
	}

	if version < 1 || version > 40 {
		t.Errorf("Expected version between 1-40, got %d", version)
	}

	// Verify the QR code
	if err := Verify(png, data); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}

func TestEncodeAndVerifyInvalidData(t *testing.T) {
	// Empty string should fail
	_, _, err := encodeAndVerify("", Medium, 256)
	if err == nil {
		t.Fatal("Expected error for empty data, got nil")
	}
}

func TestEncodeWithAllRecoveryLevels(t *testing.T) {
	data := "test all recovery levels"

	recoveryLevels := []Recovery{Low, Medium, High, Highest}
	for _, r := range recoveryLevels {
		t.Run(r.String(), func(t *testing.T) {
			opts := &EncodeOptions{Recovery: r}
			png, err := Encode(data, opts)
			if err != nil {
				t.Fatalf("Encode with %v failed: %v", r, err)
			}

			if len(png) == 0 {
				t.Fatal("Expected non-empty PNG data")
			}

			// Verify the generated QR code
			if err := Verify(png, data); err != nil {
				t.Errorf("Verification failed for %v: %v", r, err)
			}
		})
	}
}

func TestEncodeDetailedDefaultRecovery(t *testing.T) {
	data := "default recovery detailed"

	// Don't specify recovery to use default Medium
	result, err := EncodeDetailed(data, &EncodeOptions{Size: 512})
	if err != nil {
		t.Fatalf("EncodeDetailed with default recovery failed: %v", err)
	}

	if result.Size != 512 {
		t.Errorf("Expected Size=512, got %d", result.Size)
	}

	// Should have used default Medium recovery
	if result.Recovery != Medium {
		t.Errorf("Expected Recovery to be Medium, got %v", result.Recovery)
	}

	if err := Verify(result.Image, data); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}

func TestEncodeVerificationFailure(t *testing.T) {
	// This test ensures that if verification fails, we get an error
	// We can't easily force a verification failure with real QR codes,
	// but we can test the error path by using data that's too large
	largeData := strings.Repeat("X", 4000)

	_, err := Encode(largeData, &EncodeOptions{Recovery: Low})
	if err == nil {
		t.Fatal("Expected error for oversized data, got nil")
	}
}

func TestMultipleEncodings(t *testing.T) {
	// Test that multiple encodings of the same data work correctly
	data := "consistency test"

	png1, err := Encode(data, nil)
	if err != nil {
		t.Fatalf("First encode failed: %v", err)
	}

	png2, err := Encode(data, nil)
	if err != nil {
		t.Fatalf("Second encode failed: %v", err)
	}

	// Both should verify correctly
	if err := Verify(png1, data); err != nil {
		t.Errorf("First verification failed: %v", err)
	}

	if err := Verify(png2, data); err != nil {
		t.Errorf("Second verification failed: %v", err)
	}

	// Note: PNG bytes may not be identical due to timestamps or other metadata,
	// but both should decode to the same data
}

func TestEncodeSizeVariations(t *testing.T) {
	data := "size test"

	sizes := []int{64, 128, 256, 512, 1024}
	for _, size := range sizes {
		t.Run(string(rune(size)), func(t *testing.T) {
			opts := &EncodeOptions{Size: size}
			png, err := Encode(data, opts)
			if err != nil {
				t.Fatalf("Encode with size %d failed: %v", size, err)
			}

			if len(png) == 0 {
				t.Fatal("Expected non-empty PNG data")
			}

			if err := Verify(png, data); err != nil {
				t.Errorf("Verification failed for size %d: %v", size, err)
			}
		})
	}
}

func TestEncodeCompareWithVerify(t *testing.T) {
	data := "comparison test"

	// Encode with our function
	png, err := Encode(data, nil)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Verify should succeed
	if err := Verify(png, data); err != nil {
		t.Errorf("Verify failed: %v", err)
	}

	// Verify with wrong data should fail
	if err := Verify(png, "wrong data"); err == nil {
		t.Error("Expected verification to fail with wrong data")
	}
}

func TestEncodeNilOptions(t *testing.T) {
	data := "nil options test"

	// Encode with nil options should use defaults
	png, err := Encode(data, nil)
	if err != nil {
		t.Fatalf("Encode with nil options failed: %v", err)
	}

	if len(png) == 0 {
		t.Fatal("Expected non-empty PNG data")
	}

	if err := Verify(png, data); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}

func TestEncodeEmptyOptions(t *testing.T) {
	data := "empty options test"

	// Encode with empty options struct should use defaults
	png, err := Encode(data, &EncodeOptions{})
	if err != nil {
		t.Fatalf("Encode with empty options failed: %v", err)
	}

	if len(png) == 0 {
		t.Fatal("Expected non-empty PNG data")
	}

	if err := Verify(png, data); err != nil {
		t.Errorf("Verification failed: %v", err)
	}
}

func BenchmarkEncode(b *testing.B) {
	data := "benchmark test data"
	for i := 0; i < b.N; i++ {
		_, err := Encode(data, nil)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
	}
}

func BenchmarkEncodeWithOptions(b *testing.B) {
	data := "benchmark test data"
	opts := &EncodeOptions{
		Size:     512,
		Recovery: High,
	}
	for i := 0; i < b.N; i++ {
		_, err := Encode(data, opts)
		if err != nil {
			b.Fatalf("Encode failed: %v", err)
		}
	}
}

func BenchmarkQuick(b *testing.B) {
	data := "benchmark quick test"
	for i := 0; i < b.N; i++ {
		_, err := Quick(data)
		if err != nil {
			b.Fatalf("Quick failed: %v", err)
		}
	}
}

// TestLongURLs tests encoding of realistic long URLs
func TestLongURLs(t *testing.T) {
	urls := []string{
		"https://example.com/very/long/path/to/resource?param1=value1&param2=value2&param3=value3",
		"https://api.example.com/v1/users/12345/posts/67890/comments?page=1&limit=100&sort=desc&include=author,likes",
	}

	for _, url := range urls {
		t.Run(url[:50], func(t *testing.T) {
			png, err := Encode(url, nil)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			if err := Verify(png, url); err != nil {
				t.Errorf("Verification failed: %v", err)
			}
		})
	}
}

// TestBinaryData tests that invalid UTF-8 sequences are handled correctly
func TestBinaryData(t *testing.T) {
	// QR codes encode data as UTF-8 strings, so invalid UTF-8 sequences
	// won't round-trip correctly. This is expected behavior.
	// Test that we get an error rather than silent corruption.
	binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}
	data := string(binaryData)

	// This should fail because the data contains invalid UTF-8
	_, err := Encode(data, nil)
	if err == nil {
		t.Fatal("Expected error for invalid UTF-8 data, got nil")
	}

	// The error should be a verification error, indicating that
	// the data didn't round-trip correctly
	if _, ok := err.(*VerificationError); !ok {
		// It might be wrapped, check the error message
		if err != nil && !bytes.Contains([]byte(err.Error()), []byte("verification failed")) {
			t.Errorf("Expected VerificationError, got: %v", err)
		}
	}
}

// TestVeryShortData ensures single-character inputs work
func TestVeryShortData(t *testing.T) {
	shortInputs := []string{"A", "1", "!", " ", "\n"}

	for _, input := range shortInputs {
		t.Run(input, func(t *testing.T) {
			png, err := Encode(input, nil)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			if err := Verify(png, input); err != nil {
				t.Errorf("Verification failed: %v", err)
			}
		})
	}
}

// TestResultImageNotModified ensures the Result.Image is the actual PNG
func TestResultImageNotModified(t *testing.T) {
	data := "result image test"

	result, err := EncodeDetailed(data, nil)
	if err != nil {
		t.Fatalf("EncodeDetailed failed: %v", err)
	}

	// The Image field should be valid PNG data
	if !bytes.HasPrefix(result.Image, []byte{0x89, 0x50, 0x4E, 0x47}) {
		t.Error("Result.Image does not appear to be valid PNG (missing PNG header)")
	}

	// Should verify correctly
	if err := Verify(result.Image, data); err != nil {
		t.Errorf("Verification of Result.Image failed: %v", err)
	}
}

// TestMaxBytes tests all recovery levels in maxBytes function
func TestMaxBytes(t *testing.T) {
	tests := []struct {
		name     string
		recovery Recovery
		want     int
	}{
		{
			name:     "Low recovery",
			recovery: Low,
			want:     MaxBytesLow,
		},
		{
			name:     "Medium recovery",
			recovery: Medium,
			want:     MaxBytesMedium,
		},
		{
			name:     "High recovery",
			recovery: High,
			want:     MaxBytesHigh,
		},
		{
			name:     "Highest recovery",
			recovery: Highest,
			want:     MaxBytesHighest,
		},
		{
			name:     "Invalid recovery defaults to Medium",
			recovery: Recovery(99),
			want:     MaxBytesMedium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := maxBytes(tt.recovery)
			if got != tt.want {
				t.Errorf("maxBytes(%v) = %d, want %d", tt.recovery, got, tt.want)
			}
		})
	}
}

// TestEncodeDetailedDataTooLarge tests EncodeDetailed with oversized data
func TestEncodeDetailedDataTooLarge(t *testing.T) {
	tests := []struct {
		name     string
		recovery Recovery
		dataSize int
	}{
		{
			name:     "exceeds Low recovery limit",
			recovery: Low,
			dataSize: MaxBytesLow + 100,
		},
		{
			name:     "exceeds Medium recovery limit",
			recovery: Medium,
			dataSize: MaxBytesMedium + 100,
		},
		{
			name:     "exceeds High recovery limit",
			recovery: High,
			dataSize: MaxBytesHigh + 100,
		},
		{
			name:     "exceeds Highest recovery limit",
			recovery: Highest,
			dataSize: MaxBytesHighest + 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			largeData := strings.Repeat("X", tt.dataSize)
			opts := &EncodeOptions{Recovery: tt.recovery}

			_, err := EncodeDetailed(largeData, opts)
			if err == nil {
				t.Fatal("Expected error for oversized data, got nil")
			}

			// Verify error message mentions data too large
			if !strings.Contains(err.Error(), "data too large") {
				t.Errorf("Expected 'data too large' error, got: %v", err)
			}
		})
	}
}

// TestEncodeEdgeCaseSizes tests encoding with edge case sizes
func TestEncodeEdgeCaseSizes(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "negative size handled by library",
			size: -1,
		},
		{
			name: "zero size defaults to 256",
			size: 0,
		},
		{
			name: "very small size",
			size: 1,
		},
		{
			name: "very large size",
			size: 2048,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := "test"
			opts := &EncodeOptions{Size: tt.size}
			png, err := Encode(data, opts)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if len(png) == 0 {
				t.Error("Expected non-empty PNG data")
			}
		})
	}
}

// TestEncodeToFileErrors tests error conditions in EncodeToFile
func TestEncodeToFileErrors(t *testing.T) {
	t.Run("encode error propagates", func(t *testing.T) {
		// Use data that's too large to trigger encode error
		largeData := strings.Repeat("X", MaxBytesLow+100)
		tempDir := t.TempDir()
		filename := filepath.Join(tempDir, "test.png")

		err := EncodeToFile(largeData, filename, &EncodeOptions{Recovery: Low})
		if err == nil {
			t.Fatal("Expected error for oversized data, got nil")
		}

		// File should not have been created
		if _, statErr := os.Stat(filename); statErr == nil {
			t.Error("File should not exist after encode error")
		}
	})
}
