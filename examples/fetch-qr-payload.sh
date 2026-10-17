#!/bin/bash
# Example: Fetch QR registration payload and generate QR code locally
#
# Usage:
#   ./fetch-qr-payload.sh <headscale-url> <auth-id>
#
# Example:
#   ./fetch-qr-payload.sh https://headscale.example.com hskey-authreq-abc123

set -euo pipefail

if [ $# -ne 2 ]; then
    echo "Usage: $0 <headscale-url> <auth-id>"
    echo "Example: $0 https://headscale.example.com hskey-authreq-abc123"
    exit 1
fi

HEADSCALE_URL="$1"
AUTH_ID="$2"

# Fetch the QR payload
echo "Fetching QR registration payload from ${HEADSCALE_URL}..."
PAYLOAD=$(curl -sf "${HEADSCALE_URL}/api/v1/registration/qr/${AUTH_ID}")

if [ $? -ne 0 ]; then
    echo "Error: Failed to fetch QR payload"
    exit 1
fi

# Display the payload
echo "Registration Payload:"
echo "$PAYLOAD" | jq .

# Generate QR code (requires qrencode to be installed)
if command -v qrencode &> /dev/null; then
    echo ""
    echo "QR Code (scan this with Headplane):"
    echo ""
    echo "$PAYLOAD" | qrencode -t ansiutf8
else
    echo ""
    echo "Install 'qrencode' to display the QR code in terminal:"
    echo "  sudo apt-get install qrencode  # Debian/Ubuntu"
    echo "  brew install qrencode          # macOS"
fi
