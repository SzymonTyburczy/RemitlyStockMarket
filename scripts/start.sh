#!/usr/bin/env bash
# Usage: ./scripts/start.sh <PORT>
# Example: ./scripts/start.sh 8080
set -euo pipefail

PORT="${1:?Usage: $0 <port>}"
export PORT

# Generate nginx.conf from template — substitutes ${PORT}
sed "s/\${PORT}/${PORT}/g" nginx/nginx.conf.template > nginx/nginx.conf

docker compose up --build -d

echo ""
echo "✓ Stock Market running at http://localhost:${PORT}"
echo "  Instances: stock-service-1, stock-service-2, stock-service-3"
echo "  Shared state: Redis"
