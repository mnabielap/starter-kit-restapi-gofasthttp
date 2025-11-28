#!/bin/bash
set -e

echo "----------------------------------------"
echo "🚀 Starting GoFastHTTP Application..."
echo "----------------------------------------"

# Optional: You could add logic here to wait for Postgres if needed,
# but GORM usually retries connections or fails fast.

# Run the application
# Using exec allows the app to receive signals (like SIGTERM) correctly
exec ./main