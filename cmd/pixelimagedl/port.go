package main

import (
	"context"
	"log"

	"github.com/jalavosus/pixelimagedl/pkg/pixelimagedl"
	"github.com/urfave/cli/v3"
)

var portCmd = cli.Command{
	Name:  "port",
	Usage: "Execute Android 16 porting and upload to gofile.io",
	Description: `Execute Android 16 firmware porting between devices and upload to gofile.io:
- Downloads Android 16 firmware for source device
- Ports firmware to target device
- Uploads ported firmware to gofile.io in chunks
- Supports Pixel 9 ↔ Pixel 7 Pro porting`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "source",
			Usage:    "Source device (e.g., pixel9, pixel7pro)",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "target",
			Usage:    "Target device (e.g., pixel7pro, pixel9)",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "upload",
			Value: true,
			Usage: "Upload to gofile.io after porting",
		},
		&cli.IntFlag{
			Name:  "chunk-size",
			Value: 50,
			Usage: "Chunk size in MB for gofile.io upload",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		sourceDevice := c.String("source")
		targetDevice := c.String("target")
		uploadEnabled := c.Bool("upload")
		chunkSizeMB := c.Int("chunk-size")

		log.Printf("🚀 Starting Android 16 porting execution")
		log.Printf("📱 Source device: %s", sourceDevice)
		log.Printf("📱 Target device: %s", targetDevice)
		log.Printf("☁️ Upload to gofile.io: %v", uploadEnabled)
		log.Printf("📦 Chunk size: %d MB", chunkSizeMB)

		// Validate device compatibility
		if !isValidDevicePair(sourceDevice, targetDevice) {
			log.Printf("❌ Invalid device pair: %s → %s", sourceDevice, targetDevice)
			log.Printf("✅ Supported pairs: pixel9↔pixel7pro, tokay↔cheetah")
			return nil
		}

		// Execute Android 16 porting
		if err := pixelimagedl.ExecuteAndroid16PortingWithOptions(sourceDevice, targetDevice, uploadEnabled); err != nil {
			log.Printf("❌ Porting failed: %v", err)
			return err
		}

		log.Printf("🎉 Android 16 porting execution completed successfully!")
		return nil
	},
}

// isValidDevicePair checks if the device pair is supported for porting
func isValidDevicePair(source, target string) bool {
	validPairs := map[string][]string{
		"pixel9":    {"pixel7pro", "cheetah"},
		"tokay":     {"pixel7pro", "cheetah"},
		"pixel7pro": {"pixel9", "tokay"},
		"cheetah":   {"pixel9", "tokay"},
	}

	targets, exists := validPairs[source]
	if !exists {
		return false
	}

	for _, validTarget := range targets {
		if target == validTarget {
			return true
		}
	}

	return false
}
