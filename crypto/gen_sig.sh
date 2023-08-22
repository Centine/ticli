#!/bin/bash

# Check if the filename is provided
if [ -z "$1" ]; then
  echo "Please provide the filename to sign as an argument."
  exit 1
fi

# Check if the private key file exists
PRIVATE_KEY="ticli_private_key.pem"
if [ ! -f "$PRIVATE_KEY" ]; then
  echo "Private key file $PRIVATE_KEY does not exist."
  exit 1
fi

# Input file from the argument
INPUT_FILE="$1"
# Output signature file
OUTPUT_FILE="${INPUT_FILE}.sig"

# Sign the file using OpenSSL and the private key
openssl dgst -sha256 -sign "$PRIVATE_KEY" -out "$OUTPUT_FILE" "$INPUT_FILE"

if [ $? -eq 0 ]; then
  echo "File successfully signed. Signature saved to $OUTPUT_FILE."
else
  echo "An error occurred while signing the file."
  exit 1
fi
