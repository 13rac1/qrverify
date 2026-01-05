package qrverify_test

import (
	"bytes"
	"fmt"
	"image/png"
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

func ExampleDecode() {
	// First create a QR code
	pngBytes, _ := qrverify.Encode("Hello, World!", nil)

	// Decode the PNG to an image
	img, _ := png.Decode(bytes.NewReader(pngBytes))

	// Decode the QR code
	data, err := qrverify.Decode(img)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(data)
	// Output: Hello, World!
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
