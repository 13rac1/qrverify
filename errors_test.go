package qrverify

import (
	"testing"
)

func TestVerificationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		original string
		decoded  string
		want     string
	}{
		{
			name:     "simple mismatch",
			original: "hello",
			decoded:  "helo",
			want:     `verification failed: decoded length 4 does not match original length 5`,
		},
		{
			name:     "empty strings",
			original: "",
			decoded:  "",
			want:     `verification failed: decoded length 0 does not match original length 0`,
		},
		{
			name:     "unicode content",
			original: "Hello 世界",
			decoded:  "Hello World",
			want:     `verification failed: decoded length 11 does not match original length 12`,
		},
		{
			name:     "same length different content",
			original: "abcd",
			decoded:  "abce",
			want:     `verification failed: decoded length 4 does not match original length 4`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &VerificationError{
				Original: tt.original,
				Decoded:  tt.decoded,
			}
			if got := err.Error(); got != tt.want {
				t.Errorf("VerificationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerificationError_Detail(t *testing.T) {
	tests := []struct {
		name     string
		original string
		decoded  string
		want     string
	}{
		{
			name:     "simple mismatch",
			original: "hello",
			decoded:  "helo",
			want:     `verification failed: decoded "helo" does not match original "hello"`,
		},
		{
			name:     "unicode content",
			original: "Hello 世界",
			decoded:  "Hello World",
			want:     `verification failed: decoded "Hello World" does not match original "Hello 世界"`,
		},
		{
			name:     "special characters",
			original: `line1\nline2`,
			decoded:  "line1line2",
			want:     `verification failed: decoded "line1line2" does not match original "line1\\nline2"`,
		},
		{
			name:     "quotes and escapes",
			original: `"quoted"`,
			decoded:  "quoted",
			want:     `verification failed: decoded "quoted" does not match original "\"quoted\""`,
		},
		{
			name:     "emoji",
			original: "Hello 😀",
			decoded:  "Hello",
			want:     `verification failed: decoded "Hello" does not match original "Hello 😀"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &VerificationError{
				Original: tt.original,
				Decoded:  tt.decoded,
			}
			if got := err.Detail(); got != tt.want {
				t.Errorf("VerificationError.Detail() = %v, want %v", got, tt.want)
			}
		})
	}
}
