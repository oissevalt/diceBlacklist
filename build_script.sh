#!/bin/bash
# Default output filename if not provided
OUTPUT_FILENAME="myapp.exe"

# If the first argument is provided, use it as the output filename
if [ $# -gt 0 ]; then
  OUTPUT_FILENAME=$1
fi

# Build the application
GOOS=windows GOARCH=amd64 go build -o "$OUTPUT_FILENAME" .
