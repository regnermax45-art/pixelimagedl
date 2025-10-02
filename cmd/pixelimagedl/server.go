package main

import (
	"context"
	"log"

	"github.com/jalavosus/pixelimagedl/pkg/pixelimagedl"
	"github.com/urfave/cli/v3"
)

var serverCmd = cli.Command{
	Name:  "server",
	Usage: "Start custom firmware server for Android 16 & 17 NOW + device porting",
	Description: `Start a custom server that allows:
- Manual Android 16 & 17 firmware requests
- Cross-device porting between Pixel 9 and Pixel 7 Pro
- Real-time firmware downloads
- Custom firmware URL generation for both versions`,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "port",
			Value: 8080,
			Usage: "Server port",
		},
		&cli.BoolFlag{
			Name:  "custom",
			Value: true,
			Usage: "Enable custom firmware requests",
		},
		&cli.BoolFlag{
			Name:  "porting",
			Value: true,
			Usage: "Enable device porting",
		},
	},
	Action: func(ctx context.Context, c *cli.Command) error {
		config := pixelimagedl.CustomServerConfig{
			EnableCustomRequests: c.Bool("custom"),
			EnableDevicePorting:  c.Bool("porting"),
			ServerPort:          int(c.Int("port")),
			AllowedDevices:      []string{"pixel9", "tokay", "pixel7pro", "cheetah", "panther"},
		}

		log.Printf("🚀 Starting custom firmware server...")
		log.Printf("📱 Custom requests: %v", config.EnableCustomRequests)
		log.Printf("🔄 Device porting: %v", config.EnableDevicePorting)
		log.Printf("🌐 Port: %d", config.ServerPort)
		log.Printf("🔥 Ready for Android 16 & 17 NOW!")

		return pixelimagedl.StartCustomServer(config)
	},
}
