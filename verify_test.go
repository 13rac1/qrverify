package qrverify

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/skip2/go-qrcode"
)

// mockBrokenImage is a custom image type with completely empty bounds
// to trigger bitmap creation failure
type mockBrokenImage struct{}

func (m *mockBrokenImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (m *mockBrokenImage) Bounds() image.Rectangle {
	// Return empty rectangle (Min == Max) which should cause bitmap creation to fail
	return image.Rect(0, 0, 0, 0)
}

func (m *mockBrokenImage) At(x, y int) color.Color {
	return color.RGBA{0, 0, 0, 255}
}

func TestVerify(t *testing.T) {
	tests := []struct {
		name         string
		setupImage   func() []byte
		expectedData string
		wantErr      bool
		checkVerErr  bool
	}{
		{
			name: "valid QR code matches expected data",
			setupImage: func() []byte {
				qr, err := qrcode.Encode("test data", qrcode.Medium, 256)
				if err != nil {
					t.Fatalf("failed to generate test QR: %v", err)
				}
				return qr
			},
			expectedData: "test data",
			wantErr:      false,
		},
		{
			name: "valid QR code does not match expected data",
			setupImage: func() []byte {
				qr, err := qrcode.Encode("actual data", qrcode.Medium, 256)
				if err != nil {
					t.Fatalf("failed to generate test QR: %v", err)
				}
				return qr
			},
			expectedData: "expected data",
			wantErr:      true,
			checkVerErr:  true,
		},
		{
			name: "invalid PNG data",
			setupImage: func() []byte {
				return []byte("not a PNG image")
			},
			expectedData: "test data",
			wantErr:      true,
			checkVerErr:  false,
		},
		{
			name: "valid PNG but no QR code",
			setupImage: func() []byte {
				// Create a simple solid color image
				img := image.NewRGBA(image.Rect(0, 0, 100, 100))
				for y := 0; y < 100; y++ {
					for x := 0; x < 100; x++ {
						img.Set(x, y, color.RGBA{255, 0, 0, 255})
					}
				}
				var buf bytes.Buffer
				if err := png.Encode(&buf, img); err != nil {
					t.Fatalf("failed to encode test image: %v", err)
				}
				return buf.Bytes()
			},
			expectedData: "test data",
			wantErr:      true,
			checkVerErr:  false,
		},
		{
			name: "single space matches",
			setupImage: func() []byte {
				qr, err := qrcode.Encode(" ", qrcode.Medium, 256)
				if err != nil {
					t.Fatalf("failed to generate test QR: %v", err)
				}
				return qr
			},
			expectedData: " ",
			wantErr:      false,
		},
		{
			name: "unicode content matches",
			setupImage: func() []byte {
				qr, err := qrcode.Encode("Hello 世界 🌍", qrcode.Medium, 256)
				if err != nil {
					t.Fatalf("failed to generate test QR: %v", err)
				}
				return qr
			},
			expectedData: "Hello 世界 🌍",
			wantErr:      false,
		},
		{
			name: "URL data matches",
			setupImage: func() []byte {
				url := "https://example.com/path?query=value&other=test"
				qr, err := qrcode.Encode(url, qrcode.Medium, 256)
				if err != nil {
					t.Fatalf("failed to generate test QR: %v", err)
				}
				return qr
			},
			expectedData: "https://example.com/path?query=value&other=test",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qrImage := tt.setupImage()
			err := Verify(qrImage, tt.expectedData)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Verify() expected error, got nil")
					return
				}

				if tt.checkVerErr {
					var verErr *VerificationError
					if !errors.As(err, &verErr) {
						t.Errorf("Verify() expected VerificationError, got: %v", err)
						return
					}
					if verErr.Original != tt.expectedData {
						t.Errorf("VerificationError.Original = %q, want %q", verErr.Original, tt.expectedData)
					}
					if verErr.Recovery != Medium {
						t.Errorf("VerificationError.Recovery = %v, want Medium", verErr.Recovery)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Verify() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		setupImg func() image.Image
		wantData string
		wantErr  bool
	}{
		{
			name: "decode valid QR code",
			setupImg: func() image.Image {
				qr, err := qrcode.Encode("decode test", qrcode.Medium, 256)
				if err != nil {
					t.Fatalf("failed to generate test QR: %v", err)
				}
				img, err := png.Decode(bytes.NewReader(qr))
				if err != nil {
					t.Fatalf("failed to decode PNG: %v", err)
				}
				return img
			},
			wantData: "decode test",
			wantErr:  false,
		},
		{
			name: "decode fails on non-QR image",
			setupImg: func() image.Image {
				img := image.NewRGBA(image.Rect(0, 0, 50, 50))
				for y := 0; y < 50; y++ {
					for x := 0; x < 50; x++ {
						img.Set(x, y, color.RGBA{0, 0, 255, 255})
					}
				}
				return img
			},
			wantData: "",
			wantErr:  true,
		},
		{
			name: "decode QR with special characters",
			setupImg: func() image.Image {
				qr, err := qrcode.Encode("!@#$%^&*()", qrcode.Medium, 256)
				if err != nil {
					t.Fatalf("failed to generate test QR: %v", err)
				}
				img, err := png.Decode(bytes.NewReader(qr))
				if err != nil {
					t.Fatalf("failed to decode PNG: %v", err)
				}
				return img
			},
			wantData: "!@#$%^&*()",
			wantErr:  false,
		},
		{
			name: "decode QR with newlines",
			setupImg: func() image.Image {
				qr, err := qrcode.Encode("line1\nline2\nline3", qrcode.Medium, 256)
				if err != nil {
					t.Fatalf("failed to generate test QR: %v", err)
				}
				img, err := png.Decode(bytes.NewReader(qr))
				if err != nil {
					t.Fatalf("failed to decode PNG: %v", err)
				}
				return img
			},
			wantData: "line1\nline2\nline3",
			wantErr:  false,
		},
		{
			name: "decode numeric QR",
			setupImg: func() image.Image {
				qr, err := qrcode.Encode("12345", qrcode.Medium, 256)
				if err != nil {
					t.Fatalf("failed to generate test QR: %v", err)
				}
				img, err := png.Decode(bytes.NewReader(qr))
				if err != nil {
					t.Fatalf("failed to decode PNG: %v", err)
				}
				return img
			},
			wantData: "12345",
			wantErr:  false,
		},
		{
			name: "decode fails with broken image",
			setupImg: func() image.Image {
				// Use mock image with invalid bounds
				return &mockBrokenImage{}
			},
			wantData: "",
			wantErr:  true,
		},
		{
			name: "decode fails with 1x1 pixel image",
			setupImg: func() image.Image {
				// Create a tiny 1x1 image that can't contain a QR code
				img := image.NewRGBA(image.Rect(0, 0, 1, 1))
				img.Set(0, 0, color.RGBA{0, 0, 0, 255})
				return img
			},
			wantData: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img := tt.setupImg()
			got, err := decode(img)

			if tt.wantErr {
				if err == nil {
					t.Errorf("decode() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("decode() unexpected error: %v", err)
				return
			}

			if got != tt.wantData {
				t.Errorf("decode() = %q, want %q", got, tt.wantData)
			}
		})
	}
}

func TestVerify_EdgeCases(t *testing.T) {
	t.Run("nil image data", func(t *testing.T) {
		err := Verify(nil, "test")
		if err == nil {
			t.Error("Verify() with nil image expected error, got nil")
		}
	})

	t.Run("empty image data", func(t *testing.T) {
		err := Verify([]byte{}, "test")
		if err == nil {
			t.Error("Verify() with empty image expected error, got nil")
		}
	})

	t.Run("case sensitivity", func(t *testing.T) {
		qr, err := qrcode.Encode("Test", qrcode.Medium, 256)
		if err != nil {
			t.Fatalf("failed to generate test QR: %v", err)
		}

		err = Verify(qr, "test")
		if err == nil {
			t.Error("Verify() expected case-sensitive mismatch error, got nil")
		}

		var verErr *VerificationError
		if !errors.As(err, &verErr) {
			t.Errorf("Expected VerificationError for case mismatch, got: %v", err)
		}
	})

	t.Run("whitespace sensitivity", func(t *testing.T) {
		qr, err := qrcode.Encode("test", qrcode.Medium, 256)
		if err != nil {
			t.Fatalf("failed to generate test QR: %v", err)
		}

		err = Verify(qr, "test ")
		if err == nil {
			t.Error("Verify() expected whitespace mismatch error, got nil")
		}

		var verErr *VerificationError
		if !errors.As(err, &verErr) {
			t.Errorf("Expected VerificationError for whitespace mismatch, got: %v", err)
		}
	})
}
