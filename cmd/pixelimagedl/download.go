package main

import (
	"context"

	cli "github.com/urfave/cli/v3"

	"github.com/jalavosus/pixelimagedl/pkg/pixelimagedl"
)

var downloadCmd = cli.Command{
	Name: "download",
	Usage: "Download firmware with optional device porting",
	Flags: []cli.Flag{
		&downloadTypeFlag,
		&deviceNameFlag,
		&downloadTimeoutFlag,
		&outDirFlag,
		&cli.StringFlag{
			Name:  "port-from",
			Usage: "Source device for porting (e.g., pixel9)",
		},
		&cli.StringFlag{
			Name:  "port-to",
			Usage: "Target device for porting (e.g., pixel7pro)",
		},
		&cli.BoolFlag{
			Name:  "android17-now",
			Usage: "Force Android 17 download NOW (custom server mode)",
		},
	},
	Action: WithFlags(downloadCmdAction),
}

func downloadCmdAction(ctx context.Context, parsedFlags ParsedFlags) error {
	ctx, cancel := context.WithTimeout(ctx, parsedFlags.DownloadTimeout)
	defer cancel()

	return pixelimagedl.DownloadLatest(
		ctx,
		parsedFlags.Device,
		parsedFlags.DownloadType,
		parsedFlags.OutDir,
	)
}
