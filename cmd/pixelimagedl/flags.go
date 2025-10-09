package main

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/jalavosus/pixelimagedl/pkg/pixelimagedl"
)

var (
	downloadTypeFlag = cli.StringFlag{
		Name:     "imagetype",
		Usage:    "`type` of image zip to download (factory or ota)",
		Aliases:  []string{"t", "type"},
		Required: false,
		Value:    pixelimagedl.Factory.String(),
	}
	deviceNameFlag = cli.StringFlag{
		Name:     "device",
		Usage:    "`name` (or codename) of the device to download an image for",
		Aliases:  []string{"d"},
		Required: true,
	}
	sourceDeviceFlag = cli.StringFlag{
		Name:     "source",
		Usage:    "`source` device (or codename) to port firmware from",
		Aliases:  []string{"s"},
		Required: true,
	}
	targetDeviceFlag = cli.StringFlag{
		Name:     "target",
		Usage:    "`target` device (or codename) to port firmware to",
		Aliases:  []string{"tgt"},
		Required: true,
	}
	portingAlgorithmFlag = cli.StringFlag{
		Name:     "algorithm",
		Usage:    "`algorithm` to use for porting (basic, advanced, experimental)",
		Aliases:  []string{"a", "algo"},
		Required: false,
		Value:    "basic",
	}
	downloadTimeoutFlag = cli.DurationFlag{
		Name:     "timeout",
		Usage:    "`timeout` for file downloads",
		Value:    15 * time.Minute,
		Required: false,
	}
	outDirFlag = cli.StringFlag{
		Name:     "outdir",
		Usage:    "`dir`ectory to place downloaded image file in",
		Aliases:  []string{"o"},
		Required: false,
		Value:    absPath(),
	}
	streamingFlag = cli.BoolFlag{
		Name:     "streaming",
		Usage:    "Enable streaming download with real-time porting and tqdm progress",
		Aliases:  []string{"stream"},
		Required: false,
		Value:    true, // Default to streaming mode
	}
)

type ParsedFlags struct {
	OutDir            string
	Device            pixelimagedl.Pixel
	SourceDevice      pixelimagedl.Pixel
	TargetDevice      pixelimagedl.Pixel
	DownloadTimeout   time.Duration
	DownloadType      pixelimagedl.DownloadType
	PortingAlgorithm  string
	StreamingEnabled  bool
}

func WithFlags(fn func(context.Context, ParsedFlags) error) cli.ActionFunc {
	return func(ctx context.Context, command *cli.Command) error {
		parsedFlags, flagsErr := parseFlags(command)
		if flagsErr != nil {
			return flagsErr
		}

		return fn(ctx, parsedFlags)
	}
}

func parseFlags(cmd *cli.Command) (ParsedFlags, error) {
	var parsedFlags ParsedFlags

	rawDeviceName := cmd.String(deviceNameFlag.Name)
	rawImageKind := cmd.String(downloadTypeFlag.Name)
	rawSourceDevice := cmd.String(sourceDeviceFlag.Name)
	rawTargetDevice := cmd.String(targetDeviceFlag.Name)
	portingAlgorithm := cmd.String(portingAlgorithmFlag.Name)

	// For port command, use source and target devices
	if cmd.Name == "port" {
		sourceDevice, ok := validateDevice(rawSourceDevice)
		if !ok {
			return parsedFlags, errors.Errorf("invalid source device name %[1]s. Allowed values: %[2]s", rawSourceDevice, strings.Join(allowedDeviceNames, ", "))
		}

		targetDevice, ok := validateDevice(rawTargetDevice)
		if !ok {
			return parsedFlags, errors.Errorf("invalid target device name %[1]s. Allowed values: %[2]s", rawTargetDevice, strings.Join(allowedDeviceNames, ", "))
		}

		parsedFlags.SourceDevice = sourceDevice
		parsedFlags.TargetDevice = targetDevice
		parsedFlags.PortingAlgorithm = portingAlgorithm
	} else if cmd.Name == "specific" {
		// For specific command, device is hardcoded (no validation needed)
		parsedFlags.Device = pixelimagedl.Pixel7Pro // Hardcoded to cheetah
	} else {
		// For download/list commands, use the device flag
		deviceName, ok := validateDevice(rawDeviceName)
		if !ok {
			return parsedFlags, errors.Errorf("invalid device name %[1]s. Allowed values: %[2]s", rawDeviceName, strings.Join(allowedDeviceNames, ", "))
		}
		parsedFlags.Device = deviceName
	}

	downloadKind, ok := validateImageKind(rawImageKind)
	if !ok {
		return parsedFlags, errors.Errorf("invalid download kind %[1]s. Allowed values: %[2]s", rawImageKind, strings.Join(allowedDownloadTypes, ", "))
	}

	downloadTimeout := cmd.Duration(downloadTimeoutFlag.Name)
	outDir := cmd.String(outDirFlag.Name)
	streamingEnabled := cmd.Bool(streamingFlag.Name)

	parsedFlags.DownloadType = downloadKind
	parsedFlags.DownloadTimeout = downloadTimeout
	parsedFlags.OutDir = outDir
	parsedFlags.StreamingEnabled = streamingEnabled

	return parsedFlags, nil
}
