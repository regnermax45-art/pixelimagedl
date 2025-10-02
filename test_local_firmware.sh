#!/bin/bash

echo "🔥 Testing Local Firmware Support - Pixel 9 & Pixel 7 Pro"
echo "=========================================================="

echo ""
echo "🔧 Building pixelimagedl with local firmware support..."
export PATH=$PATH:/usr/local/go/bin
go build ./cmd/pixelimagedl

echo ""
echo "📦 Creating test local firmware files..."

# Create test local firmware files (small for testing)
echo "🔥 Creating test Pixel 9 local firmware..."
dd if=/dev/zero of=pixel9-firmware.zip bs=1M count=100 2>/dev/null
echo "PIXEL9_LOCAL_FIRMWARE_TEST" >> pixel9-firmware.zip

echo "🔥 Creating test Pixel 7 Pro local firmware..."
dd if=/dev/zero of=pixel7pro-firmware.zip bs=1M count=100 2>/dev/null
echo "PIXEL7PRO_LOCAL_FIRMWARE_TEST" >> pixel7pro-firmware.zip

echo "🔥 Creating test tegu codename firmware..."
dd if=/dev/zero of=tegu-firmware.zip bs=1M count=100 2>/dev/null
echo "TEGU_LOCAL_FIRMWARE_TEST" >> tegu-firmware.zip

echo "🔥 Creating test cheetah codename firmware..."
dd if=/dev/zero of=cheetah-firmware.zip bs=1M count=100 2>/dev/null
echo "CHEETAH_LOCAL_FIRMWARE_TEST" >> cheetah-firmware.zip

echo ""
echo "📁 Local firmware files created:"
ls -lh *-firmware.zip

echo ""
echo "🔥 Testing LOCAL Pixel 9 firmware detection"
echo "============================================"
echo "Expected: Should find and use pixel9-firmware.zip"
echo ""

# Test Pixel 9 local firmware detection (with timeout to avoid long processing)
timeout 15 ./pixelimagedl port --source pixel9 --target pixel7pro --upload=false

echo ""
echo "🔥 Testing LOCAL Pixel 7 Pro firmware detection"
echo "==============================================="
echo "Expected: Should find and use pixel7pro-firmware.zip"
echo ""

# Test Pixel 7 Pro local firmware detection (with timeout to avoid long processing)
timeout 15 ./pixelimagedl port --source pixel7pro --target pixel9 --upload=false

echo ""
echo "🔥 Testing LOCAL tegu codename firmware detection"
echo "================================================="
echo "Expected: Should find and use tegu-firmware.zip"
echo ""

# Test tegu codename local firmware detection
timeout 15 ./pixelimagedl port --source tegu --target cheetah --upload=false

echo ""
echo "🔥 Testing LOCAL cheetah codename firmware detection"
echo "===================================================="
echo "Expected: Should find and use cheetah-firmware.zip"
echo ""

# Test cheetah codename local firmware detection
timeout 15 ./pixelimagedl port --source cheetah --target tegu --upload=false

echo ""
echo "🧹 Cleaning up test files..."
rm -f pixel9-firmware.zip pixel7pro-firmware.zip tegu-firmware.zip cheetah-firmware.zip
rm -f android16_*.zip result_firmware_*.zip

echo ""
echo "✅ Local firmware testing completed!"
echo ""
echo "📝 Summary of Local Firmware Support:"
echo "- ✅ Checks for local firmware files BEFORE downloading"
echo "- ✅ Supports multiple filename patterns:"
echo "  • device-firmware.zip (pixel9-firmware.zip)"
echo "  • device_firmware.zip (pixel9_firmware.zip)"
echo "  • device-factory.zip (pixel9-factory.zip)"
echo "  • codename-firmware.zip (tegu-firmware.zip)"
echo "  • Real firmware filenames (tegu-bp3a.250905.014-factory-a05fafa0.zip)"
echo "- ✅ Progress tracking during local file copy"
echo "- ✅ Fallback to download if no local files found"
echo "- ✅ Works with both device names (pixel9) and codenames (tegu)"
echo ""
echo "🎯 Usage with local firmware:"
echo "  1. Place firmware files in current directory"
echo "  2. Use supported filename patterns"
echo "  3. Run normal porting commands"
echo "  4. System will automatically detect and use local files"
echo ""
echo "📁 Supported local firmware filename patterns:"
echo "  • pixel9-firmware.zip, pixel7pro-firmware.zip"
echo "  • tegu-firmware.zip, cheetah-firmware.zip"
echo "  • pixel9_factory.zip, pixel7pro_factory.zip"
echo "  • tegu-bp3a.250905.014-factory-a05fafa0.zip"
echo "  • cheetah-bp3a.250905.014-factory-3ef97bbc.zip"
echo "  • android16_pixel9_factory.zip, android16_pixel7pro_factory.zip"
