package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

var fixZipCmd = cli.Command{
	Name:  "fix-zip",
	Usage: "Fix ZIP extraction issues by recreating with better compatibility",
	Description: "Recreates ZIP files with no compression and proper headers for better extraction compatibility",
	Action: runFixZip,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "input",
			Aliases: []string{"i"},
			Usage:   "Input ZIP file to fix",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output fixed ZIP file",
			Required: true,
		},
	},
}

func runFixZip(ctx context.Context, cmd *cli.Command) error {
	inputFile := cmd.String("input")
	outputFile := cmd.String("output")

	log.Printf("🔧 Fixing ZIP file: %s → %s", inputFile, outputFile)

	// Open the problematic ZIP file
	reader, err := zip.OpenReader(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input ZIP: %v", err)
	}
	defer reader.Close()

	// Create the fixed output ZIP file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	writer := zip.NewWriter(outFile)
	defer writer.Close()

	log.Printf("📦 Processing %d files...", len(reader.File))

	for i, file := range reader.File {
		log.Printf("🔄 [%d/%d] Processing: %s", i+1, len(reader.File), file.Name)

		// Create new header with no compression for better compatibility
		header := &zip.FileHeader{
			Name:               file.Name,
			Method:             zip.Store, // No compression
			UncompressedSize64: file.UncompressedSize64,
		}
		header.SetMode(file.Mode())
		header.SetModTime(file.Modified)

		// Create writer for this file
		w, err := writer.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create header for %s: %v", file.Name, err)
		}

		// Open the file from the original ZIP
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open %s: %v", file.Name, err)
		}

		// Copy the content
		_, err = io.Copy(w, rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("failed to copy %s: %v", file.Name, err)
		}

		log.Printf("✅ Fixed: %s", file.Name)
	}

	log.Printf("🎉 ZIP file fixed successfully!")
	log.Printf("📁 Output: %s", outputFile)

	// Show file sizes
	inputStat, _ := os.Stat(inputFile)
	outputStat, _ := os.Stat(outputFile)
	
	log.Printf("📊 Size comparison:")
	log.Printf("   Original: %.2f MB", float64(inputStat.Size())/(1024*1024))
	log.Printf("   Fixed:    %.2f MB", float64(outputStat.Size())/(1024*1024))

	return nil
}
