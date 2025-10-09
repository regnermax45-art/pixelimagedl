package pixelimagedl

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

func PortFirmware(ctx context.Context, sourceDevice, targetDevice Pixel, downloadType DownloadType, outDir, algorithmStr string) error {
	// Parse the algorithm
	algorithm, err := parseAlgorithm(algorithmStr)
	if err != nil {
		return err
	}

	// Validate device compatibility
	if err := validatePortingCompatibility(sourceDevice, targetDevice); err != nil {
		return err
	}

	config := PortingConfig{
		SourceDevice:    sourceDevice,
		TargetDevice:    targetDevice,
		Algorithm:       algorithm,
		DownloadType:    downloadType,
		OutputDirectory: outDir,
	}

	log.Printf("Starting firmware porting from %s to %s using %s algorithm\n", sourceDevice.String(), targetDevice.String(), algorithm.String())
	log.Printf("⚠️  WARNING: Firmware porting is experimental and may result in unstable or unusable firmware\n")

	result, err := executePorting(ctx, config)
	if err != nil {
		return errors.Wrap(err, "failed to execute porting")
	}

	log.Printf("✅ Firmware porting completed successfully!\n")
	log.Printf("📁 Ported firmware saved to: %s\n", result.PortedImagePath)
	log.Printf("🔧 Modifications applied: %s\n", strings.Join(result.Modifications, ", "))

	if len(result.Warnings) > 0 {
		log.Printf("⚠️  Warnings: %s\n", strings.Join(result.Warnings, "; "))
	}

	return nil
}

func parseAlgorithm(algorithmStr string) (PortingAlgorithm, error) {
	switch strings.ToLower(algorithmStr) {
	case "basic":
		return BasicAlgorithm, nil
	case "advanced":
		return AdvancedAlgorithm, nil
	case "experimental":
		return ExperimentalAlgorithm, nil
	default:
		return BasicAlgorithm, fmt.Errorf("unknown algorithm: %s. Supported: basic, advanced, experimental", algorithmStr)
	}
}

func validatePortingCompatibility(source, target Pixel) error {
	// Basic compatibility check - same generation devices are more likely to be compatible
	sourceGen := getDeviceGeneration(source)
	targetGen := getDeviceGeneration(target)

	if sourceGen != targetGen {
		log.Printf("⚠️  Porting between different generations (%s → %s) may have limited compatibility\n", sourceGen, targetGen)
	}

	// Check if devices are the same
	if source == target {
		return errors.New("source and target devices cannot be the same")
	}

	return nil
}

func getDeviceGeneration(device Pixel) string {
	switch device {
	case Pixel4, Pixel4XL, Pixel4a, Pixel4a5G:
		return "Pixel 4"
	case Pixel5, Pixel5a:
		return "Pixel 5"
	case Pixel6, Pixel6Pro, Pixel6a:
		return "Pixel 6"
	case Pixel7, Pixel7Pro, Pixel7a:
		return "Pixel 7"
	case Pixel8, Pixel8Pro, Pixel8a:
		return "Pixel 8"
	case Pixel9, Pixel9Pro, Pixel9ProXL, Pixel9ProFold:
		return "Pixel 9"
	default:
		return "Unknown"
	}
}
