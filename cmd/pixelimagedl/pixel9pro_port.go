package main

import (
	"context"

	cli "github.com/urfave/cli/v3"

	"github.com/jalavosus/pixelimagedl/pkg/pixelimagedl"
)

var pixel9ProPortCmd = cli.Command{
	Name:  "pixel9pro-port",
	Usage: "Download Pixel 9 Pro firmware ported to Pixel 7 Pro using custom porting algorithm",
	Flags: []cli.Flag{
		&downloadTypeFlag,
		&downloadTimeoutFlag,
		&outDirFlag,
	},
	Action: WithFlags(pixel9ProPortCmdAction),
}

func pixel9ProPortCmdAction(ctx context.Context, parsedFlags ParsedFlags) error {
	ctx, cancel := context.WithTimeout(ctx, parsedFlags.DownloadTimeout)
	defer cancel()

	// Use Pixel 9 Pro as source, Pixel 7 Pro as target (hardcoded for this specific porting algorithm)
	sourceDevice := pixelimagedl.Pixel9Pro  // caiman
	targetDevice := pixelimagedl.Pixel7Pro  // cheetah
	
	result, err := pixelimagedl.DownloadPixel9ProPortedFirmware(
		sourceDevice,
		targetDevice, 
		parsedFlags.DownloadType,
		parsedFlags.DownloadTimeout,
		parsedFlags.OutDir,
	)
	
	if err != nil {
		return err
	}

	result.PrettyPrint()
	return nil
}
