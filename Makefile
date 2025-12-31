.PHONY: build test lint fmt vet cover clean

build:
	go build ./...
	go build -o qrverify ./cmd/qrverify

test:
	go test ./...

lint: fmt vet

fmt:
	go fmt ./...

vet:
	go vet ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -f qrverify coverage.out coverage.html
