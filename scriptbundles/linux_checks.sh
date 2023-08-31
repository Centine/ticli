#!/bin/bash

# Linux specific checks

# Call checks shared between *nix platforms
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

nix_retn_text=$( "$SCRIPT_DIR/nix_checks.sh" )
nix_retn_code=$?

echo "$nix_retn_text"
if [ $nix_retn_code -ne 0 ]; then
    echo "nix_checks.sh failed;FAIL;Error code $nix_retn_code"
    exit 1
fi

echo "nixchecks.sh;SUCCESS;Internal check passed"