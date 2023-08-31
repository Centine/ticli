#!/bin/bash
set -u # exit out on script on first use of an uninitialised variable

# Array of binaries to check
binaries=("kubectl" "go-aws-sso")

# Loop through binaries and check if they are in the path and executable
for binary in "${binaries[@]}"; do
  if command -v $binary > /dev/null 2>&1; then
    echo "$binary binary;SUCCESS;Found in path and executable"
  else
    echo "$binary binary;FAIL;Not found in path or not executable"
  fi
done

# Check if the .aws directory exists
aws_dir="$HOME/.aws"
if [ -d "$aws_dir" ]; then
  echo ".aws directory;SUCCESS;Exists in home directory"

  # Check for credentials and config files
  for file_path in "$aws_dir/credentials" "$aws_dir/config"; do
    if [ -f "$file_path" ] && [ -s "$file_path" ] && file -b --mime-encoding "$file_path" | grep -q 'ascii'; then
        grep -E -vn '^\s*([#;]|$)|^\s*\[.*\]\s*$|^\s*[^#;=\s]+\s*=\s*[^#;]*' "$file_path" > /dev/null 2>&1
        if [ $? -eq 0 ]; then
          echo "$file_path;SUCCESS;Appears to be a valid ASCII file in INI format"
        else
          echo "$file_path;FAIL;Not a valid ASCII file or not in INI format"
        fi
    else
      echo "$file_path;FAIL;Not a valid ASCII file or does not exist"
    fi
  done
else
  echo ".aws directory;FAIL;Does not exist in home directory"
fi
