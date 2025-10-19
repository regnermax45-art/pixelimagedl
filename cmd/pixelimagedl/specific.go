package main

import (
	"context"

	cli "github.com/urfave/cli/v3"

	"github.com/jalavosus/pixelimagedl/pkg/pixelimagedl"
)

var specificCmd = cli.Command{
	Name:  "specific",
	Usage: "Download specific firmware version (cheetah BD3A.251005.003.W3, Oct 2025) via custom request",
	Flags: []cli.Flag{
		&downloadTypeFlag,
		&downloadTimeoutFlag,
		&outDirFlag,
	},
	Action: WithFlags(specificCmdAction),
}

func specificCmdAction(ctx context.Context, parsedFlags ParsedFlags) error {
	ctx, cancel := context.WithTimeout(ctx, parsedFlags.DownloadTimeout)
	defer cancel()

	// Custom request for specific BD build that may not be in standard listings
	return pixelimagedl.DownloadCustomBuild(
		ctx,
		pixelimagedl.Pixel7Pro,    // cheetah
		"BD3A.251005.003.W3",      // exact build number requested
		"16.0.0",                  // Android version
		"Oct 2025",                // build date
		parsedFlags.DownloadType,
		parsedFlags.OutDir,
	)
}
