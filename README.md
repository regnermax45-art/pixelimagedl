# pixelimagedl

## Install as CLI app

Make sure you have Go 1.18 or later installed, then:

`go install github.com/jalavosus/pixelimagedl/cmd/pixelimagedl@latest`

## Usage

To list available OTA or Factory images for your pixel device:

`pixelimagedl list -d [device] -t [image type]`

Example using OTA images:

```
$> pixelimagedl list -d oriole -t ota ### fetch OTA images for Pixel 6 (oriole)

Output:

{
  "version": "12.0.0",
  "build_number": "SD1A.210817.015.A4",
  "build_date": "Oct 2021",
  "build_comment": "",
  "download_uri": "https://dl.google.com/dl/android/aosp/oriole-ota-sd1a.210817.015.a4-19a77b62.zip",
  "sha256_sum": "19a77b62a80732b0e20d3ac369288fc827554c1a8c109c6e9fa2a420998dd826"
}
{
  "version": "12.0.0",
  "build_number": "SD1A.210817.019.B1",
  "build_date": "Oct 2021",
  "build_comment": "AT&T",
  "download_uri": "https://dl.google.com/dl/android/aosp/oriole-ota-sd1a.210817.019.b1-67195ca0.zip",
  "sha256_sum": "67195ca0fa1f9cfc9f331ab3e4f44fdbf819896059dcb1819c1125edef854c25"
}
[...]
{
  "version": "12.1.0",
  "build_number": "SP2A.220405.004",
  "build_date": "Apr 2022",
  "build_comment": "",
  "download_uri": "https://dl.google.com/dl/android/aosp/oriole-ota-sp2a.220405.004-f019343f.zip",
  "sha256_sum": "f019343fc56cf466249580dfbdee30ae341a88d6429ecedda727097f51d41c28"
}
{
  "version": "12.1.0",
  "build_number": "SP2A.220505.002",
  "build_date": "May 2022",
  "build_comment": "",
  "download_uri": "https://dl.google.com/dl/android/aosp/oriole-ota-sp2a.220505.002-513d254d.zip",
  "sha256_sum": "513d254dc29c6a61f26fda316eededf265b4fc4b0b528d761a852d074632ec71"
}
```

To download an OTA or Factory image for a device:

`pixelimagedl download -d [device] -t [image type]`

Example:

```
$> pixelimagedl download -d pixel4a -t factory

Output:

2022/05/11 13:41:33 latest stable factory image for Pixel 4a is 12.1.0 (SP2A.220505.002)
2022/05/11 13:41:33 downloading factory image from https://dl.google.com/dl/android/aosp/sunfish-sp2a.220505.002-factory-2ca902f1.zip
2022/05/11 13:41:33 saving factory image to sunfish-sp2a.220505.002-factory-2ca902f1.zip
 100% |████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████| (2.0/2.0 GB, 9.546 MB/s)          
2022/05/11 13:08:36 saved 2.0Gb to sunfish-sp2a.220505.002-factory-2ca902f1.zip
2022/05/11 13:08:41 SHA256 sum 2ca902f10574fa45806c598e8cf10584aa2dc3dd0fb1cf75882e50426187221c of downloaded file matches expected
```

### Device names

Currently supported devices (and their codenames), and their corresponding CLI values:

- Pixel 6 ("Oriole")
  - `pixel6`
  - `oriole`
- Pixel 6 Pro ("Raven")
  - `pixel6pro`
  - `raven`
- Pixel 5 ("Redfin")
  - `pixel5`
  - `redfin`
- Pixel 5a ("Barbet")
  - `pixel5a`
  - `barbet`
- Pixel 4 ("Flame")
  - `pixel4`
  - `flame`
- Pixel 4 XL ("Coral")
  - `pixel4xl`
  - `coral`
- Pixel 4a ("Sunfish")
  - `pixel4a`
  - `sunfish`
- Pixel 4a (5G) ("Bramble")
  - `pixel4a5g`
  - `bramble`

## Firmware Porting

⚠️ **WARNING: Firmware porting is experimental and can result in boot failures or device bricking. Always backup your device before flashing ported firmware.**

The `port` command allows you to download firmware from one Pixel device and apply custom algorithms to make it compatible with another device. This is useful for testing new features or bringing newer Android versions to older devices.

### Port Command Usage

`pixelimagedl port -s [source device] -tgt [target device] -t [image type] -a [algorithm]`

#### Parameters:
- `-s, --source`: Source device to port firmware from (required)
- `-tgt, --target`: Target device to port firmware to (required)
- `-t, --imagetype`: Type of image (factory or ota)
- `-a, --algorithm`: Porting algorithm to use (basic, advanced, experimental) - defaults to basic
- `-o, --outdir`: Output directory for the ported firmware

#### Example - Port Pixel 9 Pro firmware to Pixel 7 Pro:

```
$> pixelimagedl port -s pixel9pro -tgt pixel7pro -t factory -a basic

Output:
Starting firmware porting from Pixel 9 Pro to Pixel 7 Pro using basic algorithm
⚠️  WARNING: Firmware porting is experimental and may result in unstable or unusable firmware
⚠️  CROSS-GENERATION PORTING: Porting from Pixel 9 Pro (Pixel 9) to Pixel 7 Pro (Pixel 7)
⚠️  This operation has higher risk of compatibility issues
ℹ️  BASIC ALGORITHM SELECTED
ℹ️  This provides minimal modifications with lower risk
ℹ️  May require additional manual configuration
📥 Downloading source firmware from Pixel 9 Pro...
📦 Using latest factory image: 14.0.0 (AP2A.240605.024)
🔧 Applying basic porting algorithm...
✅ Firmware porting completed successfully!
📁 Ported firmware saved to: ./caiman-ap2a.240605.024-factory-ported-to-pixel7pro-20241009-143052.zip
🔧 Modifications applied: Added target device metadata, Updated build fingerprint, Modified device-specific configurations
⚠️  Warnings: Basic algorithm provides minimal compatibility, May require additional manual modifications, Boot success not guaranteed
```

#### Porting Algorithms:

- **basic**: Minimal modifications with lower risk. Good for testing basic compatibility.
- **advanced**: More extensive modifications including device tree and kernel adjustments. Higher compatibility but increased risk.
- **experimental**: Aggressive modifications for cross-generation porting. Highest risk - use only for development/testing.

#### Supported Devices for Porting:

Porting is supported between Pixel 6, 7, 8, and 9 series devices. Cross-generation porting (e.g., Pixel 9 to Pixel 7) is possible but has higher risk of compatibility issues.

## TODO

- Implement functionality for listing/downloading available beta versions
- Add support for custom ROM repositories
- Implement more sophisticated porting algorithms
