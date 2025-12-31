package qrverify

import "fmt"

// VerificationError indicates decoded data does not match input.
type VerificationError struct {
	Original string   // What was encoded
	Decoded  string   // What was decoded
	Recovery Recovery // Recovery level used
}

// Error returns a safe error message without exposing data content.
func (e *VerificationError) Error() string {
	return fmt.Sprintf("verification failed: decoded length %d does not match original length %d (recovery: %v)",
		len(e.Decoded), len(e.Original), e.Recovery)
}

// Detail returns an error message with full data content for debugging.
func (e *VerificationError) Detail() string {
	return fmt.Sprintf("verification failed: decoded %q does not match original %q (recovery: %v)",
		e.Decoded, e.Original, e.Recovery)
}
