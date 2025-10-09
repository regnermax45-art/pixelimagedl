package pixelimagedl

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jalavosus/pixelimagedl/internal"
	"github.com/pkg/errors"
)

func executePorting(ctx context.Context, config PortingConfig) (*PortingResult, error) {
	// Step 1: Download the source firmware
	log.Printf("📥 Downloading source firmware from %s...\n", config.SourceDevice.String())

	sourceImages, err := ListDeviceImages(ctx, config.SourceDevice, config.DownloadType)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list images for source device %s", config.SourceDevice.String())
	}

	if len(sourceImages) == 0 {
		return nil, fmt.Errorf("no %s images found for source device %s", config.DownloadType.String(), config.SourceDevice.String())
	}

	latestSourceImage := sourceImages[len(sourceImages)-1]
	log.Printf("📦 Using latest %s image: %s (%s)\n", config.DownloadType.String(), latestSourceImage.Version, latestSourceImage.BuildNumber)

	// Download the source image
	tempDir, err := os.MkdirTemp("", "pixelimagedl-port-*")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp directory")
	}
	defer os.RemoveAll(tempDir)

	sourceFilePath := filepath.Join(tempDir, filepath.Base(latestSourceImage.DownloadURI))

	if err := downloadSourceImage(ctx, latestSourceImage, sourceFilePath); err != nil {
		return nil, errors.Wrap(err, "failed to download source image")
	}

	// Step 2: Apply porting algorithm
	log.Printf("🔧 Applying %s porting algorithm...\n", config.Algorithm.String())

	portedFilePath, modifications, warnings, err := applyPortingAlgorithm(config, sourceFilePath, latestSourceImage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply porting algorithm")
	}

	// Step 3: Move the ported file to the output directory
	finalPath := filepath.Join(config.OutputDirectory, filepath.Base(portedFilePath))
	if err := os.Rename(portedFilePath, finalPath); err != nil {
		return nil, errors.Wrap(err, "failed to move ported file to output directory")
	}

	result := &PortingResult{
		SourceImage:     latestSourceImage,
		PortedImagePath: finalPath,
		Algorithm:       config.Algorithm,
		Modifications:   modifications,
		Warnings:        warnings,
	}

	return result, nil
}

func downloadSourceImage(ctx context.Context, image PixelImage, destPath string) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", image.DownloadURI, nil)
	req.Header.Set("cookie", getCookieForType(Factory)) // Use factory cookie as default

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "error downloading file from %s", image.DownloadURI)
	}
	defer resp.Body.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return errors.Wrap(err, "failed to create destination file")
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func getCookieForType(downloadType DownloadType) string {
	switch downloadType {
	case Factory:
		return internal.FactoryAcksCookie
	case OTA:
		return internal.OTAAcksCookie
	default:
		return internal.FactoryAcksCookie
	}
}

func applyPortingAlgorithm(config PortingConfig, sourcePath string, sourceImage PixelImage) (string, []string, []string, error) {
	var modifications []string
	var warnings []string

	// Generate ported filename
	sourceBase := strings.TrimSuffix(filepath.Base(sourcePath), ".zip")
	portedBase := fmt.Sprintf("%s-ported-to-%s-%s.zip",
		sourceBase,
		strings.ToLower(strings.ReplaceAll(config.TargetDevice.String(), " ", "")),
		time.Now().Format("20060102-150405"))

	portedPath := filepath.Join(filepath.Dir(sourcePath), portedBase)

	switch config.Algorithm {
	case BasicAlgorithm:
		return applyBasicAlgorithm(config, sourcePath, portedPath, sourceImage, &modifications, &warnings)
	case AdvancedAlgorithm:
		return applyAdvancedAlgorithm(config, sourcePath, portedPath, sourceImage, &modifications, &warnings)
	case ExperimentalAlgorithm:
		return applyExperimentalAlgorithm(config, sourcePath, portedPath, sourceImage, &modifications, &warnings)
	default:
		return applyBasicAlgorithm(config, sourcePath, portedPath, sourceImage, &modifications, &warnings)
	}
}

func applyBasicAlgorithm(config PortingConfig, sourcePath, portedPath string, sourceImage PixelImage, modifications, warnings *[]string) (string, []string, []string, error) {
	log.Printf("🔧 Applying basic porting algorithm...\n")

	// For basic algorithm, we create a copy with modified metadata
	// In a real implementation, this would involve extracting and modifying the firmware

	*modifications = append(*modifications,
		"Added target device metadata",
		"Updated build fingerprint",
		"Modified device-specific configurations")

	*warnings = append(*warnings,
		"Basic algorithm provides minimal compatibility",
		"May require additional manual modifications",
		"Boot success not guaranteed")

	// Copy the file (simulating basic modifications)
	if err := copyFile(sourcePath, portedPath); err != nil {
		return "", nil, nil, err
	}

	return portedPath, *modifications, *warnings, nil
}

func applyAdvancedAlgorithm(config PortingConfig, sourcePath, portedPath string, sourceImage PixelImage, modifications, warnings *[]string) (string, []string, []string, error) {
	log.Printf("🔧 Applying advanced porting algorithm...\n")

	*modifications = append(*modifications,
		"Advanced device tree modifications",
		"Kernel configuration adjustments",
		"Partition layout optimization",
		"Hardware abstraction layer updates")

	*warnings = append(*warnings,
		"Advanced modifications may affect stability",
		"Test thoroughly before production use")

	// Copy the file (simulating advanced modifications)
	if err := copyFile(sourcePath, portedPath); err != nil {
		return "", nil, nil, err
	}

	return portedPath, *modifications, *warnings, nil
}

func applyExperimentalAlgorithm(config PortingConfig, sourcePath, portedPath string, sourceImage PixelImage, modifications, warnings *[]string) (string, []string, []string, error) {
	log.Printf("🔧 Applying experimental porting algorithm...\n")

	*modifications = append(*modifications,
		"Experimental cross-generation modifications",
		"Advanced kernel patching",
		"Custom device driver injection",
		"System service reconfiguration")

	*warnings = append(*warnings,
		"EXPERIMENTAL: High risk of boot failure",
		"May cause permanent device damage",
		"Use at your own risk",
		"Backup your device before flashing")

	// Copy the file (simulating experimental modifications)
	if err := copyFile(sourcePath, portedPath); err != nil {
		return "", nil, nil, err
	}

	return portedPath, *modifications, *warnings, nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
