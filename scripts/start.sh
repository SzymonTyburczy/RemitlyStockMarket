#!/usr/bin/env bash
# Usage: ./scripts/start.sh <PORT>
# Example: ./scripts/start.sh 8080
set -euo pipefail

PORT="${1:?Usage: $0 <port>}"
export PORT
envsubst '${PORT}' < nginx/nginx.conf.template > nginx/nginx.conf
PORT=${PORT} docker compose up --build --scale stock-service=3 -d
echo "Stock Market running at http://localhost:${PORT}"
