#!/bin/bash

# This script demonstrates how to use scallop to scan a Docker image

# Check if an image name is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <image_name>"
  echo "Example: $0 nginx:latest"
  exit 1
fi

IMAGE_NAME=$1
OUTPUT_FORMAT=${2:-text}  # Default to text format if not specified

echo "Scanning Docker image: $IMAGE_NAME"
echo "Output format: $OUTPUT_FORMAT"

# Run scallop
../scallop --image "$IMAGE_NAME" --format "$OUTPUT_FORMAT"
