#!/bin/bash

echo "🚀 Testing Android 16 Porting + GoFile.io Upload"
echo "================================================="

echo ""
echo "🔧 Building pixelimagedl with porting support..."
go build ./cmd/pixelimagedl

echo ""
echo "📱 Testing Android 16 porting: Pixel 9 → Pixel 7 Pro"
echo "======================================================"
./pixelimagedl port --source pixel9 --target pixel7pro

echo ""
echo "📱 Testing Android 16 porting: Pixel 7 Pro → Pixel 9"
echo "======================================================"
./pixelimagedl port --source pixel7pro --target pixel9

echo ""
echo "📱 Testing invalid device pair (should fail gracefully)"
echo "======================================================="
./pixelimagedl port --source pixel8 --target pixel9

echo ""
echo "📋 Listing generated files..."
ls -lh *.zip 2>/dev/null || echo "No firmware files found"

echo ""
echo "✅ Android 16 porting test completed!"
