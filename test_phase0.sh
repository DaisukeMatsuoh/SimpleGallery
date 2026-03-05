#!/bin/bash

# Test Phase 0 deliverables

echo "Testing Phase 0: Basic Project Structure"
echo "=========================================="

# 1. Check if binary exists
if [ -f bin/simple-gallery ]; then
    echo "✓ Binary built: bin/simple-gallery"
else
    echo "✗ Binary not found"
    exit 1
fi

# 2. Create directories for storage
mkdir -p /tmp/simple-gallery/media /tmp/simple-gallery/thumbs

# 3. Start server in background
echo "Starting server..."
./bin/simple-gallery -config config.toml &
SERVER_PID=$!

# Wait for server to start
sleep 2

# 4. Test healthz endpoint
echo "Testing /healthz endpoint..."
RESPONSE=$(curl -s http://localhost:8080/healthz)
if echo "$RESPONSE" | grep -q '"status":"ok"'; then
    echo "✓ /healthz endpoint works: $RESPONSE"
else
    echo "✗ /healthz endpoint failed: $RESPONSE"
    kill $SERVER_PID
    exit 1
fi

# 5. Check if database was created
echo "Checking database creation..."
if [ -f /tmp/simple-gallery/simple-gallery.db ]; then
    echo "✓ Database created: /tmp/simple-gallery/simple-gallery.db"
else
    echo "✗ Database not found"
    kill $SERVER_PID
    exit 1
fi

# 6. Stop server
echo "Stopping server..."
kill $SERVER_PID 2>/dev/null || true
sleep 1

echo ""
echo "=========================================="
echo "All Phase 0 tests passed! ✓"
