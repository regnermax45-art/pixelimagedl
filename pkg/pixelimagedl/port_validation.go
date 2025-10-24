package pixelimagedl

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

// ValidatePortingPrerequisites checks if all necessary conditions are met for porting
func ValidatePortingPrerequisites(sourceDevice, targetDevice Pixel, algorithm PortingAlgorithm) error {
	// Check if devices are supported
	if !isDeviceSupported(sourceDevice) {
		return fmt.Errorf("source device %s is not supported for porting", sourceDevice.String())
	}

	if !isDeviceSupported(targetDevice) {
		return fmt.Errorf("target device %s is not supported for porting", targetDevice.String())
	}

	// Check algorithm compatibility
	if err := validateAlgorithmCompatibility(sourceDevice, targetDevice, algorithm); err != nil {
		return err
	}

	// Log warnings for risky operations
	logRiskWarnings(sourceDevice, targetDevice, algorithm)

	return nil
}

// isDeviceSupported checks if a device is supported for porting operations
func isDeviceSupported(device Pixel) bool {
	supportedDevices := []Pixel{
		Pixel6, Pixel6Pro, Pixel6a,
		Pixel7, Pixel7Pro, Pixel7a,
		Pixel8, Pixel8Pro, Pixel8a,
		Pixel9, Pixel9Pro, Pixel9ProXL, Pixel9ProFold,
	}

	for _, supported := range supportedDevices {
		if device == supported {
			return true
		}
	}

	return false
}

// validateAlgorithmCompatibility checks if the chosen algorithm is appropriate for the device pair
func validateAlgorithmCompatibility(source, target Pixel, algorithm PortingAlgorithm) error {
	sourceGen := getDeviceGeneration(source)
	targetGen := getDeviceGeneration(target)

	// Experimental algorithm should only be used within the same generation
	if algorithm == ExperimentalAlgorithm && sourceGen != targetGen {
		return errors.New("experimental algorithm can only be used within the same device generation")
	}

	// Warn about cross-generation porting with advanced algorithm
	if algorithm == AdvancedAlgorithm && sourceGen != targetGen {
		log.Printf("⚠️  WARNING: Using advanced algorithm for cross-generation porting (%s → %s) is not recommended\n", sourceGen, targetGen)
	}

	return nil
}

// logRiskWarnings logs appropriate warnings based on the porting operation
func logRiskWarnings(sourceDevice, targetDevice Pixel, algorithm PortingAlgorithm) {
	sourceGen := getDeviceGeneration(sourceDevice)
	targetGen := getDeviceGeneration(targetDevice)

	if sourceGen != targetGen {
		log.Printf("⚠️  CROSS-GENERATION PORTING: Porting from %s (%s) to %s (%s)\n",
			sourceDevice.String(), sourceGen, targetDevice.String(), targetGen)
		log.Printf("⚠️  This operation has higher risk of compatibility issues\n")
	}

	switch algorithm {
	case ExperimentalAlgorithm:
		log.Printf("🚨 EXPERIMENTAL ALGORITHM SELECTED\n")
		log.Printf("🚨 This algorithm makes significant system modifications\n")
		log.Printf("🚨 BOOT FAILURE OR DEVICE BRICKING IS POSSIBLE\n")
		log.Printf("🚨 ENSURE YOU HAVE A COMPLETE BACKUP BEFORE PROCEEDING\n")

	case AdvancedAlgorithm:
		log.Printf("⚠️  ADVANCED ALGORITHM SELECTED\n")
		log.Printf("⚠️  This algorithm modifies core system components\n")
		log.Printf("⚠️  Test the ported firmware in a safe environment first\n")

	case BasicAlgorithm:
		log.Printf("ℹ️  BASIC ALGORITHM SELECTED\n")
		log.Printf("ℹ️  This provides minimal modifications with lower risk\n")
		log.Printf("ℹ️  May require additional manual configuration\n")
	}
}

// ValidatePortingResult performs post-porting validation
func ValidatePortingResult(result *PortingResult) error {
	if result == nil {
		return errors.New("porting result is nil")
	}

	// Check if the ported file exists
	if result.PortedImagePath == "" {
		return errors.New("ported image path is empty")
	}

	// Additional validation could be added here
	// For example: checksum validation, file size checks, etc.

	return nil
}
