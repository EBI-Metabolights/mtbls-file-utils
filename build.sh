#!/bin/bash

rm -rf build

# macOS - Intel (amd64)
GOOS=darwin GOARCH=amd64 go build -o build/mac/amd64/mtbls-compress compress.go

# macOS - Apple Silicon (arm64)
GOOS=darwin GOARCH=arm64 go build -o build/mac/arm64/mtbls-compress compress.go

# Linux - Intel (amd64)
GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/mtbls-compress compress.go

# Linux - ARM (arm64)
GOOS=linux GOARCH=arm64 go build -o build/linux/arm64/mtbls-compress compress.go

# Windows - Intel (amd64)
GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/mtbls-compress.exe compress.go

# Windows - ARM (arm64)
GOOS=windows GOARCH=arm64 go build -o build/windows/arm64/mtbls-compress.exe compress.go


# macOS - Intel (amd64)
GOOS=darwin GOARCH=amd64 go build -o build/mac/amd64/mtbls-rename rename.go

# macOS - Apple Silicon (arm64)
GOOS=darwin GOARCH=arm64 go build -o build/mac/arm64/mtbls-rename rename.go

# Linux - Intel (amd64)
GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/mtbls-rename rename.go

# Linux - ARM (arm64)
GOOS=linux GOARCH=arm64 go build -o build/linux/arm64/mtbls-rename rename.go

# Windows - Intel (amd64)
GOOS=windows GOARCH=amd64 go build -o build/windows/amd64/mtbls-rename.exe rename.go

# Windows - ARM (arm64)
GOOS=windows GOARCH=arm64 go build -o build/windows/arm64/mtbls-rename.exe rename.go


chmod +x build/*