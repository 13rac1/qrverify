package qrverify

import "fmt"

// VerificationError indicates decoded data does not match input.
type VerificationError struct {
	Original string   // What was encoded
	Decoded  string   // What was decoded
	Recovery Recovery // Recovery level used
}

func (e *VerificationError) Error() string {
	return fmt.Sprintf("verification failed: decoded %q does not match original %q (recovery: %v)",
		e.Decoded, e.Original, e.Recovery)
}
