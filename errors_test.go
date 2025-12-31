package qrverify

import (
	"testing"
)

func TestVerificationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		original string
		decoded  string
		recovery Recovery
		want     string
	}{
		{
			name:     "simple mismatch",
			original: "hello",
			decoded:  "helo",
			recovery: Medium,
			want:     `verification failed: decoded length 4 does not match original length 5 (recovery: Medium)`,
		},
		{
			name:     "empty strings",
			original: "",
			decoded:  "",
			recovery: Low,
			want:     `verification failed: decoded length 0 does not match original length 0 (recovery: Low)`,
		},
		{
			name:     "unicode content",
			original: "Hello 世界",
			decoded:  "Hello World",
			recovery: High,
			want:     `verification failed: decoded length 11 does not match original length 12 (recovery: High)`,
		},
		{
			name:     "same length different content",
			original: "abcd",
			decoded:  "abce",
			recovery: Highest,
			want:     `verification failed: decoded length 4 does not match original length 4 (recovery: Highest)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &VerificationError{
				Original: tt.original,
				Decoded:  tt.decoded,
				Recovery: tt.recovery,
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
		recovery Recovery
		want     string
	}{
		{
			name:     "simple mismatch",
			original: "hello",
			decoded:  "helo",
			recovery: Medium,
			want:     `verification failed: decoded "helo" does not match original "hello" (recovery: Medium)`,
		},
		{
			name:     "unicode content",
			original: "Hello 世界",
			decoded:  "Hello World",
			recovery: High,
			want:     `verification failed: decoded "Hello World" does not match original "Hello 世界" (recovery: High)`,
		},
		{
			name:     "special characters",
			original: `line1\nline2`,
			decoded:  "line1line2",
			recovery: Highest,
			want:     `verification failed: decoded "line1line2" does not match original "line1\\nline2" (recovery: Highest)`,
		},
		{
			name:     "quotes and escapes",
			original: `"quoted"`,
			decoded:  "quoted",
			recovery: Low,
			want:     `verification failed: decoded "quoted" does not match original "\"quoted\"" (recovery: Low)`,
		},
		{
			name:     "emoji",
			original: "Hello 😀",
			decoded:  "Hello",
			recovery: High,
			want:     `verification failed: decoded "Hello" does not match original "Hello 😀" (recovery: High)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &VerificationError{
				Original: tt.original,
				Decoded:  tt.decoded,
				Recovery: tt.recovery,
			}
			if got := err.Detail(); got != tt.want {
				t.Errorf("VerificationError.Detail() = %v, want %v", got, tt.want)
			}
		})
	}
}
