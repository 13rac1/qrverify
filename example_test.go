package qrverify_test

import (
	"fmt"
	"os"

	"github.com/13rac1/qrverify"
)

func ExampleEncode() {
	opts := &qrverify.EncodeOptions{
		Recovery: qrverify.High,
		Size:     512,
	}
	png, err := qrverify.Encode("Hello, World!", opts)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Generated %d byte QR code\n", len(png))
	// Output: Generated 1781 byte QR code
}

func ExampleEncodeDetailed() {
	result, err := qrverify.EncodeDetailed("https://example.com", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Recovery: %v\n", result.Recovery)
	// Output: Recovery: Medium
}

func ExampleVerify() {
	// First create a QR code
	png, _ := qrverify.Encode("test data", nil)

	// Then verify it
	err := qrverify.Verify(png, "test data")
	if err != nil {
		fmt.Println("Verification failed:", err)
		return
	}
	fmt.Println("Verification passed")
	// Output: Verification passed
}

func ExampleEncodeToFile() {
	tmpfile, _ := os.CreateTemp("", "qr-*.png")
	defer func() { _ = os.Remove(tmpfile.Name()) }()
	_ = tmpfile.Close()

	err := qrverify.EncodeToFile("https://example.com", tmpfile.Name(), nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("File created successfully")
	// Output: File created successfully
}
