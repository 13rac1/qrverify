package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/13rac1/qrverify"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "encode":
		encodeCommand(os.Args[2:])
	case "verify":
		verifyCommand(os.Args[2:])
	case "demo":
		demoCommand(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: qrverify <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  encode  Generate a verified QR code")
	fmt.Println("  verify  Verify a QR code image")
	fmt.Println("  demo    Demonstrate encode/verify workflow")
	fmt.Println()
	fmt.Println("Run 'qrverify <command> -h' for command help.")
}

func encodeCommand(args []string) {
	fs := flag.NewFlagSet("encode", flag.ExitOnError)
	output := fs.String("o", "qr.png", "Output file")
	recovery := fs.String("r", "medium", "Recovery level: low, medium, high, highest")
	size := fs.Int("s", 256, "Size in pixels")

	fs.Usage = func() {
		fmt.Println("Usage: qrverify encode <data> [-o output.png] [-r recovery] [-s size]")
		fmt.Println()
		fmt.Println("Flags:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: data argument required")
		fs.Usage()
		os.Exit(1)
	}

	data := fs.Arg(0)

	// Parse recovery level
	r, err := parseRecovery(*recovery)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	opts := &qrverify.EncodeOptions{
		Recovery: r,
		Size:     *size,
	}

	// Use EncodeDetailed to get metadata for output
	result, err := qrverify.EncodeDetailed(data, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding QR code: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*output, result.Image, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created %s (%dx%d, recovery: %s)\n",
		*output, result.Size, result.Size, strings.ToLower(result.Recovery.String()))
}

func verifyCommand(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage: qrverify verify <file.png> <expected-data>")
		fmt.Println()
		fmt.Println("Reads the PNG file and verifies it decodes to the expected data.")
		fmt.Println("Exit 0 on success, exit 1 on failure.")
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if fs.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Error: file and expected-data arguments required")
		fs.Usage()
		os.Exit(1)
	}

	filename := fs.Arg(0)
	expectedData := fs.Arg(1)

	qrImage, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	if err := qrverify.Verify(qrImage, expectedData); err != nil {
		fmt.Fprintf(os.Stderr, "Verification failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Verification passed")
}

func demoCommand(args []string) {
	fs := flag.NewFlagSet("demo", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("Usage: qrverify demo")
		fmt.Println()
		fmt.Println("Demonstrates the encode/verify workflow:")
		fmt.Println("1. Encodes \"Hello, QR World!\" to a temp file")
		fmt.Println("2. Verifies it")
		fmt.Println("3. Shows the result metadata")
		fmt.Println("4. Cleans up")
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	data := "Hello, QR World!"
	fmt.Printf("Demo: Encoding %q\n", data)

	// Create temp file
	tmpfile, err := os.CreateTemp("", "qrverify-demo-*.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp file: %v\n", err)
		os.Exit(1)
	}
	tmpname := tmpfile.Name()
	_ = tmpfile.Close()

	defer func() {
		fmt.Println("Cleaning up...")
		_ = os.Remove(tmpname)
	}()

	// Encode with detailed metadata
	result, err := qrverify.EncodeDetailed(data, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(tmpname, result.Image, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing temp file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created temporary QR code at %s\n", tmpname)
	fmt.Println("Verifying...")

	// Verify
	if err := qrverify.Verify(result.Image, data); err != nil {
		fmt.Fprintf(os.Stderr, "Verification failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Verification passed!")
	fmt.Println("QR Code Details:")
	fmt.Printf("  Recovery: %s\n", result.Recovery.String())
	fmt.Printf("  Size: %dx%d\n", result.Size, result.Size)
	fmt.Println("Done!")
}

func parseRecovery(s string) (qrverify.Recovery, error) {
	switch strings.ToLower(s) {
	case "low":
		return qrverify.Low, nil
	case "medium":
		return qrverify.Medium, nil
	case "high":
		return qrverify.High, nil
	case "highest":
		return qrverify.Highest, nil
	default:
		return 0, fmt.Errorf("invalid recovery level %q, must be: low, medium, high, highest", s)
	}
}
