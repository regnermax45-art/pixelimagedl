#!/bin/bash

echo "🚀 Testing Android 16 & 17 NOW Custom Server"
echo "=============================================="

# Start server in background
echo "Starting custom server..."
./pixelimagedl server --port 8080 &
SERVER_PID=$!
sleep 3

echo ""
echo "📊 Testing server status..."
curl -s http://localhost:8080/status | jq '.'

echo ""
echo "🔥 Requesting Android 17 for Pixel 7 Pro NOW..."
curl -X POST http://localhost:8080/request-firmware \
  -d "device=pixel7pro" \
  -d "version=android17" \
  -d "build=17dp1" | jq '.'

echo ""
echo "🔥 Requesting Android 16 for Pixel 9 NOW..."
curl -X POST http://localhost:8080/request-firmware \
  -d "device=pixel9" \
  -d "version=android16" \
  -d "build=16dp1" | jq '.'

echo ""
echo "🔄 Testing device porting: Pixel 9 → Pixel 7 Pro..."
curl -X POST http://localhost:8080/port-device \
  -d "source_device=pixel9" \
  -d "target_device=pixel7pro" \
  -d "version=android17" | jq '.'

echo ""
echo "🔄 Testing device porting: Pixel 7 Pro → Pixel 9..."
curl -X POST http://localhost:8080/port-device \
  -d "source_device=pixel7pro" \
  -d "target_device=pixel9" \
  -d "version=android17" | jq '.'

# Wait a bit for background processes
sleep 5

# Kill server
kill $SERVER_PID 2>/dev/null

echo ""
echo "🔄 Testing Android 16 device porting: Pixel 7 Pro → Pixel 9..."
curl -X POST http://localhost:8080/port-device \
  -d "source_device=pixel7pro" \
  -d "target_device=pixel9" \
  -d "version=android16" | jq '.'

echo ""
echo "✅ Android 16 & 17 NOW custom server test completed!"
