package main

import (
	"context"

	cli "github.com/urfave/cli/v3"

	"github.com/jalavosus/pixelimagedl/pkg/pixelimagedl"
)

var portCmd = cli.Command{
	Name:  "port",
	Usage: "Download and port firmware from one Pixel device to another using custom algorithms",
	Flags: []cli.Flag{
		&sourceDeviceFlag,
		&targetDeviceFlag,
		&downloadTypeFlag,
		&downloadTimeoutFlag,
		&outDirFlag,
		&portingAlgorithmFlag,
	},
	Action: WithFlags(portCmdAction),
}

func portCmdAction(ctx context.Context, parsedFlags ParsedFlags) error {
	ctx, cancel := context.WithTimeout(ctx, parsedFlags.DownloadTimeout)
	defer cancel()

	return pixelimagedl.PortFirmware(
		ctx,
		parsedFlags.SourceDevice,
		parsedFlags.TargetDevice,
		parsedFlags.DownloadType,
		parsedFlags.OutDir,
		parsedFlags.PortingAlgorithm,
	)
}
